package main

import (
	"fmt"
	"github.com/juju/loggo"
	"github.com/vsartor/gvm/gvm"
	"github.com/vsartor/gvm/gvm/compiler"
	"github.com/vsartor/gvm/gvm/vm"
	"os"
	"path"
	"strconv"
	"time"
)

var appLogger loggo.Logger

func init() {
	appLogger = loggo.GetLogger("gvm")
	appLogger.SetLogLevel(loggo.INFO)
}

func currentTimestamp() string {
	nanoseconds := time.Now().UnixNano()
	return strconv.FormatInt(nanoseconds, 10)
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		appLogger.Criticalf("Expected at least one argument. Call 'gvm help' for usage.\n")
		os.Exit(1)
	}

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
		appLogger.Criticalf("Only received a debug flag.\n")
		os.Exit(1)
	}

	// Remove logging argument from args
	if hasLogging {
		args = args[1:]
	}

	// Dispatch into the correct routine
	switch runMode := args[0]; runMode {
	case "c", "compile":
		if len(args) != 3 {
			appLogger.Criticalf("Expected two files after 'compile': <source_path>, <object_path>\n")
			os.Exit(1)
		}
		compiler.Compile(args[1], args[2])

	case "r", "run":
		if len(args) < 2 {
			appLogger.Criticalf("Expected one file after 'run': <object_path>\n")
			os.Exit(1)
		}
		vm.Run(args[1], args[2:])

	case "d", "disassemble":
		if len(args) != 2 {
			appLogger.Criticalf("Expected one file after 'disassemble': <object_path>\n")
			os.Exit(1)
		}
		vm.Disassemble(args[1])

	case "cr":
		// For composite commands, if output directory is not given write to tmpdir
		if len(args) == 2 {
			args = append(args, path.Join(os.TempDir(), currentTimestamp()))
		}

		if len(args) < 3 {
			appLogger.Criticalf("Expected two files after 'cr': <source_path>, <object_path>\n")
			os.Exit(1)
		}
		compiler.Compile(args[1], args[2])
		vm.Run(args[2], args[3:])

	case "cd":
		// For composite commands, if output directory is not given write to tmpdir
		if len(args) == 2 {
			args = append(args, path.Join(os.TempDir(), currentTimestamp()))
		}

		if len(args) != 3 {
			appLogger.Criticalf("Expected two files after 'cd': <source_path>, <object_path>\n")
			os.Exit(1)
		}
		compiler.Compile(args[1], args[2])
		vm.Disassemble(args[2])

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
		appLogger.Criticalf("Unknown command '%s'.\n", runMode)
		os.Exit(1)
	}
}
