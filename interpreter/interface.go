package interpreter

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/calico32/goose/ast"
	"github.com/calico32/goose/token"
)

func RunFile(file *ast.File, fset *token.FileSet, trace bool, doPanic bool, stdout io.Writer, stderr io.Writer) (exitCode int, err error) {
	i := &interp{
		modules:   make(map[string]*Module),
		global:    NewGlobalScope(builtins),
		trace:     trace,
		fset:      fset,
		stdout:    stdout,
		stderr:    stderr,
		gooseRoot: os.Getenv("GOOSEROOT"),
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
		return 1, err
	}

	module := &Module{
		File:    file,
		Scope:   i.global.Fork(ScopeOwnerModule),
		Exports: make(map[string]*Variable),
	}

	module.Scope.module = module
	i.modules[file.Name] = module

	return i.run(doPanic)
}

func EvalExpr(expr ast.Expr, scope *Scope) (ret Value, err error) {
	i := &interp{
		fset:      token.NewFileSet(),
		modules:   make(map[string]*Module),
		global:    NewGlobalScope(builtins),
		gooseRoot: os.Getenv("GOOSEROOT"),
	}

	module := &Module{
		File: &ast.File{
			Stmts: []ast.Stmt{
				&ast.ExprStmt{
					X: expr,
				},
			},
		},
		Scope:   scope,
		Exports: make(map[string]*Variable),
	}

	if module.Scope == nil {
		module.Scope = i.global.Fork(ScopeOwnerModule)
	}

	module.Scope.module = module

	i.modules["<eval>"] = module
	i.executionStack = append(i.executionStack, module)

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", r)
		}
	}()

	return i.evalExpr(scope, expr), nil
}
