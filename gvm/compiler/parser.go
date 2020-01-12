package compiler

import (
	"github.com/vsartor/gvm/gvm"
	"github.com/vsartor/gvm/gvm/lang"
	"os"
	"strconv"
)

func parseInstruction(repr string, ctxt gvm.Context) gvm.Code {
	instruction, err := lang.ParseInstruction(repr)
	if err != nil {
		gvm.Logger.Criticalf("%s.%d: Parsing instruction: %s\n",
			ctxt.FileName, ctxt.LineNum, err.Error())
		os.Exit(1)
	}

	return instruction
}

func parseRegister(repr string, ctxt gvm.Context) gvm.Code {
	if repr[0] != 'r' {
		gvm.Logger.Criticalf("%s.%d: Parsing register: Expected 'r', got '%c'.\n",
			ctxt.FileName, ctxt.LineNum, repr[0])
		os.Exit(1)
	}

	reg, err := strconv.ParseInt(repr[1:], 10, 64)
	if err != nil {
		gvm.Logger.Criticalf("%s.%d: Parsing register: Expected integer but got '%s'.\n",
			ctxt.FileName, ctxt.LineNum, repr[1:])
		os.Exit(1)
	}

	return gvm.Code(reg)
}

func parseInt(repr string, ctxt gvm.Context) gvm.Code {
	val, err := strconv.ParseInt(repr, 10, 64)
	if err != nil {
		gvm.Logger.Criticalf("%s.%d: Parsing integer: Expected integer but got '%s'.\n",
			ctxt.FileName, ctxt.LineNum, repr[1:])
		os.Exit(1)
	}

	return gvm.Code(val)
}
