package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func expectArgN(tokens string, n, tot int, l *log.Logger) {
	if n != tot {
		l.Fatalf("Instruction `%s` expected %d arguments but got %d.\n",
			tokens, n, tot)
	}
}

func parseRegister(token string, l *log.Logger) int {
	if token[0] != 'r' && token[0] != 'f' {
		l.Fatalf("While parsing register expected 'r' or 'f' and got '%c'.",
			token[0])
	}

	reg, err := strconv.Atoi(token[1:])
	if err != nil {
		l.Printf("Failed to parse register.\n")
		l.Printf("err: %v\n", err.Error())
		os.Exit(1)
	}

	return reg
}

func parseInt(token string, l *log.Logger) int {
	val, err := strconv.Atoi(token)
	if err != nil {
		l.Printf("Failed to parse integer.\n")
		l.Printf("err: %v\n", err.Error())
		os.Exit(1)
	}

	return val
}

func parse(src *os.File, l *log.Logger) []int {
	l.Printf("Parsing file '%s'.\n", src.Name())

	code := make([]int, 0, 64)

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

func compile(srcPath, dstPath string, l *log.Logger) {
	l.Printf("Trying to open '%s'.\n", srcPath)
	src, err := os.Open(srcPath)
	if err != nil {
		l.Printf("Error while opening '%s'.\n", srcPath)
		l.Printf("err: %v\n", err.Error())
		os.Exit(1)
	}
	defer src.Close()

	//TODO: In the future, open the file so we can write to it
	l.Printf("Ignoring destination file '%s' for now.\n", dstPath)

	// Parse the file
	code := parse(src, l)
	fmt.Printf("%v\n", code)

	// Write it in binary form into dst
}
