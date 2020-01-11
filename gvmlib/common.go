package gvmlib

import "log"

type Context struct {
	Logger    *log.Logger
	IsVerbose bool
	LineNum   int
}
