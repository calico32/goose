package lib

import (
	"io"

	"github.com/calico32/goose/token"
)

type Interpreter interface {
	Stdin() io.Reader
	Stdout() io.Writer
	Stderr() io.Writer
	Modules() map[string]*Module
	Fset() *token.FileSet
	ExecutionStack() []*Module
	Global() *Scope
	GooseRoot() string

	CurrentModule() *Module

	Run() (exitCode int, err error)

	Throw(format string, args ...interface{})
}
