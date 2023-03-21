package interpreter

import "github.com/calico32/goose/ast"

type Module struct {
	*ast.File
	Scope   *Scope
	Exports map[string]*Variable
}
