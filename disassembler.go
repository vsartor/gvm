package main

import (
	"fmt"
	"log"
	"os"
)

func disassemble(filePath string, l *log.Logger) {
	l.Println("Starting to disassemble.")

	l.Printf("Opening file '%s' now.\n", filePath)
	file, err := os.Open(filePath)
	if err != nil {
		l.Printf("Error while opening '%s'.\n", filePath)
		l.Fatalf("err: %v\n", err.Error())
	}
	defer file.Close()

	code := readCode(file, l)

	l.Println("Pretty printing now.")

	ptr := 0
	for ptr < len(code) {
		switch code[ptr] {
		case IThalt:
			ptr++
		case ITset:
			fmt.Printf("set r%d %d\n", code[ptr+1], code[ptr+2])
			ptr += 3
		case ITadd:
			fmt.Printf("add r%d r%d\n", code[ptr+1], code[ptr+2])
			ptr += 3
		case ITsub:
			fmt.Printf("sub r%d r%d\n", code[ptr+1], code[ptr+2])
			ptr += 3
		case ITmul:
			fmt.Printf("mul r%d r%d\n", code[ptr+1], code[ptr+2])
			ptr += 3
		case ITdiv:
			fmt.Printf("div r%d r%d\n", code[ptr+1], code[ptr+2])
			ptr += 3
		case ITrem:
			fmt.Printf("rem r%d r%d\n", code[ptr+1], code[ptr+2])
			ptr += 3
		case ITshow:
			fmt.Printf("show r%d\n", code[ptr+1])
			ptr += 2
		default:
			l.Fatalf("Unexpected instruction %d.\n", code[ptr])
		}
	}
}
