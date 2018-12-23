package main

import (
	"bufio"
	"encoding/binary"
	"os"
	"strconv"
	"strings"
)

const binHeader int64 = 201812241

func expectArgN(tokens string, n, tot int, ctxt context) {
	if n != tot {
		ctxt.l.Fatalf("l%d: Token `%s` expected %d args, got %d.\n",
			ctxt.lineNum, tokens, n, tot)
	}
}

func parseRegister(token string, ctxt context) int64 {
	if token[0] != 'r' && token[0] != 'f' {
		ctxt.l.Fatalf("l%d: Parsing register: expected 'r|f', got '%c'.",
			ctxt.lineNum, token[0])
	}

	reg, err := strconv.ParseInt(token[1:], 10, 64)
	if err != nil {
		ctxt.l.Printf("Failed to parse register.\n")
		ctxt.l.Fatalf("err: %v\n", err.Error())
	}

	return reg
}

func parseInt(token string, ctxt context) int64 {
	val, err := strconv.ParseInt(token, 10, 64)
	if err != nil {
		ctxt.l.Printf("Failed to parse integer.\n")
		ctxt.l.Fatalf("err: %v\n", err.Error())
	}

	return val
}

func parse(src *os.File, ctxt context) []int64 {
	ctxt.l.Printf("Parsing file '%s'.\n", src.Name())

	code := make([]int64, 0, 64)

	label2pos := make(map[string]int64)
	pos2label := make(map[int64]string)

	label2pos["_zero"] = 0 // Adds the implicitly defined '_zero' label
	lastLabel := "_zero"

	pos := int64(0)

	ctxt.lineNum = 0

	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		ctxt.lineNum++

		line := scanner.Text()
		tokens := strings.Fields(line)

		if len(tokens) == 0 {
			continue
		}

		// Remove comments
		for idx, token := range tokens {
			if token[0] == ';' {
				tokens = tokens[:idx]
			}
		}
		if len(tokens) == 0 {
			continue
		}

		// Check if it's a label
		if tok := tokens[0]; tok[len(tok)-1] == ':' {
			label := tok[:len(tok)-1]

			// Check it it's a sublabel
			if label[0] == '.' {
				label = lastLabel + label
				if ctxt.verbose {
					ctxt.l.Printf("Read sublabel %s at %d.", label, pos)
				}
			} else {
				if ctxt.verbose {
					ctxt.l.Printf("Read label %s at %d.", label, pos)
				}
				lastLabel = label
			}

			// Do not allow labels to be rewritten
			if labPos, ok := label2pos[label]; ok {
				ctxt.l.Fatalf("l%d: Rewritten label '%s' from %d at %d.",
					ctxt.lineNum, label, labPos, pos)
			}

			label2pos[label] = pos
			continue
		}

		// Parse instruction tokens
		switch inst := tokens[0]; inst {
		case "halt", "ret":
			expectArgN(inst, 1, len(tokens), ctxt)
			code = append(code, str2it[inst])
			pos++
		case "set":
			expectArgN(inst, 3, len(tokens), ctxt)
			code = append(code, str2it[inst])
			code = append(code, parseRegister(tokens[1], ctxt))
			code = append(code, parseInt(tokens[2], ctxt))
			pos += 3
		case "mov", "add", "sub", "mul", "div", "rem", "cmp":
			expectArgN(inst, 3, len(tokens), ctxt)
			code = append(code, str2it[inst])
			code = append(code, parseRegister(tokens[1], ctxt))
			code = append(code, parseRegister(tokens[2], ctxt))
			pos += 3
		case "jmp", "jeq", "jne", "jgt", "jlt", "jge", "jle", "call":
			expectArgN(inst, 2, len(tokens), ctxt)
			code = append(code, str2it[inst])
			// Add placeholder for the label and make a pending parse
			code = append(code, 0)
			// Expand if it's a sublabel
			if tokens[1][0] == '.' {
				pos2label[pos+1] = lastLabel + tokens[1]
			} else {
				pos2label[pos+1] = tokens[1]
			}
			pos += 2
		case "show", "inc", "dec", "push", "pop":
			expectArgN(inst, 2, len(tokens), ctxt)
			code = append(code, str2it[inst])
			code = append(code, parseRegister(tokens[1], ctxt))
			pos += 2
		default:
			ctxt.l.Fatalf("l%d: Unknown instruction '%s'.", ctxt.lineNum, inst)
		}
	}

	if err := scanner.Err(); err != nil {
		ctxt.l.Printf("Error while reading the file.\n")
		ctxt.l.Fatalf("err: %v\n", err.Error())
	}

	// Always append a halt into the end of the code
	code = append(code, IThalt)

	ctxt.l.Println("Setting labels to values.")

	// Check if main was defined, if it is two tokens are added at the
	// beginning. Therefore an offset of 2 should be added to every position
	// to account for it.
	mainPos, ok := label2pos["main"]
	mainOffset := int64(2)
	if !ok {
		mainOffset = 0
		ctxt.l.Println("WARNING: The label `main` has not been defined.")
	} else {
		code = append([]int64{ITjmp, mainPos + mainOffset}, code...)
	}

	// Set labels to their values
	for codePos, label := range pos2label {
		jmpPos, ok := label2pos[label]
		if !ok {
			ctxt.l.Fatalf("At position %d reference to unknown label '%s'.",
				codePos, label)
		}
		// l.Printf("Setting address at %d for label %s.", codePos, label)
		code[codePos+mainOffset] = jmpPos + mainOffset
	}

	ctxt.l.Println("Finished parsing.")

	return code
}

