package vm

import (
	"bufio"
	"fmt"
	"github.com/vsartor/gvm/gvm"
	"github.com/vsartor/gvm/gvm/compiler"
	"os"
	"strconv"
	"strings"
)

type debugContext struct {
	breakPoints      []int64
	currentDirective string
	reader           *bufio.Reader
}

func (ctxt *debugContext) isBreakpoint(codePosition int64) bool {
	for _, bp := range ctxt.breakPoints {
		if codePosition == bp {
			return true
		}
	}
	return false
}

func (ctxt *debugContext) isNotBreakpoint(codePosition int64) bool {
	return !ctxt.isBreakpoint(codePosition)
}

func debugStep(vm *virtualMachine, code []gvm.Code, ctxt *debugContext) {
	// Show current position
	disassembleStep(*vm, code)

	// Decide what to do
	promptInput := false
	switch ctxt.currentDirective {
	case "", "n":
		// We were either told to walk a single step or nothing, either way
		// prompt for another input
		promptInput = true
	case "c":
		// We were told to continue until we reach a breakpoint, so check if
		// we reached a breakpoint.
		if ctxt.isBreakpoint(vm.codePosition) {
			promptInput = true
		}
	default:
		panic("This path should be impossible.")
	}

	// Get user input if required
	if promptInput {
		fmt.Printf("gvm.Debugger: ")
		userInput, err := ctxt.reader.ReadString('\n')
		if err != nil {
			gvm.Logger.Errorf("Problem reading input.\n")
			return
		}

		// Validate input
		tokens := strings.Fields(userInput)

		if len(tokens) > 0 {
			switch command := tokens[0]; command {
			case "n", "c":
				ctxt.currentDirective = command

			case "x", "exit":
				vm.codePosition = int64(len(code))
				return

			case "bp":
				if len(tokens) != 2 {
					gvm.Logger.Errorf("`bp` requires one argument.\n")
					return
				}

				value, parseErr := strconv.ParseInt(tokens[1], 10, 64)
				if parseErr != nil {
					gvm.Logger.Errorf("Could not process integer '%s'.\n", value)
					return
				}

				if ctxt.isNotBreakpoint(value) {
					ctxt.breakPoints = append(ctxt.breakPoints, value)
				}
				return

			case "p":
				if len(tokens) != 2 {
					gvm.Logger.Errorf("`p` requires one argument.\n")
					return
				}

				switch tokens[1] {
				case "stack":
					fmt.Printf("%v\n", vm.stack[:vm.stackPtr])
				case "reg":
					fmt.Printf("%v\n", vm.reg)
				case "code":
					disassemble(code)
				default:
					gvm.Logger.Errorf("Unknown argument '%s'.\n", command)
				}
				return

			default:
				gvm.Logger.Errorf("Unexpected command '%s'.\n", command)
				return
			}
		}
	}

	// Actual execution step
	executeStep(vm, code)
}

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

	ctxt := debugContext{reader: bufio.NewReader(os.Stdin)}
	for vm.codePosition < int64(len(code)) {
		debugStep(&vm, code, &ctxt)
	}
}
