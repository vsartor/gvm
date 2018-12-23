package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Fprint(os.Stderr, "gvm: Expected at least one argument.\n")
		fmt.Fprint(os.Stderr, "gvm: Call 'gvm h' for usage.\n")
		os.Exit(1)
	}

	// Check early if we're actually doing logging
	logger := log.New(ioutil.Discard, "", log.Ldate|log.Ltime|log.Lshortfile)
	debugOffset := 0
	if args[0] == "d" || args[0] == "debug" {
		debugOffset++
		logger.SetOutput(os.Stdout)
	}
	if debugOffset == 1 && len(args) == 1 {
		fmt.Fprintln(os.Stderr, "gvm: Only received a debug flag.")
		os.Exit(1)
	}

	// Dispatch into the correct routine
	switch runMode := args[0+debugOffset]; runMode {
	case "c":
		logger.Println("Compilation mode has been set.")
		if len(args) != 3+debugOffset {
			fmt.Fprint(os.Stderr, "gvm: Expected two files after 'c'.\n")
			os.Exit(1)
		}
		compile(args[1+debugOffset], args[2+debugOffset], logger)

	case "r":
		logger.Println("Execution mode has been set.")
		if len(args) != 2+debugOffset {
			fmt.Fprint(os.Stderr, "gvm: Expected a file after 'r'.\n")
			os.Exit(1)
		}
		run(args[1+debugOffset], logger)

	case "d":
		logger.Println("Disassembly mode has been set.")
		if len(args) != 2+debugOffset {
			fmt.Fprint(os.Stderr, "gvm: Expected a file after 'o'.\n")
			os.Exit(1)
		}
		disassemble(args[1+debugOffset], logger)

	case "cr":
		logger.Println("Compilation and running mode has been set.")
		if len(args) != 3+debugOffset {
			fmt.Fprint(os.Stderr, "gvm: Expected two files after 'cr'.\n")
			os.Exit(1)
		}
		logger.Println("Starting compilation mode.")
		compile(args[1+debugOffset], args[2+debugOffset], logger)
		logger.Println("Starting running mode.")
		run(args[2+debugOffset], logger)

	case "h":
		fmt.Println("gvm [d|debug] <run mode> [file] [output]")
		fmt.Println("Run mode options:")
		fmt.Println("- h:  Shows usage.")
		fmt.Println("- c:  Compiles a file.")
		fmt.Println("- r:  Runs a compiled file.")
		fmt.Println("- d:  Disassemble and pretty print a compiled file.")
		fmt.Println("- cr: Compiles a file and then runs it.")

	default:
		fmt.Fprintf(os.Stderr,
			"gvm: Unknown run mode '%s' was attempted to be set.\n",
			runMode)
		os.Exit(1)
	}
}
