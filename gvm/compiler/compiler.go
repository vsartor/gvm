package compiler

import (
	"bufio"
	"encoding/binary"
	"github.com/vsartor/gvm/gvm"
	"github.com/vsartor/gvm/gvm/lang"
	"os"
	"strings"
)

func assertArgCount(instruction gvm.Code, expectedCount, argCount int, ctxt gvm.Context) {
	if expectedCount != argCount {
		gvm.Logger.Fatalf("l%d: Token `%s` expected %d arguments, got %d.\n",
			ctxt.LineNum, lang.ToString(instruction), expectedCount, argCount)
	}
}

func expandSublabel(sublabel, lastLabel string) string {
	return lastLabel + sublabel
}

func compile(src *os.File, ctxt gvm.Context) []gvm.Code {
	gvm.Logger.Printf("Parsing file '%s'.\n", src.Name())

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

	gvm.Logger.Println("Parser pass.")

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
					gvm.Logger.Fatalf("l%d: Orphan sublabel '%s' found.", ctxt.LineNum, labelName)
				}
				labelName = expandSublabel(labelName, lastLabel)

				if ctxt.IsVerbose {
					gvm.Logger.Printf("l%d: Read sublabel %s at %d.",
						ctxt.LineNum, labelName, currCodePosition)
				}
			} else {
				// Remember the last labelName
				if ctxt.IsVerbose {
					gvm.Logger.Printf("l%d: Read label %s at %d.",
						ctxt.LineNum, labelName, currCodePosition)
				}
				lastLabel = labelName
			}

			// Do not allow multiple instances of the same labelName
			if _, ok := labelToPosition[labelName]; ok {
				gvm.Logger.Fatalf("l%d: Attempt to overwrite label '%s'.", ctxt.LineNum, labelName)
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
		case lang.Set:
			assertArgCount(instruction, 3, len(tokens), ctxt)
			code = append(code, instruction)
			code = append(code, parseRegister(tokens[1], ctxt))
			code = append(code, parseInt(tokens[2], ctxt))
			currCodePosition += 3
		case lang.Mov, lang.Add, lang.Sub, lang.Mul, lang.Div, lang.Rem, lang.Cmp:
			assertArgCount(instruction, 3, len(tokens), ctxt)
			code = append(code, instruction)
			code = append(code, parseRegister(tokens[1], ctxt))
			code = append(code, parseRegister(tokens[2], ctxt))
			currCodePosition += 3
		case lang.Jmp, lang.Jeq, lang.Jne, lang.Jgt, lang.Jlt, lang.Jge, lang.Jle, lang.Call:
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
		case lang.Show, lang.Inc, lang.Dec, lang.Push, lang.Pop:
			assertArgCount(instruction, 2, len(tokens), ctxt)
			code = append(code, instruction)
			code = append(code, parseRegister(tokens[1], ctxt))
			currCodePosition += 2
		default:
			gvm.Logger.Fatalf("l%d: Unknown instructionType '%v'.", ctxt.LineNum, instruction)
		}
	}

	if err := scanner.Err(); err != nil {
		gvm.Logger.Fatalf("Error while reading the file: %s.\n", err.Error())
	}

	gvm.Logger.Println("Label pass.")

	// Check if `main` was defined. If it has, change the `noop`s into the correct
	// instruction.
	if mainPosition, ok := labelToPosition["main"]; ok {
		code[0] = lang.Jmp
		code[1] = gvm.Code(mainPosition)
	} else {
		gvm.Logger.Println("The label `main` was not defined.")
	}

	// Fill in code positions based on the labels
	for srcPosition, label := range positionToLabel {
		dstPosition, ok := labelToPosition[label]
		if !ok {
			gvm.Logger.Fatalf("At srcPosition %d reference to unknown label '%s'.", srcPosition, label)
		}
		code[srcPosition] = gvm.Code(dstPosition)
	}

	gvm.Logger.Println("Finished parsing.")

	return code
}

func writeCode(code []gvm.Code, output *os.File) {
	gvm.Logger.Println("Writing binary file.")

	// Header
	err := binary.Write(output, binary.LittleEndian, gvm.BinaryFileHeader)
	if err != nil {
		gvm.Logger.Fatalf("Failed writting header: %s\n", err.Error())
	}

	// Size of the code
	err = binary.Write(output, binary.LittleEndian, int64(len(code)))
	if err != nil {
		gvm.Logger.Fatalf("Failed writting size of code array: %s\n", err.Error())
	}

	// Write the code
	for _, tok := range code {
		err = binary.Write(output, binary.LittleEndian, tok)
		if err != nil {
			gvm.Logger.Fatalf("Failed writting code token: %s\n", err.Error())
		}
	}

	gvm.Logger.Println("Finished writting binary file.")
}

func ReadCode(file *os.File) []gvm.Code {
	gvm.Logger.Println("Reading binary file.")

	// Read and validate header
	var header int64
	err := binary.Read(file, binary.LittleEndian, &header)
	if err != nil {
		gvm.Logger.Fatalf("Failed reading header: %s\n", err.Error())
	}
	if header != gvm.BinaryFileHeader {
		gvm.Logger.Fatalf("Expected header %d but got %d.\n", gvm.BinaryFileHeader, header)
	}

	// Read code size and allocate slice
	var codeSize int64
	err = binary.Read(file, binary.LittleEndian, &codeSize)
	if err != nil {
		gvm.Logger.Fatalf("Failed reading code size: %s\n", err.Error())
	}
	code := make([]gvm.Code, codeSize)

	// Read the code
	for i := 0; i < int(codeSize); i++ {
		err = binary.Read(file, binary.LittleEndian, &code[i])
		if err != nil {
			gvm.Logger.Fatalf("Failed reading code token: %s\n", err.Error())
		}
	}

	gvm.Logger.Println("Finished reading binary file.")

	return code
}

func Compile(srcPath, dstPath string, ctxt gvm.Context) {
	// Open source file
	gvm.Logger.Printf("Opening '%s'.\n", srcPath)
	input, err := os.Open(srcPath)
	if err != nil {
		gvm.Logger.Fatalf("Failed opening '%s': %s\n", srcPath, err.Error())
	}
	defer input.Close()

	// Parse the file
	code := compile(input, ctxt)

	// Create object file
	gvm.Logger.Printf("Opening '%s'.\n", dstPath)
	output, err := os.Create(dstPath)
	if err != nil {
		gvm.Logger.Fatalf("Failed opening '%s': %s\n", dstPath, err.Error())
	}
	defer output.Close()

	// Write it in binary form into output
	writeCode(code, output)
}
