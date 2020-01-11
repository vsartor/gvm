package gvmlib

import (
	"fmt"
	"os"
)

const (
	regCount      = 16
	stackSize     = 128
	callStackSize = 64
)

func Run(filePath string, ctxt Context) {
	ctxt.Logger.Println("Starting to disassemble.")

	ctxt.Logger.Printf("Opening file '%s' now.\n", filePath)
	file, err := os.Open(filePath)
	if err != nil {
		ctxt.Logger.Printf("Error while opening '%s'.\n", filePath)
		ctxt.Logger.Fatalf("err: %v\n", err.Error())
	}
	defer file.Close()

	code := readCode(file, ctxt)

	ctxt.Logger.Println("Starting execution now.")

	stack := make([]int64, stackSize)
	stackPtr := 0
	callStack := make([]int64, callStackSize)
	callStackPtr := 0
	reg := make([]int64, regCount)
	var cmpFlag int64

	ptr := int64(0)
	for ptr < int64(len(code)) {
		switch code[ptr] {
		case IThalt:
			ptr = int64(len(code))
		case ITset:
			reg[code[ptr+1]] = code[ptr+2]
			ptr += 3
		case ITpush:
			stack[stackPtr] = reg[code[ptr+1]]
			stackPtr++
			ptr += 2
		case ITpop:
			stackPtr--
			reg[code[ptr+1]] = stack[stackPtr]
			ptr += 2
		case ITinc:
			reg[code[ptr+1]]++
			ptr += 2
		case ITdec:
			reg[code[ptr+1]]--
			ptr += 2
		case ITmov:
			reg[code[ptr+2]] = reg[code[ptr+1]]
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
		case ITcmp:
			cmpFlag = reg[code[ptr+2]] - reg[code[ptr+1]]
			ptr += 3
		case ITjmp:
			ptr = code[ptr+1]
		case ITjeq:
			if cmpFlag == 0 {
				ptr = code[ptr+1]
			} else {
				ptr += 2
			}
		case ITjne:
			if cmpFlag != 0 {
				ptr = code[ptr+1]
			} else {
				ptr += 2
			}
		case ITjgt:
			if cmpFlag > 0 {
				ptr = code[ptr+1]
			} else {
				ptr += 2
			}
		case ITjlt:
			if cmpFlag < 0 {
				ptr = code[ptr+1]
			} else {
				ptr += 2
			}
		case ITjge:
			if cmpFlag >= 0 {
				ptr = code[ptr+1]
			} else {
				ptr += 2
			}
		case ITjle:
			if cmpFlag <= 0 {
				ptr = code[ptr+1]
			} else {
				ptr += 2
			}
		case ITshow:
			fmt.Printf("%d\n", reg[code[ptr+1]])
			ptr += 2
		case ITcall:
			callStack[callStackPtr] = ptr + 2
			ptr = code[ptr+1]
			callStackPtr++
		case ITret:
			callStackPtr--
			ptr = callStack[callStackPtr]
		default:
			ctxt.Logger.Fatalf("Unexpected instruction %d.\n", code[ptr])
		}
	}
}
