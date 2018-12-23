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

		switch inst := tokens[0]; inst {
		case "halt":
			expectArgN(inst, 1, len(tokens), l)
			code = append(code, IThalt)
		case "set":
			expectArgN(inst, 3, len(tokens), l)
			code = append(code, ITset)
			code = append(code, parseRegister(tokens[1], l))
			code = append(code, parseInt(tokens[2], l))
		case "add":
			expectArgN(inst, 3, len(tokens), l)
			code = append(code, ITadd)
			code = append(code, parseRegister(tokens[1], l))
			code = append(code, parseRegister(tokens[2], l))
		case "sub":
			expectArgN(inst, 3, len(tokens), l)
			code = append(code, ITsub)
			code = append(code, parseRegister(tokens[1], l))
			code = append(code, parseRegister(tokens[2], l))
		case "mul":
			expectArgN(inst, 3, len(tokens), l)
			code = append(code, ITmul)
			code = append(code, parseRegister(tokens[1], l))
			code = append(code, parseRegister(tokens[2], l))
		case "div":
			expectArgN(inst, 3, len(tokens), l)
			code = append(code, ITdiv)
			code = append(code, parseRegister(tokens[1], l))
			code = append(code, parseRegister(tokens[2], l))
		case "rem":
			expectArgN(inst, 3, len(tokens), l)
			code = append(code, ITrem)
			code = append(code, parseRegister(tokens[1], l))
			code = append(code, parseRegister(tokens[2], l))
		case "show":
			expectArgN(inst, 2, len(tokens), l)
			code = append(code, ITshow)
			code = append(code, parseRegister(tokens[1], l))
		default:
			l.Fatalf("Unknown instruction '%s'.", inst)
		}
	}
	l.Println("Finished parsing.")

	if err := scanner.Err(); err != nil {
		l.Printf("Error while reading the file.\n")
		l.Printf("err: %v\n", err.Error())
		os.Exit(1)
	}

	// Always append a halt into the end of the code
	code = append(code, IThalt)

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
