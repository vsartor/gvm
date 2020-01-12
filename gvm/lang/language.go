package lang

import (
	"errors"
	"fmt"
	"github.com/vsartor/gvm/gvm"
)

// Instructions
const (
	Halt gvm.Code = iota
	Const
	Push
	Pop
	Inc
	Dec
	Mov
	Add
	Sub
	Mul
	Div
	Rem
	Cmp
	Jmp
	Jeq
	Jne
	Jgt
	Jlt
	Jge
	Jle
	Jerr
	Show
	Call
	Ret
	Noop
	Iarg
)

// Mappings between instructions and their string representations
var (
	reprFromIns         map[gvm.Code]string
	instructionFromRepr map[string]gvm.Code
)

// Initialize the mappings
func init() {
	reprFromIns = make(map[gvm.Code]string)
	instructionFromRepr = make(map[string]gvm.Code)

	reprFromIns[Halt] = "halt"
	reprFromIns[Const] = "const"
	reprFromIns[Push] = "push"
	reprFromIns[Pop] = "pop"
	reprFromIns[Inc] = "inc"
	reprFromIns[Dec] = "dec"
	reprFromIns[Mov] = "mov"
	reprFromIns[Add] = "add"
	reprFromIns[Sub] = "sub"
	reprFromIns[Mul] = "mul"
	reprFromIns[Div] = "div"
	reprFromIns[Rem] = "rem"
	reprFromIns[Cmp] = "cmp"
	reprFromIns[Jmp] = "jmp"
	reprFromIns[Jeq] = "jeq"
	reprFromIns[Jne] = "jne"
	reprFromIns[Jgt] = "jgt"
	reprFromIns[Jlt] = "jlt"
	reprFromIns[Jge] = "jge"
	reprFromIns[Jle] = "jle"
	reprFromIns[Jerr] = "jerr"
	reprFromIns[Show] = "show"
	reprFromIns[Call] = "call"
	reprFromIns[Ret] = "ret"
	reprFromIns[Noop] = "noop"
	reprFromIns[Iarg] = "iarg"

	for instruction, repr := range reprFromIns {
		instructionFromRepr[repr] = instruction
	}
}

func ToString(ins gvm.Code) string {
	return reprFromIns[ins]
}

func ParseInstruction(repr string) (gvm.Code, error) {
	if instruction, ok := instructionFromRepr[repr]; ok {
		return instruction, nil
	} else {
		errorMessage := fmt.Sprintf("Unexpected instruction '%s'.", repr)
		return Halt, errors.New(errorMessage)
	}
}
