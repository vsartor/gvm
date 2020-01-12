package vm

import (
	"fmt"
	"github.com/vsartor/gvm/gvm"
	"github.com/vsartor/gvm/gvm/compiler"
	"github.com/vsartor/gvm/gvm/lang"
	"os"
)

func disassemble(code []gvm.Code, ctxt gvm.Context) {
	codePosition := 0
	for codePosition < len(code) {
		switch instruction := code[codePosition]; instruction {
		case lang.Halt, lang.Ret:
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
		case lang.Jmp, lang.Jeq, lang.Jne, lang.Jgt, lang.Jlt, lang.Jge, lang.Jle, lang.Call:
			fmt.Printf("%04d: %s %d\n", codePosition, lang.ToString(instruction), code[codePosition+1])
			codePosition += 2
		case lang.Show, lang.Inc, lang.Dec, lang.Push, lang.Pop:
			fmt.Printf("%04d: %s r%d\n", codePosition, lang.ToString(instruction), code[codePosition+1])
			codePosition += 2
		default:
			gvm.Logger.Fatalf("Unexpected instruction code %d.\n", instruction)
		}
	}
}

func Disassemble(filePath string, ctxt gvm.Context) {
	gvm.Logger.Println("Starting to disassemble.")

	gvm.Logger.Printf("Opening file '%s'.\n", filePath)
	file, err := os.Open(filePath)
	if err != nil {
		gvm.Logger.Fatalf("Failed opening '%s': %s\n", filePath, err.Error())
	}
	defer file.Close()

	code := compiler.ReadCode(file)

	disassemble(code, ctxt)
}
