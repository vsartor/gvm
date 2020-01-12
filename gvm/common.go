package gvm

import (
	"github.com/juju/loggo"
	"github.com/juju/loggo/loggocolor"
	"os"
)

type Code int64

type Context struct {
	FileName string
	LineNum  int
}

var Logger loggo.Logger

func init() {
	_, err := loggo.ReplaceDefaultWriter(loggocolor.NewWriter(os.Stdout))
	if err != nil {
		// If we're failing during init, just panic
		panic(err.Error())
	}

	// Initialize the logger with no output (ioutil.Discard).
	// The cli entry point should set a different output in case an
	// appropriate flag is passed during program invokation.
	Logger = loggo.GetLogger("gvm")
	Logger.SetLogLevel(loggo.ERROR)
}
