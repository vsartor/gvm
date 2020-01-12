package vm

import (
	"github.com/vsartor/gvm/gvm"
	"github.com/vsartor/gvm/gvm/compiler"
	"os"
)

func Debug(filePath string, args []string) {
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
		codePosition := vm.codePosition
		instruction := code[codePosition]
		executeStep(instruction, code, &vm)
		disassembleStep(instruction, code, codePosition)
	}
}
