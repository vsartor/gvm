package main

import (
	"fmt"
	"log"
	"os"
)

const regCount = 16

func run(filePath string, l *log.Logger) {
	l.Println("Starting to disassemble.")

	l.Printf("Opening file '%s' now.\n", filePath)
	file, err := os.Open(filePath)
	if err != nil {
		l.Printf("Error while opening '%s'.\n", filePath)
		l.Fatalf("err: %v\n", err.Error())
	}
	defer file.Close()

	code := readCode(file, l)

	l.Println("Starting execution now.")

	reg := make([]int64, regCount)

	ptr := 0
	for ptr < len(code) {
		switch code[ptr] {
		case IThalt:
			ptr = len(code)
		case ITset:
			reg[code[ptr+1]] = code[ptr+2]
			ptr += 3
		case ITadd:
			reg[code[ptr+2]] += reg[code[ptr+1]]
			ptr += 3
		case ITsub:
			reg[code[ptr+2]] -= reg[code[ptr+1]]
			ptr += 3
		case ITmul:
			reg[code[ptr+2]] *= reg[code[ptr+1]]
			ptr += 3
		case ITdiv:
			reg[code[ptr+2]] /= reg[code[ptr+1]]
			ptr += 3
		case ITrem:
			reg[code[ptr+2]] %= reg[code[ptr+1]]
			ptr += 3
		case ITshow:
			fmt.Printf("%d\n", reg[code[ptr+1]])
			ptr += 2
		default:
			l.Fatalf("Unexpected instruction %d.\n", code[ptr])
		}
	}
}