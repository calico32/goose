package interpreter

import "github.com/calico32/goose/ast"

func (i *interp) evalSelectorExpr(scope *Scope, expr *ast.SelectorExpr) Value {
	defer un(trace(i, "selector expr"))

	x := i.evalExpr(scope, expr.X)

	if _, ok := x.(*Null); ok {
		i.throw("cannot access property %s of null", expr.Sel.Name)
	}

	return GetProperty(x, &String{expr.Sel.Name})
}

func (i *interp) evalBracketSelectorExpr(scope *Scope, expr *ast.BracketSelectorExpr) Value {
	defer un(trace(i, "bracket selector expr"))

	x := i.evalExpr(scope, expr.X)

	if _, ok := x.(*Null); ok {
		i.throw("cannot access property of null")
	}

	sel := i.evalExpr(scope, expr.Sel)

	if _, ok := sel.(PropertyKey); !ok {
		i.throw("cannot use %s as property key", sel.Type())
	}

	return GetProperty(x, sel.(PropertyKey))
}
