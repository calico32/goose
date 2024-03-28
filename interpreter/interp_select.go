package interpreter

import (
	"github.com/calico32/goose/ast"
	. "github.com/calico32/goose/interpreter/lib"
)

func (i *interp) evalSelectorExpr(scope *Scope, expr *ast.SelectorExpr) Value {
	defer un(trace(i, "selector expr"))

	x := i.evalExpr(scope, expr.X)

	if _, ok := x.(*Null); ok {
		i.Throw("cannot access property %s of null", expr.Sel.Name)
	}

	return GetProperty(x, NewString(expr.Sel.Name))
}

func (i *interp) evalBracketSelectorExpr(scope *Scope, expr *ast.BracketSelectorExpr) Value {
	defer un(trace(i, "bracket selector expr"))

	x := i.evalExpr(scope, expr.X)

	if _, ok := x.(*Null); ok {
		i.Throw("cannot access property of null")
	}

	if _, ok := x.(*String); ok {
		sel := i.evalExpr(scope, expr.Sel)
		if _, ok := sel.(*Integer); !ok {
			i.Throw("cannot use %s as index", sel.Type())
		}
		char := x.(*String).Value[sel.(Numeric).Int()]
		return Wrap(char)
	}

	sel := i.evalExpr(scope, expr.Sel)

	if _, ok := sel.(PropertyKey); !ok {
		i.Throw("cannot use %s as property key", sel.Type())
	}

	return GetProperty(x, sel.(PropertyKey))
}
