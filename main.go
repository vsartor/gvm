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
	hasLogging := false
	if args[0] == "l" {
		gvm.Logger.SetLogLevel(loggo.INFO)
		hasLogging = true
	} else if args[0] == "L" {
		gvm.Logger.SetLogLevel(loggo.DEBUG)
		hasLogging = true
	}

	// Error out in case only a debug flag was passed
	if hasLogging && len(args) == 1 {
		appLogger.Fatalln("Only received a debug flag.")
	}

	// Remove logging argument from args
	if hasLogging {
		args = args[1:]
	}

	// Dispatch into the correct routine
	switch runMode := args[0]; runMode {
	case "c", "compile":
		if len(args) != 3 {
			appLogger.Fatalln("Expected two files after 'compile': <source_path>, <object_path>")
		}
		compiler.Compile(args[1], args[2], ctxt)

	case "r", "run":
		if len(args) != 2 {
			appLogger.Fatalln("Expected one file after 'run': <object_path>")
		}
		vm.Run(args[1], ctxt)

	case "d", "disassemble":
		if len(args) != 2 {
			appLogger.Fatalln("Expected one file after 'disassemble': <object_path>")
		}
		vm.Disassemble(args[1], ctxt)

	case "cr":
		if len(args) != 3 {
			appLogger.Fatalln("Expected two files after 'cr': <source_path>, <object_path>")
		}
		compiler.Compile(args[1], args[2], ctxt)
		vm.Run(args[2], ctxt)

	case "cd":
		if len(args) != 3 {
			appLogger.Fatalln("Expected two files after 'cd': <source_path>, <object_path>")
		}
		compiler.Compile(args[1], args[2], ctxt)
		vm.Disassemble(args[2], ctxt)

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
