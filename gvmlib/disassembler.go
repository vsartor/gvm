package gvmlib

import (
	"fmt"
	"os"
)

func Disassemble(filePath string, ctxt Context) {
	ctxt.Logger.Println("Starting to disassemble.")

	ctxt.Logger.Printf("Opening file '%s' now.\n", filePath)
	file, err := os.Open(filePath)
	if err != nil {
		ctxt.Logger.Printf("Error while opening '%s'.\n", filePath)
		ctxt.Logger.Fatalf("err: %v\n", err.Error())
	}
	defer file.Close()

	code := readCode(file, ctxt)

	ctxt.Logger.Println("Pretty printing now.")

	ptr := 0
	for ptr < len(code) {
		switch it := code[ptr]; it {
		case IThalt, ITret:
			fmt.Printf("%04d: %s\n", ptr, it2str[it])
			ptr++
		case ITset:
			fmt.Printf("%04d: %s r%d %d\n",
				ptr, it2str[it], code[ptr+1], code[ptr+2])
			ptr += 3
		case ITmov, ITadd, ITsub, ITmul, ITdiv, ITrem, ITcmp:
			fmt.Printf("%04d: %s r%d r%d\n",
				ptr, it2str[it], code[ptr+1], code[ptr+2])
			ptr += 3
		case ITjmp, ITjeq, ITjne, ITjgt, ITjlt, ITjge, ITjle, ITcall:
			fmt.Printf("%04d: %s %d\n", ptr, it2str[it], code[ptr+1])
			ptr += 2
		case ITshow, ITinc, ITdec, ITpush, ITpop:
			fmt.Printf("%04d: %s r%d\n", ptr, it2str[it], code[ptr+1])
			ptr += 2
		default:
			ctxt.Logger.Fatalf("Unexpected instruction %d.\n", code[ptr])
		}
	}
}
