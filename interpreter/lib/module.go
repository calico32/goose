package lib

import "github.com/calico32/goose/ast"

type Module struct {
	*ast.Module
	Scope   *Scope
	Exports map[string]*Variable
}
