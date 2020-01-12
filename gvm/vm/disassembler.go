package vm

import (
	"fmt"
	"github.com/vsartor/gvm/gvm"
	"github.com/vsartor/gvm/gvm/compiler"
	"github.com/vsartor/gvm/gvm/lang"
	"os"
	"path"
)

func disassembleStep(instruction gvm.Code, code []gvm.Code, codePosition int64) int64 {
	switch instruction {
	case lang.Halt, lang.Ret, lang.Noop:
		fmt.Printf("%04d: %s\n", codePosition, lang.ToString(instruction))
		codePosition++
	case lang.Const:
		fmt.Printf("%04d: %s %d r%d\n",
			codePosition, lang.ToString(instruction), code[codePosition+1], code[codePosition+2])
		codePosition += 3
	case lang.Mov, lang.Add, lang.Sub, lang.Mul, lang.Div, lang.Rem, lang.Cmp:
		fmt.Printf("%04d: %s r%d r%d\n",
			codePosition, lang.ToString(instruction), code[codePosition+1], code[codePosition+2])
		codePosition += 3
	case lang.Jmp, lang.Jeq, lang.Jne, lang.Jgt, lang.Jlt, lang.Jge, lang.Jle, lang.Jerr, lang.Call:
		fmt.Printf("%04d: %s %d\n", codePosition, lang.ToString(instruction), code[codePosition+1])
		codePosition += 2
	case lang.Show, lang.Inc, lang.Dec, lang.Push, lang.Pop, lang.Iarg:
		fmt.Printf("%04d: %s r%d\n", codePosition, lang.ToString(instruction), code[codePosition+1])
		codePosition += 2
	default:
		gvm.Logger.Criticalf("Unexpected instruction code %d.\n", instruction)
		os.Exit(1)
	}

	return codePosition
}

func disassemble(code []gvm.Code, ctxt gvm.Context) {
	codePosition := int64(0)
	for codePosition < int64(len(code)) {
		instruction := code[codePosition]
		codePosition = disassembleStep(instruction, code, codePosition)
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

	ctxt := gvm.Context{FileName: path.Base(filePath)}
	disassemble(code, ctxt)
}
