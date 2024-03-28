package interpreter

import (
	"io"
	"os"

	"github.com/calico32/goose/ast"
	. "github.com/calico32/goose/interpreter/lib"
	"github.com/calico32/goose/token"
)

type interp struct {
	fset           *token.FileSet
	executionStack []*Module
	callStack      []*CallFrame
	modules        map[string]*Module
	global         *Scope
	stdin          io.Reader
	stdout         io.Writer
	stderr         io.Writer
	gooseRoot      string

	// internal state
	trace    bool
	indent   int
	lastPos  token.Pos
	posStack []token.Pos
}

type CallFrame struct {
	Module *Module
	Node   ast.Node
}

type Exception struct {
	Message string
	Cause   *Exception
	Stack   []*CallFrame
}

func (i *interp) Fset() *token.FileSet        { return i.fset }
func (i *interp) ExecutionStack() []*Module   { return i.executionStack }
func (i *interp) CallStack() []*CallFrame     { return i.callStack }
func (i *interp) Modules() map[string]*Module { return i.modules }
func (i *interp) Global() *Scope              { return i.global }
func (i *interp) Stdin() io.Reader            { return i.stdin }
func (i *interp) Stdout() io.Writer           { return i.stdout }
func (i *interp) Stderr() io.Writer           { return i.stderr }
func (i *interp) GooseRoot() string           { return i.gooseRoot }

func (i *interp) CurrentModule() *Module {
	if len(i.executionStack) == 0 {
		if len(i.modules) == 1 {
			for _, m := range i.modules {
				i.executionStack = append(i.executionStack, m)
				return m
			}
		}
		i.Throw("no current module")
	}

	return i.executionStack[len(i.executionStack)-1]
}

const PANIC int = 128

type gooseExit struct{ code int }

func (i *interp) Run() (exitCode int, err error) {
	defer func() {
		if r := recover(); r != nil {
			if exit, ok := r.(gooseExit); ok {
				exitCode = int(exit.code)
			} else {
				panic(r)
			}
		}
	}()
	for k, v := range Globals {
		i.global.Set(k, &Variable{
			Constant: true,
			Value:    &Func{Executor: v},
		})
	}

	i.runBuiltins()
	i.runModule(i.CurrentModule())
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
			i.Throw("cannot return from top-level")
		case *Break:
			i.Throw("cannot break from top-level")
		case *Continue:
			i.Throw("cannot continue from top-level")
		}
	}

	i.executionStack = i.executionStack[:len(i.executionStack)-1]
}

func (i *interp) runBuiltins() {
	specs := []string{
		"std:language/builtin.goose",
	}

	for _, spec := range specs {
		mod := i.loadStdModule(spec)
		i.copyModuleExportsToGlobal(mod)
	}
}
