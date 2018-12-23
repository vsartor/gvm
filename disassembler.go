package main

import (
	"fmt"
	"os"
)

func disassemble(filePath string, ctxt context) {
	ctxt.l.Println("Starting to disassemble.")

	ctxt.l.Printf("Opening file '%s' now.\n", filePath)
	file, err := os.Open(filePath)
	if err != nil {
		ctxt.l.Printf("Error while opening '%s'.\n", filePath)
		ctxt.l.Fatalf("err: %v\n", err.Error())
	}
	defer file.Close()

	code := readCode(file, ctxt)

	ctxt.l.Println("Pretty printing now.")

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
			ctxt.l.Fatalf("Unexpected instruction %d.\n", code[ptr])
		}
	}
}
