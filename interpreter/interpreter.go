package interpreter

import (
	"fmt"
	"io"
	"os"

	"github.com/calico32/goose/token"
)

type interp struct {
	fset           *token.FileSet
	executionStack []*Module
	modules        map[string]*Module
	global         *Scope
	stdout         io.Writer
	stderr         io.Writer
	gooseRoot      string

	// internal state
	trace    bool
	indent   int
	lastPos  token.Pos
	posStack []token.Pos
}

func (i *interp) currentModule() *Module {
	if len(i.executionStack) == 0 {
		if len(i.modules) == 1 {
			for _, m := range i.modules {
				i.executionStack = append(i.executionStack, m)
				return m
			}
		}
		i.throw("no current module")
	}

	return i.executionStack[len(i.executionStack)-1]
}

const PANIC int = 128

type gooseExit struct{ code int }

func (i *interp) run(doPanic bool) (exitCode int, err error) {
	defer func() {
		if doPanic {
			return // don't catch panics if we're panicking
		}
		if r := recover(); r != nil {
			if exit, ok := r.(gooseExit); ok {
				exitCode = int(exit.code)
			} else {
				exitCode = PANIC
				err = fmt.Errorf("%s", r)
			}
		}
	}()
	for k, v := range globals {
		i.global.Set(k, &Variable{
			Constant: false, // allow overwriting of builtins
			Value:    &Func{Executor: v},
		})
	}

	i.runModule(i.currentModule())
	return 0, nil
}

func (i *interp) runModule(module *Module) {
	i.executionStack = append(i.executionStack, module)

	if i.stdout == nil {
		i.stdout = os.Stdout
	}
	if i.stderr == nil {
		i.stderr = os.Stderr
	}

	for _, stmt := range module.Stmts {
		result := i.runStmt(module.Scope, stmt)

		switch result.(type) {
		case *Return:
			i.throw("cannot return from top-level")
		case *Break:
			i.throw("cannot break from top-level")
		case *Continue:
			i.throw("cannot continue from top-level")
		}
	}

	i.executionStack = i.executionStack[:len(i.executionStack)-1]
}
