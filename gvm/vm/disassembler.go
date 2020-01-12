package vm

import (
	"fmt"
	"github.com/vsartor/gvm/gvm"
	"github.com/vsartor/gvm/gvm/compiler"
	"github.com/vsartor/gvm/gvm/lang"
	"os"
)

func disassembleStep(vm virtualMachine, code []gvm.Code) virtualMachine {
	// Because performance is not a concern for this step and this disassembly step will
	// be used in  conjuction with other functions for debugging purposes, let's make our
	// life easier  by not mutating the virtualMachine, but instead returning a mutated
	// copy that the caller can opt to use or not.

	switch instruction := code[vm.codePosition]; instruction {
	case lang.Halt, lang.Ret, lang.Noop:
		fmt.Printf("%04d: %s\n", vm.codePosition, lang.ToString(instruction))
		vm.codePosition++
	case lang.Const:
		fmt.Printf("%04d: %s %d r%d\n",
			vm.codePosition, lang.ToString(instruction), code[vm.codePosition+1], code[vm.codePosition+2])
		vm.codePosition += 3
	case lang.Mov, lang.Add, lang.Sub, lang.Mul, lang.Div, lang.Rem, lang.Cmp:
		fmt.Printf("%04d: %s r%d r%d\n",
			vm.codePosition, lang.ToString(instruction), code[vm.codePosition+1], code[vm.codePosition+2])
		vm.codePosition += 3
	case lang.Jmp, lang.Jeq, lang.Jne, lang.Jgt, lang.Jlt, lang.Jge, lang.Jle, lang.Jerr, lang.Call:
		fmt.Printf("%04d: %s %d\n", vm.codePosition, lang.ToString(instruction), code[vm.codePosition+1])
		vm.codePosition += 2
	case lang.Show, lang.Inc, lang.Dec, lang.Push, lang.Pop, lang.Iarg:
		fmt.Printf("%04d: %s r%d\n", vm.codePosition, lang.ToString(instruction), code[vm.codePosition+1])
		vm.codePosition += 2
	default:
		gvm.Logger.Criticalf("Unexpected instruction code %d.\n", instruction)
		os.Exit(1)
	}

	return vm
}

func disassemble(code []gvm.Code) {
	vm := virtualMachine{}
	for vm.codePosition < int64(len(code)) {
		vm = disassembleStep(vm, code)
	}
}

func Disassemble(filePath string) {
	gvm.Logger.Infof("Starting to disassemble.\n")

	gvm.Logger.Infof("Opening file '%s'.\n", filePath)
	file, err := os.Open(filePath)
	if err != nil {
		gvm.Logger.Criticalf("Failed opening '%s': %s\n", filePath, err.Error())
		os.Exit(1)
	}
	defer file.Close()

	code := compiler.ReadCode(file)

	disassemble(code)
}
