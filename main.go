package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type context struct {
	l       *log.Logger
	verbose bool
	lineNum int
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Fprint(os.Stderr, "gvm: Expected at least one argument.\n")
		fmt.Fprint(os.Stderr, "gvm: Call 'gvm h' for usage.\n")
		os.Exit(1)
	}

	var ctxt context

	// Check early if we're actually doing logging
	ctxt.l = log.New(ioutil.Discard, "", log.Ltime|log.Lshortfile)
	debugOffset := 0
	if args[0] == "d" {
		debugOffset++
		ctxt.l.SetOutput(os.Stderr)
	} else if args[0] == "D" {
		debugOffset++
		ctxt.l.SetOutput(os.Stdout)
		ctxt.verbose = true
	}
	if debugOffset == 1 && len(args) == 1 {
		fmt.Fprintln(os.Stderr, "gvm: Only received a debug flag.")
		os.Exit(1)
	}

	// Dispatch into the correct routine
	switch runMode := args[0+debugOffset]; runMode {
	case "c":
		ctxt.l.Println("Compilation mode has been set.")
		if len(args) != 3+debugOffset {
			fmt.Fprint(os.Stderr, "gvm: Expected two files after 'c'.\n")
			os.Exit(1)
		}
		compile(args[1+debugOffset], args[2+debugOffset], ctxt)

	case "r":
		ctxt.l.Println("Execution mode has been set.")
		if len(args) != 2+debugOffset {
			fmt.Fprint(os.Stderr, "gvm: Expected a file after 'r'.\n")
			os.Exit(1)
		}
		run(args[1+debugOffset], ctxt)

	case "d":
		ctxt.l.Println("Disassembly mode has been set.")
		if len(args) != 2+debugOffset {
			fmt.Fprint(os.Stderr, "gvm: Expected a file after 'o'.\n")
			os.Exit(1)
		}
		disassemble(args[1+debugOffset], ctxt)

	case "cr":
		ctxt.l.Println("Compilation and running mode has been set.")
		if len(args) != 3+debugOffset {
			fmt.Fprint(os.Stderr, "gvm: Expected two files after 'cr'.\n")
			os.Exit(1)
		}
		ctxt.l.Println("Starting compilation mode.")
		compile(args[1+debugOffset], args[2+debugOffset], ctxt)
		ctxt.l.Println("Starting running mode.")
		run(args[2+debugOffset], ctxt)

	case "cd":
		ctxt.l.Println("Compilation and disassembly mode has been set.")
		if len(args) != 3+debugOffset {
			fmt.Fprint(os.Stderr, "gvm: Expected two files after 'cd'.\n")
			os.Exit(1)
		}
		ctxt.l.Println("Starting compilation mode.")
		compile(args[1+debugOffset], args[2+debugOffset], ctxt)
		ctxt.l.Println("Starting disassembly mode.")
		disassemble(args[2+debugOffset], ctxt)

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
