package vm

import (
	"fmt"
	"github.com/vsartor/gvm/gvm"
	"github.com/vsartor/gvm/gvm/compiler"
	"github.com/vsartor/gvm/gvm/lang"
	"os"
)

func Run(filePath string, ctxt gvm.Context) {
	gvm.Logger.Println("Starting to disassemble.")

	gvm.Logger.Printf("Opening file '%s'.\n", filePath)
	file, err := os.Open(filePath)
	if err != nil {
		gvm.Logger.Fatalf("Failed opening '%s': %s\n", filePath, err.Error())
	}
	defer file.Close()

	code := compiler.ReadCode(file)

	gvm.Logger.Println("Starting execution.")

	stack := make([]int64, gvm.StackSize)
	stackPtr := 0
	callStack := make([]int64, gvm.CallStackSize)
	callStackPtr := 0
	reg := make([]int64, gvm.RegisterCount)
	var cmpFlag int64

	codePosition := int64(0)
	for codePosition < int64(len(code)) {
		switch code[codePosition] {
		case lang.Halt:
			// Program needs to stop. Do so by making the loop condition false.
			codePosition = int64(len(code))
		case lang.Set:
			dstRegIdx := code[codePosition+1]
			intConst := int64(code[codePosition+2])
			reg[dstRegIdx] = intConst
			codePosition += 3
		case lang.Push:
			srcRegIdx := code[codePosition+1]
			stack[stackPtr] = reg[srcRegIdx]
			stackPtr++
			codePosition += 2
		case lang.Pop:
			stackPtr--
			dstRegIdx := code[codePosition+1]
			reg[dstRegIdx] = stack[stackPtr]
			codePosition += 2
		case lang.Inc:
			dstRegIdx := code[codePosition+1]
			reg[dstRegIdx]++
			codePosition += 2
		case lang.Dec:
			dstRegIdx := code[codePosition+1]
			reg[dstRegIdx]--
			codePosition += 2
		case lang.Mov:
			srcRegIdx := code[codePosition+1]
			dstRegIdx := code[codePosition+2]
			reg[dstRegIdx] = reg[srcRegIdx]
			codePosition += 3
		case lang.Add:
			srcRegIdx := code[codePosition+1]
			dstRegIdx := code[codePosition+2]
			reg[dstRegIdx] += reg[srcRegIdx]
			codePosition += 3
		case lang.Sub:
			srcRegIdx := code[codePosition+1]
			dstRegIdx := code[codePosition+2]
			reg[dstRegIdx] -= reg[srcRegIdx]
			codePosition += 3
		case lang.Mul:
			srcRegIdx := code[codePosition+1]
			dstRegIdx := code[codePosition+2]
			reg[dstRegIdx] *= reg[srcRegIdx]
			codePosition += 3
		case lang.Div:
			srcRegIdx := code[codePosition+1]
			dstRegIdx := code[codePosition+2]
			reg[dstRegIdx] /= reg[srcRegIdx]
			codePosition += 3
		case lang.Rem:
			srcRegIdx := code[codePosition+1]
			dstRegIdx := code[codePosition+2]
			reg[dstRegIdx] %= reg[srcRegIdx]
			codePosition += 3
		case lang.Cmp:
			srcRegIdx := code[codePosition+1]
			dstRegIdx := code[codePosition+2]
			cmpFlag = reg[dstRegIdx] - reg[srcRegIdx]
			cmpFlag = reg[code[codePosition+2]] - reg[code[codePosition+1]]
			codePosition += 3
		case lang.Jmp:
			codePosition = int64(code[codePosition+1])
		case lang.Jeq:
			if cmpFlag == 0 {
				codePosition = int64(code[codePosition+1])
			} else {
				codePosition += 2
			}
		case lang.Jne:
			if cmpFlag != 0 {
				codePosition = int64(code[codePosition+1])
			} else {
				codePosition += 2
			}
		case lang.Jgt:
			if cmpFlag > 0 {
				codePosition = int64(code[codePosition+1])
			} else {
				codePosition += 2
			}
		case lang.Jlt:
			if cmpFlag < 0 {
				codePosition = int64(code[codePosition+1])
			} else {
				codePosition += 2
			}
		case lang.Jge:
			if cmpFlag >= 0 {
				codePosition = int64(code[codePosition+1])
			} else {
				codePosition += 2
			}
		case lang.Jle:
			if cmpFlag <= 0 {
				codePosition = int64(code[codePosition+1])
			} else {
				codePosition += 2
			}
		case lang.Show:
			srcRegIdx := code[codePosition+1]
			fmt.Printf("%d\n", reg[srcRegIdx])
			codePosition += 2
		case lang.Call:
			if callStackPtr == len(callStack) {
				gvm.Logger.Fatalf("Call stack overflow.")
			}
			callStack[callStackPtr] = codePosition + 2
			codePosition = int64(code[codePosition+1])
			callStackPtr++
		case lang.Ret:
			callStackPtr--
			codePosition = callStack[callStackPtr]
		default:
			gvm.Logger.Fatalf("Unexpected instruction code %d.\n", code[codePosition])
		}
	}
}