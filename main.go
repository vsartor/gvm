package main

import (
	"fmt"
	"github.com/vsartor/gvm/gvmlib"
	"io/ioutil"
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

	var ctxt gvmlib.Context

	// Check early if we're actually doing logging
	ctxt.Logger = log.New(ioutil.Discard, "", log.Ltime|log.Lshortfile)
	debugOffset := 0
	if args[0] == "l" {
		debugOffset++
		ctxt.Logger.SetOutput(os.Stderr)
	} else if args[0] == "L" {
		debugOffset++
		ctxt.Logger.SetOutput(os.Stdout)
		ctxt.IsVerbose = true
	}
	if debugOffset == 1 && len(args) == 1 {
		appLogger.Fatalln("Only received a debug flag.")
	}

	// Dispatch into the correct routine
	switch runMode := args[0+debugOffset]; runMode {
	case "c", "compile":
		ctxt.Logger.Println("Compilation mode has been set.")
		if len(args) != 3+debugOffset {
			appLogger.Fatalln("Expected two files after 'compile': <source_path>, <object_path>")
		}
		gvmlib.Compile(args[1+debugOffset], args[2+debugOffset], ctxt)

	case "r", "run":
		ctxt.Logger.Println("Execution mode has been set.")
		if len(args) != 2+debugOffset {
			appLogger.Fatalln("Expected one file after 'run': <object_path>")
		}
		gvmlib.Run(args[1+debugOffset], ctxt)

	case "d", "disassemble":
		ctxt.Logger.Println("Disassembly mode has been set.")
		if len(args) != 2+debugOffset {
			appLogger.Fatalln("Expected one file after 'disassemble': <object_path>")
		}
		gvmlib.Disassemble(args[1+debugOffset], ctxt)

	case "cr":
		ctxt.Logger.Println("Compilation and running mode has been set.")
		if len(args) != 3+debugOffset {
			appLogger.Fatalln("Expected two files after 'cr': <source_path>, <object_path>")
		}
		ctxt.Logger.Println("Starting compilation mode.")
		gvmlib.Compile(args[1+debugOffset], args[2+debugOffset], ctxt)
		ctxt.Logger.Println("Starting running mode.")
		gvmlib.Run(args[2+debugOffset], ctxt)

	case "cd":
		ctxt.Logger.Println("Compilation and disassembly mode has been set.")
		if len(args) != 3+debugOffset {
			appLogger.Fatalln("Expected two files after 'cd': <source_path>, <object_path>")
		}
		ctxt.Logger.Println("Starting compilation mode.")
		gvmlib.Compile(args[1+debugOffset], args[2+debugOffset], ctxt)
		ctxt.Logger.Println("Starting disassembly mode.")
		gvmlib.Disassemble(args[2+debugOffset], ctxt)

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
