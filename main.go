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
		fmt.Fprint(os.Stderr, "Expected at least one argument.\n")
		fmt.Fprint(os.Stderr, "Call 'gvm h' for usage.\n")
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
		fmt.Fprintln(os.Stderr, "Only received a debug flag.")
		os.Exit(1)
	}

	// Dispatch into the correct routine
	switch runMode := args[0+debugOffset]; runMode {
	case "c":
		logger.Print("Compilation mode has been set.")
		logger.Fatal("Compilation mode not yet implemented.")
	case "r":
		logger.Print("Execution mode has been set.")
		logger.Fatal("Execution mode not yet implemented.")
	case "h":
		fmt.Println("gvm [d|debug] <run mode> [file] [output]")
		fmt.Println("Run mode options:")
		fmt.Println("- h: Shows usage.")
		fmt.Println("- c: Compiles a file.")
		fmt.Println("- r: Runs a compiled file.")
	default:
		fmt.Printf("Unknown run mode '%s' was attempted to be set.\n", runMode)
		os.Exit(1)
	}
}
