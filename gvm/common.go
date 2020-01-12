package gvm

import (
	"io/ioutil"
	"log"
)

type Code int64

type Context struct {
	IsVerbose bool
	LineNum   int
}

var Logger *log.Logger

func init() {
	// Initialize the logger with no output (ioutil.Discard).
	// The cli entry point should set a different output in case an
	// appropriate flag is passed during program invokation.
	Logger = log.New(ioutil.Discard, "gvm: ", log.Lshortfile)
}
