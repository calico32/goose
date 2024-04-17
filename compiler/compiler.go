package compiler

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/calico32/goose/ast"
	"github.com/calico32/goose/interpreter"
	"github.com/calico32/goose/token"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"

	. "github.com/calico32/goose/interpreter/lib"
)

type Compiler struct {
	fset        *token.FileSet
	moduleStack []*Module
	modules     map[string]*Module
	global      *Scope
	stdin       io.Reader
	stdout      io.Writer
	stderr      io.Writer
	gooseRoot   string

	// internal state
	trace    bool
	indent   int
	lastNode ast.Node
	stack    []ast.Node
}

func trace(i *Compiler, msg string) *Compiler {
	if i.trace {
		i.printTrace(msg, "(")
		i.indent++
	}
	return i
}

// Usage pattern: defer un(trace(p, "..."))
func un(i *Compiler) {
	if i.trace {
		i.indent--
		i.printTrace(")")
	}
}

func (c *Compiler) printTrace(a ...any) {
	const dots = ". . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . "
	const n = len(dots)
	i := 2 * c.indent
	for i > n {
		fmt.Fprint(c.stderr, dots)
		i -= n
	}
	// i <= n
	fmt.Fprint(c.stderr, dots[0:i])
	for _, arg := range a {
		fmt.Fprintf(c.stderr, "%v ", arg)
	}
}

func (c *Compiler) Fset() *token.FileSet        { return c.fset }
func (c *Compiler) ImportStack() []*Module      { return c.moduleStack }
func (c *Compiler) Modules() map[string]*Module { return c.modules }
func (c *Compiler) Global() *Scope              { return c.global }
func (c *Compiler) Stdin() io.Reader            { return c.stdin }
func (c *Compiler) Stdout() io.Writer           { return c.stdout }
func (c *Compiler) Stderr() io.Writer           { return c.stderr }
func (c *Compiler) GooseRoot() string           { return c.gooseRoot }

func (c *Compiler) CurrentModule() *Module {
	if len(c.moduleStack) == 0 {
		if len(c.modules) == 1 {
			for _, m := range c.modules {
				c.moduleStack = append(c.moduleStack, m)
				return m
			}
		}
		c.Throw("no current module")
	}

	return c.moduleStack[len(c.moduleStack)-1]
}

func (c *Compiler) Throw(msg string, parts ...any) {
	panic(fmt.Errorf("%s: Validation error: %s", c.fset.Position(c.currentNode().Pos()), fmt.Sprintf(msg, parts...)))
}

func (c *Compiler) currentNode() ast.Node {
	if len(c.stack) == 0 {
		return c.lastNode
	}
	return c.stack[len(c.stack)-1]
}

func push(c *Compiler, n ast.Node) *Compiler {
	c.lastNode = n
	c.stack = append(c.stack, n)
	return c
}

func pop(c *Compiler) {
	c.stack = c.stack[:len(c.stack)-1]
}

func New(file *ast.Module, fset *token.FileSet, trace bool, stdin io.Reader, stdout io.Writer, stderr io.Writer) (i *Compiler, err error) {
	i = &Compiler{
		modules:     make(map[string]*Module),
		global:      NewGlobalScope(interpreter.GlobalConstants),
		trace:       trace,
		fset:        fset,
		stdin:       stdin,
		stdout:      stdout,
		stderr:      stderr,
		gooseRoot:   os.Getenv("GOOSEROOT"),
		moduleStack: make([]*Module, 0, 10),
	}

	if i.gooseRoot == "" {
		if xdgDataHome := os.Getenv("XDG_DATA_HOME"); xdgDataHome != "" {

			i.gooseRoot = filepath.Join(xdgDataHome, "goose")
		} else {
			home := os.Getenv("HOME")
			if home == "" {
				home = os.Getenv("USERPROFILE")
			}

			i.gooseRoot = filepath.Join(home, ".goose")
		}
	}

	err = interpreter.CreateGooseRoot(i.gooseRoot)
	if err != nil {
		return
	}

	module := &Module{
		Module:  file,
		Scope:   i.global.Fork(ScopeOwnerModule),
		Exports: make(map[string]*Variable),
	}

	module.Scope.SetModule(module)
	i.modules[file.Specifier] = module
	i.moduleStack = append(i.moduleStack, module)

	return
}

func (c *Compiler) Compile() string {
	m := ir.NewModule()
	m.SourceFilename = c.CurrentModule().Module.Specifier

	printf := m.NewFunc("printf", types.I32, ir.NewParam("format", types.I8Ptr), ir.NewParam("...", types.I8Ptr))

	main := m.NewFunc("main", types.I32)
	mainBlock := main.NewBlock("")
	mainBlock.NewCall(printf, constant.NewCharArrayFromString("%d\n"), constant.NewInt(types.I32, 42))
	mainBlock.NewRet(constant.NewInt(types.I32, 0))

	return m.String()
}
