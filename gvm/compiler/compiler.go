package compiler

import (
	"bufio"
	"encoding/binary"
	"github.com/vsartor/gvm/gvm"
	"github.com/vsartor/gvm/gvm/lang"
	"os"
	"path"
	"strings"
)

func assertArgCount(instruction gvm.Code, expectedCount, argCount int, ctxt gvm.Context) {
	if expectedCount != argCount {
		gvm.Logger.Criticalf("%s.%d: Token `%s` expected %d arguments, got %d.\n",
			ctxt.FileName, ctxt.LineNum, lang.ToString(instruction), expectedCount, argCount)
		os.Exit(1)
	}
}

func expandSublabel(sublabel, lastLabel string) string {
	return lastLabel + sublabel
}

func compile(src *os.File, ctxt gvm.Context) []gvm.Code {
	gvm.Logger.Infof("Parsing file '%s'.\n", src.Name())

	code := make([]gvm.Code, 0, gvm.CodeArrayInitialSize)

	// Because we may need to add the `jmp main` tokens in case a
	// `main` label is present, let's save the space for this instruction
	// by adding in two `noop`'s.
	code = append(code, lang.Noop)
	code = append(code, lang.Noop)

	labelToPosition := make(map[string]int64)
	positionToLabel := make(map[int64]string)

	lastLabel := ""
	currCodePosition := int64(2)
	ctxt.LineNum = 0

	gvm.Logger.Infof("Parser pass starting.\n")

	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		ctxt.LineNum++

		line := scanner.Text()
		tokens := strings.Fields(line)

		// Skip empty lines straight away
		if len(tokens) == 0 {
			continue
		}

		// Remove comments if there are any
		for idx, token := range tokens {
			if token[0] == ';' {
				tokens = tokens[:idx]
			}
		}

		// Remove lines which became empty after removing comments
		if len(tokens) == 0 {
			continue
		}

		// Check if line contains a label
		if tok := tokens[0]; tok[len(tok)-1] == ':' {
			labelName := tok[:len(tok)-1]

			if labelName[0] == '.' {
				//  If it's a sublabel, expand its name to include the labelName
				if len(lastLabel) == 0 {
					gvm.Logger.Criticalf("%s.%d: Orphan sublabel '%s' found.",
						ctxt.FileName, ctxt.LineNum, labelName)
					os.Exit(1)
				}
				labelName = expandSublabel(labelName, lastLabel)

				gvm.Logger.Debugf("%s.%d: Read sublabel %s at %d.",
					ctxt.FileName, ctxt.LineNum, labelName, currCodePosition)
			} else {
				// Remember the last labelName
				gvm.Logger.Debugf("%s.%d: Read label %s at %d.",
					ctxt.FileName, ctxt.LineNum, labelName, currCodePosition)
				lastLabel = labelName
			}

			// Do not allow multiple instances of the same labelName
			if _, ok := labelToPosition[labelName]; ok {
				gvm.Logger.Criticalf("%s.%d: Attempt to overwrite label '%s'.",
					ctxt.FileName, ctxt.LineNum, labelName)
				os.Exit(1)
			}

			labelToPosition[labelName] = currCodePosition
			continue
		}

		// Parse the instruction
		switch instruction := parseInstruction(tokens[0], ctxt); instruction {
		case lang.Halt, lang.Ret, lang.Noop:
			assertArgCount(instruction, 1, len(tokens), ctxt)
			code = append(code, instruction)
			currCodePosition++
		case lang.Const:
			assertArgCount(instruction, 3, len(tokens), ctxt)
			code = append(code, instruction)
			code = append(code, parseInt(tokens[1], ctxt))
			code = append(code, parseRegister(tokens[2], ctxt))
			currCodePosition += 3
		case lang.Mov, lang.Add, lang.Sub, lang.Mul, lang.Div, lang.Rem, lang.Cmp:
			assertArgCount(instruction, 3, len(tokens), ctxt)
			code = append(code, instruction)
			code = append(code, parseRegister(tokens[1], ctxt))
			code = append(code, parseRegister(tokens[2], ctxt))
			currCodePosition += 3
		case lang.Jmp, lang.Jeq, lang.Jne, lang.Jgt, lang.Jlt, lang.Jge, lang.Jle, lang.Jerr, lang.Call:
			assertArgCount(instruction, 2, len(tokens), ctxt)
			code = append(code, instruction)
			// Add a placeholder for the code srcPosition
			// There will be a pass later to fill it in
			code = append(code, 0)

			// Save information for latter pass so that we can fill in the
			// correct code srcPosition later on, making sure we expand the
			labelName := tokens[1]
			if labelName[0] == '.' {
				positionToLabel[currCodePosition+1] = expandSublabel(labelName, lastLabel)
			} else {
				positionToLabel[currCodePosition+1] = labelName
			}
			currCodePosition += 2
		case lang.Show, lang.Inc, lang.Dec, lang.Push, lang.Pop, lang.Iarg:
			assertArgCount(instruction, 2, len(tokens), ctxt)
			code = append(code, instruction)
			code = append(code, parseRegister(tokens[1], ctxt))
			currCodePosition += 2
		default:
			gvm.Logger.Criticalf("%s.%d: Unknown instruction code %d.",
				ctxt.FileName, ctxt.LineNum, instruction)
			os.Exit(1)
		}
	}

	if err := scanner.Err(); err != nil {
		gvm.Logger.Criticalf("Error while reading the file: %s.\n", err.Error())
		os.Exit(1)
	}

	gvm.Logger.Infof("Label pass starting.\n")

	// Check if `main` was defined. If it has, change the `noop`s into the correct
	// instruction.
	if mainPosition, ok := labelToPosition["main"]; ok {
		code[0] = lang.Jmp
		code[1] = gvm.Code(mainPosition)
	} else {
		gvm.Logger.Infof("The label `main` was not defined.\n")
	}

	// Fill in code positions based on the labels
	for srcPosition, label := range positionToLabel {
		dstPosition, ok := labelToPosition[label]
		if !ok {
			gvm.Logger.Criticalf("At position %d reference to unknown label '%s'.", srcPosition, label)
			os.Exit(1)
		}
		code[srcPosition] = gvm.Code(dstPosition)
	}

	gvm.Logger.Infof("Finished parsing.\n")

	return code
}

