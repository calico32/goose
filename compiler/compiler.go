package compiler

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/calico32/goose/ast"
	"github.com/calico32/goose/compiler/cstring"
	"github.com/calico32/goose/interpreter"
	"github.com/calico32/goose/token"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	. "github.com/calico32/goose/interpreter/lib"
)

type Compiler struct {
	fset        *token.FileSet
	moduleStack []*CModule
	modules     map[string]*CModule
	global      *CScope
	stdin       io.Reader
	stdout      io.Writer
	stderr      io.Writer
	gooseRoot   string
	main        *ir.Func
	m           *ir.Module
	ext         *ExternalFunctionRegistry
	block       *ir.Block

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

func (c *Compiler) Fset() *token.FileSet         { return c.fset }
func (c *Compiler) ImportStack() []*CModule      { return c.moduleStack }
func (c *Compiler) Modules() map[string]*CModule { return c.modules }
func (c *Compiler) Global() *CScope              { return c.global }
func (c *Compiler) Stdin() io.Reader             { return c.stdin }
func (c *Compiler) Stdout() io.Writer            { return c.stdout }
func (c *Compiler) Stderr() io.Writer            { return c.stderr }
func (c *Compiler) GooseRoot() string            { return c.gooseRoot }

func (c *Compiler) CurrentModule() *CModule {
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

func New(file *ast.Module, fset *token.FileSet, trace bool, stdin io.Reader, stdout io.Writer, stderr io.Writer) (c *Compiler, err error) {
	c = &Compiler{
		modules:     make(map[string]*CModule),
		global:      NewGlobalCScope(CGlobalConstants),
		trace:       trace,
		fset:        fset,
		stdin:       stdin,
		stdout:      stdout,
		stderr:      stderr,
		gooseRoot:   os.Getenv("GOOSEROOT"),
		moduleStack: make([]*CModule, 0, 10),
	}

	if c.gooseRoot == "" {
		if xdgDataHome := os.Getenv("XDG_DATA_HOME"); xdgDataHome != "" {

			c.gooseRoot = filepath.Join(xdgDataHome, "goose")
		} else {
			home := os.Getenv("HOME")
			if home == "" {
				home = os.Getenv("USERPROFILE")
			}

			c.gooseRoot = filepath.Join(home, ".goose")
		}
	}

	err = interpreter.CreateGooseRoot(c.gooseRoot)
	if err != nil {
		return
	}

	c.ext = NewExternalFunctionRegistry(c)

	module := &CModule{
		Module:  file,
		Exports: make(map[string]*CVariable),
		Scope:   c.global.Fork(ScopeOwnerModule),
	}

	module.Scope.SetModule(module)
	c.modules[file.Specifier] = module
	c.moduleStack = append(c.moduleStack, module)

	return
}

func (c *Compiler) Compile() string {
	m := ir.NewModule()
	c.m = m
	m.SourceFilename = c.CurrentModule().Module.Specifier
	c.CurrentModule().Scope.fn = c.main
	c.main = m.NewFunc("main", types.I32)

	// c.gooseValue := m.NewTypeDef("goose_value", types.NewStruct(
	// 	types.I8,  // type
	// 	types.I64, // data/length
	// 	types.NewPointer(types.I8),
	// ))

	for k, v := range CGlobals {
		c.global.Set(k, &CVariable{
			Constant: false,
			Value:    v(c),
		})
	}
	str := m.NewGlobalDef(cstring.NextStringName(), cstring.Constant("Hello, world!\n"))
	str.Immutable = true

	bb := c.main.NewBlock("entry")
	c.block = bb

	module := c.CurrentModule()
	for _, stmt := range module.Module.Stmts {
		c.CompileStmt(stmt, module.Scope)
	}

	return m.String()
}

func (c *Compiler) CompileStmt(stmt ast.Stmt, scope *CScope) {
	defer un(trace(c, "Stmt"))
	switch stmt := stmt.(type) {
	case *ast.LetStmt:
		c.CompileLetStmt(stmt, scope)
	case *ast.AssignStmt:
		c.CompileAssignStmt(stmt, scope)
	case *ast.ExprStmt:
		c.CompileExpr(stmt.X, scope)
	default:
		fmt.Printf("didn't handle statement type: %T\n", stmt)
	}
}

func (c *Compiler) CompileLetStmt(stmt *ast.LetStmt, scope *CScope) {
	defer un(trace(c, "LetStmt"))
	defer pop(push(c, stmt))

	if stmt.Ident.Name == "_" {
		c.Throw("cannot declare _")
	}

	if scope.IsDefinedInCurrentScope(stmt.Ident.Name) {
		c.Throw("cannot redefine variable %s", stmt.Ident.Name)
	}

	value := c.CompileExpr(stmt.Value, scope)

	scope.Set(stmt.Ident.Name, &CVariable{
		Constant: false,
		Value:    value,
		Source:   VariableSourceDecl,
	})
}

func (c *Compiler) CompileAssignStmt(stmt *ast.AssignStmt, scope *CScope) {
	defer un(trace(c, "AssignStmt"))
	defer pop(push(c, stmt))

}

func (c *Compiler) CompileExpr(x ast.Expr, scope *CScope) value.Value {
	defer un(trace(c, "Expr"))
	switch x := x.(type) {
	case *ast.BinaryExpr:
		return c.CompileBinaryExpr(x, scope)
	case *ast.CallExpr:
		return c.CompileCallExpr(x, scope)
	case *ast.Ident:
		return c.CompileIndent(x, scope)
	default:
		fmt.Printf("didn't handle expression type: %T\n", x)
	}

	return nil
}

func (c *Compiler) CompileIndent(x *ast.Ident, scope *CScope) value.Value {
	defer un(trace(c, "Ident"))
	defer pop(push(c, x))

	if v := scope.Get(x.Name); v != nil {
		return v.Value
	}

	if v := c.global.Get(x.Name); v != nil {
		return v.Value
	}

	c.Throw("undeclared identifier %s", x.Name)
	return nil
}

func (c *Compiler) CompileBinaryExpr(x *ast.BinaryExpr, scope *CScope) value.Value {
	defer un(trace(c, "BinaryExpr"))
	defer pop(push(c, x))

	lhs := c.CompileExpr(x.X, scope)
	rhs := c.CompileExpr(x.Y, scope)

	switch x.Op {
	case token.Add:
		return c.block.NewAdd(lhs, rhs)
	case token.Sub:
		return c.block.NewSub(lhs, rhs)
	case token.Mul:
		return c.block.NewMul(lhs, rhs)
	case token.Quo:
		return c.block.NewSDiv(lhs, rhs)
	case token.Rem:
		return c.block.NewSRem(lhs, rhs)
	default:
		fmt.Println("didn't handle binary operator:", x.Op)
		return constant.NewInt(types.I32, 0)
	}
}

func (c *Compiler) CompileCallExpr(x *ast.CallExpr, scope *CScope) value.Value {
	defer un(trace(c, "CallExpr"))
	defer pop(push(c, x))

	fn := c.CompileExpr(x.Func, scope)
	args := make([]value.Value, len(x.Args))
	for i, arg := range x.Args {
		args[i] = c.CompileExpr(arg, scope)
	}

	return c.block.NewCall(fn, args...)
}