func writeCode(code []int64, dst *os.File, ctxt context) {
	ctxt.l.Println("Writing binary file.")

	// Header
	err := binary.Write(dst, binary.LittleEndian, binHeader)
	if err != nil {
		ctxt.l.Println("Failed writting header to file.")
		ctxt.l.Fatalf("err: %s\n", err.Error())
	}

	// Size of the code
	err = binary.Write(dst, binary.LittleEndian, int64(len(code)))
	if err != nil {
		ctxt.l.Println("Failed writting size of code to file.")
		ctxt.l.Fatalf("err: %s\n", err.Error())
	}

	// Write the code
	for _, tok := range code {
		err = binary.Write(dst, binary.LittleEndian, tok)
		if err != nil {
			ctxt.l.Println("Failed writting code token to file.")
			ctxt.l.Fatalf("err: %s\n", err.Error())
		}
	}

	ctxt.l.Println("Finished writting binary file.")
}

func readCode(file *os.File, ctxt context) []int64 {
	ctxt.l.Println("Reading binary file.")

	// Read and validate header
	var header int64
	err := binary.Read(file, binary.LittleEndian, &header)
	if err != nil {
		ctxt.l.Println("Failed reading header from file.")
		ctxt.l.Fatalf("err: %s\n", err.Error())
	}
	if header != binHeader {
		ctxt.l.Fatalf("Expected header %d but got %d.\n", binHeader, header)
	}

	// Read code size and allocate slice
	var codeSize int64
	err = binary.Read(file, binary.LittleEndian, &codeSize)
	if err != nil {
		ctxt.l.Println("Failed reading code size from file.")
		ctxt.l.Fatalf("err: %s\n", err.Error())
	}
	code := make([]int64, codeSize)

	// Read the code
	for i := 0; i < int(codeSize); i++ {
		err = binary.Read(file, binary.LittleEndian, &code[i])
		if err != nil {
			ctxt.l.Println("Failed reading code token from file.")
			ctxt.l.Fatalf("err: %s\n", err.Error())
		}
	}

	ctxt.l.Println("Finished reading binary file.")

	return code
}

func compile(srcPath, dstPath string, ctxt context) {
	// Open source file
	ctxt.l.Printf("Trying to open '%s'.\n", srcPath)
	src, err := os.Open(srcPath)
	if err != nil {
		ctxt.l.Printf("Error while opening '%s'.\n", srcPath)
		ctxt.l.Fatalf("err: %v\n", err.Error())
	}
	defer src.Close()

	// Parse the file
	code := parse(src, ctxt)

	// Create object file
	ctxt.l.Printf("Trying to open '%s'.\n", dstPath)
	dst, err := os.Create(dstPath)
	if err != nil {
		ctxt.l.Printf("Error while opening '%s'.\n", dstPath)
		ctxt.l.Fatalf("err: %v\n", err.Error())
	}
	defer dst.Close()

	// Write it in binary form into dst
	writeCode(code, dst, ctxt)
}