func writeCode(code []gvm.Code, output *os.File) {
	gvm.Logger.Infof("Writing binary file.\n")

	// Header
	err := binary.Write(output, binary.LittleEndian, gvm.BinaryFileHeader)
	if err != nil {
		gvm.Logger.Criticalf("Failed writting header: %s\n", err.Error())
		os.Exit(1)
	}

	// Size of the code
	err = binary.Write(output, binary.LittleEndian, int64(len(code)))
	if err != nil {
		gvm.Logger.Criticalf("Failed writting size of code array: %s\n", err.Error())
		os.Exit(1)
	}

	// Write the code
	for _, tok := range code {
		err = binary.Write(output, binary.LittleEndian, tok)
		if err != nil {
			gvm.Logger.Criticalf("Failed writting code token: %s\n", err.Error())
			os.Exit(1)
		}
	}

	gvm.Logger.Infof("Finished writting binary file.\n")
}

func ReadCode(file *os.File) []gvm.Code {
	gvm.Logger.Infof("Reading binary file.\n")

	// Read and validate header
	var header int64
	err := binary.Read(file, binary.LittleEndian, &header)
	if err != nil {
		gvm.Logger.Criticalf("Failed reading header: %s\n", err.Error())
		os.Exit(1)
	}
	if header != gvm.BinaryFileHeader {
		gvm.Logger.Criticalf("Expected header %d but got %d.\n", gvm.BinaryFileHeader, header)
		os.Exit(1)
	}

	// Read code size and allocate slice
	var codeSize int64
	err = binary.Read(file, binary.LittleEndian, &codeSize)
	if err != nil {
		gvm.Logger.Criticalf("Failed reading code size: %s\n", err.Error())
		os.Exit(1)
	}
	code := make([]gvm.Code, codeSize)

	// Read the code
	for i := 0; i < int(codeSize); i++ {
		err = binary.Read(file, binary.LittleEndian, &code[i])
		if err != nil {
			gvm.Logger.Criticalf("Failed reading code token: %s\n", err.Error())
			os.Exit(1)
		}
	}

	gvm.Logger.Infof("Finished reading binary file.\n")

	return code
}

func Compile(srcPath, dstPath string) {
	// Open source file
	gvm.Logger.Infof("Opening '%s'.\n", srcPath)
	input, err := os.Open(srcPath)
	if err != nil {
		gvm.Logger.Criticalf("Failed opening '%s': %s\n", srcPath, err.Error())
		os.Exit(1)
	}
	defer input.Close()

	// Parse the file
	ctxt := gvm.Context{FileName: path.Base(srcPath)}
	code := compile(input, ctxt)

	// Create object file
	gvm.Logger.Infof("Opening '%s'.\n", dstPath)
	output, err := os.Create(dstPath)
	if err != nil {
		gvm.Logger.Criticalf("Failed opening '%s': %s\n", dstPath, err.Error())
		os.Exit(1)
	}
	defer output.Close()

	// Write it in binary form into output
	writeCode(code, output)
}
