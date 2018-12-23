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
		switch it := code[ptr]; it {
		case IThalt:
			fmt.Printf("%s\n", it2str[it])
			ptr++
		case ITset:
			fmt.Printf("%s r%d %d\n", it2str[it], code[ptr+1], code[ptr+2])
			ptr += 3
		case ITadd, ITsub, ITmul, ITdiv, ITrem, ITcmp:
			fmt.Printf("%s r%d r%d\n", it2str[it], code[ptr+1], code[ptr+2])
			ptr += 3
		case ITjmp, ITjeq, ITjne, ITjgt, ITjlt, ITjge, ITjle:
			fmt.Printf("%s %d\n", it2str[it], code[ptr+1])
			ptr += 2
		case ITshow:
			fmt.Printf("%s r%d\n", it2str[it], code[ptr+1])
			ptr += 2
		default:
			l.Fatalf("Unexpected instruction %d.\n", code[ptr])
		}
	}
}
