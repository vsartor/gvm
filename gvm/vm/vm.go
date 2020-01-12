package vm

import (
	"fmt"
	"github.com/vsartor/gvm/gvm"
	"github.com/vsartor/gvm/gvm/compiler"
	"github.com/vsartor/gvm/gvm/lang"
	"os"
	"strconv"
)

type virtualMachine struct {
	stack        []int64
	stackPtr     int64
	callStack    []int64
	callStackPtr int64
	reg          []int64
	codePosition int64
	cmpFlag      int64
	errFlag      int64
	args         []string
}

func executeStep(vm *virtualMachine, code []gvm.Code) {
	switch instruction := code[vm.codePosition]; instruction {
	case lang.Halt:
		// Program needs to stop. Do so by making the loop condition false.
		vm.codePosition = int64(len(code))
	case lang.Const:
		dstRegIdx := code[vm.codePosition+2]
		intConst := int64(code[vm.codePosition+1])
		vm.reg[dstRegIdx] = intConst
		vm.codePosition += 3
	case lang.Push:
		srcRegIdx := code[vm.codePosition+1]
		vm.stack[vm.stackPtr] = vm.reg[srcRegIdx]
		vm.stackPtr++
		vm.codePosition += 2
	case lang.Pop:
		if vm.stackPtr == 0 {
			gvm.Logger.Criticalf("Stack underflow.\n")
			os.Exit(1)
		}
		vm.stackPtr--
		dstRegIdx := code[vm.codePosition+1]
		vm.reg[dstRegIdx] = vm.stack[vm.stackPtr]
		vm.codePosition += 2
	case lang.Inc:
		dstRegIdx := code[vm.codePosition+1]
		vm.reg[dstRegIdx]++
		vm.codePosition += 2
	case lang.Dec:
		dstRegIdx := code[vm.codePosition+1]
		vm.reg[dstRegIdx]--
		vm.codePosition += 2
	case lang.Mov:
		srcRegIdx := code[vm.codePosition+1]
		dstRegIdx := code[vm.codePosition+2]
		vm.reg[dstRegIdx] = vm.reg[srcRegIdx]
		vm.codePosition += 3
	case lang.Add:
		srcRegIdx := code[vm.codePosition+1]
		dstRegIdx := code[vm.codePosition+2]
		vm.reg[dstRegIdx] += vm.reg[srcRegIdx]
		vm.codePosition += 3
	case lang.Sub:
		srcRegIdx := code[vm.codePosition+1]
		dstRegIdx := code[vm.codePosition+2]
		vm.reg[dstRegIdx] -= vm.reg[srcRegIdx]
		vm.codePosition += 3
	case lang.Mul:
		srcRegIdx := code[vm.codePosition+1]
		dstRegIdx := code[vm.codePosition+2]
		vm.reg[dstRegIdx] *= vm.reg[srcRegIdx]
		vm.codePosition += 3
	case lang.Div:
		srcRegIdx := code[vm.codePosition+1]
		dstRegIdx := code[vm.codePosition+2]
		vm.reg[dstRegIdx] /= vm.reg[srcRegIdx]
		vm.codePosition += 3
	case lang.Rem:
		srcRegIdx := code[vm.codePosition+1]
		dstRegIdx := code[vm.codePosition+2]
		vm.reg[dstRegIdx] %= vm.reg[srcRegIdx]
		vm.codePosition += 3
	case lang.Cmp:
		srcRegIdx := code[vm.codePosition+1]
		dstRegIdx := code[vm.codePosition+2]
		vm.cmpFlag = vm.reg[dstRegIdx] - vm.reg[srcRegIdx]
		vm.cmpFlag = vm.reg[code[vm.codePosition+2]] - vm.reg[code[vm.codePosition+1]]
		vm.codePosition += 3
	case lang.Jmp:
		vm.codePosition = int64(code[vm.codePosition+1])
	case lang.Jeq:
		if vm.cmpFlag == 0 {
			vm.codePosition = int64(code[vm.codePosition+1])
		} else {
			vm.codePosition += 2
		}
	case lang.Jne:
		if vm.cmpFlag != 0 {
			vm.codePosition = int64(code[vm.codePosition+1])
		} else {
			vm.codePosition += 2
		}
	case lang.Jgt:
		if vm.cmpFlag > 0 {
			vm.codePosition = int64(code[vm.codePosition+1])
		} else {
			vm.codePosition += 2
		}
	case lang.Jlt:
		if vm.cmpFlag < 0 {
			vm.codePosition = int64(code[vm.codePosition+1])
		} else {
			vm.codePosition += 2
		}
	case lang.Jge:
		if vm.cmpFlag >= 0 {
			vm.codePosition = int64(code[vm.codePosition+1])
		} else {
			vm.codePosition += 2
		}
	case lang.Jle:
		if vm.cmpFlag <= 0 {
			vm.codePosition = int64(code[vm.codePosition+1])
		} else {
			vm.codePosition += 2
		}
	case lang.Jerr:
		if vm.errFlag != 0 {
			vm.errFlag = 0
			vm.codePosition = int64(code[vm.codePosition+1])
		} else {
			vm.codePosition += 2
		}
	case lang.Show:
		srcRegIdx := code[vm.codePosition+1]
		fmt.Printf("%d\n", vm.reg[srcRegIdx])
		vm.codePosition += 2
	case lang.Call:
		if vm.callStackPtr == int64(len(vm.callStack)) {
			gvm.Logger.Criticalf("Call stack overflow.\n")
			os.Exit(1)
		}
		vm.callStack[vm.callStackPtr] = vm.codePosition + 2
		vm.codePosition = int64(code[vm.codePosition+1])
		vm.callStackPtr++
	case lang.Ret:
		if vm.callStackPtr == 0 {
			gvm.Logger.Criticalf("Call stack underflow.\n")
			os.Exit(1)
		}
		vm.callStackPtr--
		vm.codePosition = vm.callStack[vm.callStackPtr]
	case lang.Noop:
		vm.codePosition++
	case lang.Iarg:
		srcRegIdx := code[vm.codePosition+1]
		argIdx := vm.reg[srcRegIdx]
		if argIdx < 0 || argIdx >= int64(len(vm.args)) {
			vm.errFlag = 1
		} else {
			value, err := strconv.ParseInt(vm.args[argIdx], 10, 64)
			if err != nil {
				vm.errFlag = 1
			} else {
				vm.stack[vm.stackPtr] = value
				vm.stackPtr++
			}
		}
		vm.codePosition += 2
	default:
		gvm.Logger.Criticalf("Unexpected instruction code %d.\n", code[vm.codePosition])
		os.Exit(1)
	}
}

func Execute(filePath string, args []string) {
	gvm.Logger.Infof("Starting to disassemble.\n")

	gvm.Logger.Infof("Opening file '%s'.\n", filePath)
	file, err := os.Open(filePath)
	if err != nil {
		gvm.Logger.Criticalf("Failed opening '%s': %s\n", filePath, err.Error())
		os.Exit(1)
	}
	defer file.Close()

	code := compiler.ReadCode(file)

	gvm.Logger.Infof("Starting execution.\n")

	vm := virtualMachine{
		stack:        make([]int64, gvm.StackSize),
		stackPtr:     0,
		callStack:    make([]int64, gvm.CallStackSize),
		callStackPtr: 0,
		reg:          make([]int64, gvm.RegisterCount),
		codePosition: 0,
		cmpFlag:      0,
		errFlag:      0,
		args:         args,
	}

	for vm.codePosition < int64(len(code)) {
		executeStep(&vm, code)
	}
}
