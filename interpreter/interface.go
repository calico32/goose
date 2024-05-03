package interpreter

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/calico32/goose/ast"
	. "github.com/calico32/goose/interpreter/lib"
	"github.com/calico32/goose/token"
)

func New(file *ast.Module, fset *token.FileSet, trace bool, stdin io.Reader, stdout io.Writer, stderr io.Writer) (i *interp, err error) {
	i = &interp{
		modules:        make(map[string]*Module),
		global:         NewGlobalScope(GlobalConstants),
		trace:          trace,
		fset:           fset,
		stdin:          stdin,
		stdout:         stdout,
		stderr:         stderr,
		gooseRoot:      os.Getenv("GOOSEROOT"),
		executionStack: make([]*Module, 0, 10),
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

	err = CreateGooseRoot(i.gooseRoot)
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
	i.executionStack = append(i.executionStack, module)

	return
}

func EvalExpr(expr ast.Expr, scope *Scope) (ret Value, err error) {
	i := &interp{
		fset:      token.NewFileSet(),
		modules:   make(map[string]*Module),
		global:    NewGlobalScope(GlobalConstants),
		gooseRoot: os.Getenv("GOOSEROOT"),
	}

	module := &Module{
		Module: &ast.Module{
			Stmts: []ast.Stmt{
				&ast.ExprStmt{
					X: expr,
				},
			},
		},
		Scope:   scope.Reparent(i.global),
		Exports: make(map[string]*Variable),
	}

	if module.Scope == nil {
		module.Scope = i.global.Fork(ScopeOwnerModule)
	}

	module.Scope.SetModule(module)

	i.modules["<eval>"] = module
	i.executionStack = append(i.executionStack, module)

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", r)
		}
	}()

	return i.evalExpr(scope, expr), nil
}
