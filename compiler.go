package main

import (
	"bufio"
	"encoding/binary"
	"log"
	"os"
	"strconv"
	"strings"
)

const binHeader int64 = 201812231

func expectArgN(tokens string, n, tot int, l *log.Logger) {
	if n != tot {
		l.Fatalf("Instruction `%s` expected %d arguments but got %d.\n",
			tokens, n, tot)
	}
}

func parseRegister(token string, l *log.Logger) int64 {
	if token[0] != 'r' && token[0] != 'f' {
		l.Fatalf("While parsing register expected 'r' or 'f' and got '%c'.",
			token[0])
	}

	reg, err := strconv.ParseInt(token[1:], 10, 64)
	if err != nil {
		l.Printf("Failed to parse register.\n")
		l.Printf("err: %v\n", err.Error())
		os.Exit(1)
	}

	return reg
}

func parseInt(token string, l *log.Logger) int64 {
	val, err := strconv.ParseInt(token, 10, 64)
	if err != nil {
		l.Printf("Failed to parse integer.\n")
		l.Printf("err: %v\n", err.Error())
		os.Exit(1)
	}

	return val
}

func parse(src *os.File, l *log.Logger) []int64 {
	l.Printf("Parsing file '%s'.\n", src.Name())

	code := make([]int64, 0, 64)

	label2pos := make(map[string]int64)
	pos2label := make(map[int64]string)

	label2pos["_zero"] = 0 // Adds the implicitly defined '_zero' label
	lastLabel := "_zero"

	pos := int64(0)

	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
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
				l.Printf("Read sublabel %s for position %d.", label, pos)
			} else {
				l.Printf("Read label %s for position %d.", label, pos)
				lastLabel = label
			}

			// Do not allow labels to be rewritten
			if labPos, ok := label2pos[label]; ok {
				l.Fatalf("Rewritten label '%s' from %d at %d.",
					label, labPos, pos)
			}

			label2pos[label] = pos
			continue
		}

		// Parse instruction tokens
		switch inst := tokens[0]; inst {
		case "halt":
			expectArgN(inst, 1, len(tokens), l)
			code = append(code, str2it[inst])
			pos++
		case "set":
			expectArgN(inst, 3, len(tokens), l)
			code = append(code, str2it[inst])
			code = append(code, parseRegister(tokens[1], l))
			code = append(code, parseInt(tokens[2], l))
			pos += 3
		case "add", "sub", "mul", "div", "rem", "cmp":
			expectArgN(inst, 3, len(tokens), l)
			code = append(code, str2it[inst])
			code = append(code, parseRegister(tokens[1], l))
			code = append(code, parseRegister(tokens[2], l))
			pos += 3
		case "jmp", "jeq", "jne", "jgt", "jlt", "jge", "jle":
			expectArgN(inst, 2, len(tokens), l)
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
		case "show":
			expectArgN(inst, 2, len(tokens), l)
			code = append(code, str2it[inst])
			code = append(code, parseRegister(tokens[1], l))
			pos += 2
		default:
			l.Fatalf("Unknown instruction '%s'.", inst)
		}
	}

	if err := scanner.Err(); err != nil {
		l.Printf("Error while reading the file.\n")
		l.Printf("err: %v\n", err.Error())
		os.Exit(1)
	}

	// Always append a halt into the end of the code
	code = append(code, IThalt)

	l.Println("Setting labels to values.")

	// Set labels to their values
	for codePos, label := range pos2label {
		jmpPos, ok := label2pos[label]
		if !ok {
			l.Fatalf("At position %d invalid reference to unknown label '%s'.",
				codePos, label)
		}
		l.Printf("Setting address at %d for label %s.", codePos, label)
		code[codePos] = jmpPos
	}

	l.Println("Finished parsing.")

	return code
}

func writeCode(code []int64, dst *os.File, l *log.Logger) {
	l.Println("Writing binary file.")

	// Header
	err := binary.Write(dst, binary.LittleEndian, binHeader)
	if err != nil {
		l.Println("Failed writting header to file.")
		l.Fatalf("err: %s\n", err.Error())
	}

	// Size of the code
	err = binary.Write(dst, binary.LittleEndian, int64(len(code)))
	if err != nil {
		l.Println("Failed writting size of code to file.")
		l.Fatalf("err: %s\n", err.Error())
	}

	// Write the code
	for _, tok := range code {
		err = binary.Write(dst, binary.LittleEndian, tok)
		if err != nil {
			l.Println("Failed writting code token to file.")
			l.Fatalf("err: %s\n", err.Error())
		}
	}

	l.Println("Finished writting binary file.")
}

func readCode(file *os.File, l *log.Logger) []int64 {
	l.Println("Reading binary file.")

	// Read and validate header
	var header int64
	err := binary.Read(file, binary.LittleEndian, &header)
	if err != nil {
		l.Println("Failed reading header from file.")
		l.Fatalf("err: %s\n", err.Error())
	}
	if header != binHeader {
		l.Fatalf("Expected header %d but got %d.\n", binHeader, header)
	}

	// Read code size and allocate slice
	var codeSize int64
	err = binary.Read(file, binary.LittleEndian, &codeSize)
	if err != nil {
		l.Println("Failed reading code size from file.")
		l.Fatalf("err: %s\n", err.Error())
	}
	code := make([]int64, codeSize)

	// Read the code
	for i := 0; i < int(codeSize); i++ {
		err = binary.Read(file, binary.LittleEndian, &code[i])
		if err != nil {
			l.Println("Failed reading code token from file.")
			l.Fatalf("err: %s\n", err.Error())
		}
	}

	l.Println("Finished reading binary file.")

	return code
}

func compile(srcPath, dstPath string, l *log.Logger) {
	// Open source file
	l.Printf("Trying to open '%s'.\n", srcPath)
	src, err := os.Open(srcPath)
	if err != nil {
		l.Printf("Error while opening '%s'.\n", srcPath)
		l.Fatalf("err: %v\n", err.Error())
	}
	defer src.Close()

	// Parse the file
	code := parse(src, l)

	// Create object file
	l.Printf("Trying to open '%s'.\n", dstPath)
	dst, err := os.Create(dstPath)
	if err != nil {
		l.Printf("Error while opening '%s'.\n", dstPath)
		l.Fatalf("err: %v\n", err.Error())
	}
	defer dst.Close()

	// Write it in binary form into dst
	writeCode(code, dst, l)
}
