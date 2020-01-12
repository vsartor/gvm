package main

import (
	"fmt"
	"github.com/juju/loggo"
	"github.com/vsartor/gvm/gvm"
	"github.com/vsartor/gvm/gvm/compiler"
	"github.com/vsartor/gvm/gvm/vm"
	"log"
	"os"
)

var appLogger *log.Logger

func init() {
	appLogger = log.New(os.Stderr, "gvm: ", 0)
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		appLogger.Fatalln("Expected at least one argument. Call 'gvm help' for usage.")
	}

	var ctxt gvm.Context

	// Check early if we're actually doing logging
	debugOffset := 0
	if args[0] == "l" {
		debugOffset++
		gvm.Logger.SetLogLevel(loggo.INFO)
	} else if args[0] == "L" {
		debugOffset++
		gvm.Logger.SetLogLevel(loggo.DEBUG)
	}

	// Error out in case only a debug flag was passed
	if debugOffset == 1 && len(args) == 1 {
		appLogger.Fatalln("Only received a debug flag.")
	}

	// Dispatch into the correct routine
	switch runMode := args[0+debugOffset]; runMode {
	case "c", "compile":
		if len(args) != 3+debugOffset {
			appLogger.Fatalln("Expected two files after 'compile': <source_path>, <object_path>")
		}
		compiler.Compile(args[1+debugOffset], args[2+debugOffset], ctxt)

	case "r", "run":
		if len(args) != 2+debugOffset {
			appLogger.Fatalln("Expected one file after 'run': <object_path>")
		}
		vm.Run(args[1+debugOffset], ctxt)

	case "d", "disassemble":
		if len(args) != 2+debugOffset {
			appLogger.Fatalln("Expected one file after 'disassemble': <object_path>")
		}
		vm.Disassemble(args[1+debugOffset], ctxt)

	case "cr":
		if len(args) != 3+debugOffset {
			appLogger.Fatalln("Expected two files after 'cr': <source_path>, <object_path>")
		}
		compiler.Compile(args[1+debugOffset], args[2+debugOffset], ctxt)
		vm.Run(args[2+debugOffset], ctxt)

	case "cd":
		if len(args) != 3+debugOffset {
			appLogger.Fatalln("Expected two files after 'cd': <source_path>, <object_path>")
		}
		compiler.Compile(args[1+debugOffset], args[2+debugOffset], ctxt)
		vm.Disassemble(args[2+debugOffset], ctxt)

	case "h", "help":
		fmt.Println("gvm [logging flag] <command> [input file] [output file]")
		fmt.Println("Available commands:")
		fmt.Println("  help (h)           Shows this message.")
		fmt.Println("  compile (c)        Compiles a file.")
		fmt.Println("  run (r)            Runs a compiled file.")
		fmt.Println("  disassemble (d)    Disassembles and pretty prints a compiled file.")
		fmt.Println("  cd                 Compiles, disassembles and pretty prints a compiled file.")
		fmt.Println("  cr                 Compiles a file and then runs it.")
		fmt.Println("Available logging flags:")
		fmt.Println("  l                  Basic logging.")
		fmt.Println("  L                  Verbose logging.")

	default:
		appLogger.Fatalf("Unknown command '%s'.\n", runMode)
	}
}
