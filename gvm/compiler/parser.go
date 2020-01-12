package compiler

import (
	"github.com/vsartor/gvm/gvm"
	"github.com/vsartor/gvm/gvm/lang"
	"strconv"
)

func parseInstruction(repr string, ctxt gvm.Context) gvm.Code {
	instruction, err := lang.ParseInstruction(repr)
	if err != nil {
		gvm.Logger.Fatalf("l%d: Parsing instruction: %s\n", ctxt.LineNum, err.Error())
	}

	return instruction
}

func parseRegister(repr string, ctxt gvm.Context) gvm.Code {
	if repr[0] != 'r' {
		gvm.Logger.Fatalf("l%d: Parsing register: Expected 'r', got '%c'.\n", ctxt.LineNum, repr[0])
	}

	reg, err := strconv.ParseInt(repr[1:], 10, 64)
	if err != nil {
		gvm.Logger.Fatalf("l%d: Parsing register: Expected integer but got '%s'.\n", ctxt.LineNum, repr[1:])
	}

	return gvm.Code(reg)
}

func parseInt(repr string, ctxt gvm.Context) gvm.Code {
	val, err := strconv.ParseInt(repr, 10, 64)
	if err != nil {
		gvm.Logger.Fatalf("l%d: Parsing integer: Expected integer but got '%s'.\n", ctxt.LineNum, repr[1:])
	}

	return gvm.Code(val)
}
