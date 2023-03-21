package interpreter

import (
	"strconv"
	"strings"

	"github.com/calico32/goose/ast"
	"github.com/calico32/goose/token"
)

func (i *interp) evalLiteral(scope *Scope, expr *ast.Literal) Value {
	switch expr.Kind {
	case token.Int:
		strVal := strings.ReplaceAll(expr.Value, "_", "")
		base := 10
		switch {
		case strings.HasPrefix(expr.Value, "0x"):
			strVal = strVal[2:]
			base = 16
		case strings.HasPrefix(expr.Value, "0o"):
			strVal = strVal[2:]
			base = 8
		case strings.HasPrefix(expr.Value, "0b"):
			strVal = strVal[2:]
			base = 2
		}
		val, err := strconv.ParseInt(strVal, base, 64)
		if err != nil {
			i.throw(err.Error())
		}

		return wrap(val)

	case token.Float:
		val, err := strconv.ParseFloat(expr.Value, 64)
		if err != nil {
			i.throw(err.Error())
		}

		return wrap(val)

	case token.Null:
		return NullValue

	case token.True:
		return TrueValue

	case token.False:
		return FalseValue

	default:
		i.throw("unexpected literal kind %s", expr.Kind)
		return nil
	}
}

func (i *interp) evalIdent(scope *Scope, expr *ast.Ident) Value {
	defer un(trace(i, "ident"))

	if expr.Name[0] == '#' {
		x := scope.Get("this")
		if x == nil {
			i.throw("invalid property access: 'this' is not defined")
		}
		// property
		prop := GetProperty(x.Value, &String{expr.Name[1:]})
		if fn, ok := prop.(*Func); ok {
			fn.This = x.Value
		}

		return prop
	}

	val := scope.Get(expr.Name)
	if val == nil {
		i.throw("%s is not defined", expr.Name)
	}

	return val.Value
}

// alex was here
func (i *interp) evalString(scope *Scope, expr *ast.StringLiteral) Value {
	defer un(trace(i, "string"))

	value := expr.StringStart.Content

	for _, part := range expr.Parts {
		switch expr := part.(type) {
		case *ast.StringLiteralMiddle:
			value += expr.Content
		case *ast.StringLiteralInterpIdent:
			val := i.evalIdent(scope, &ast.Ident{Name: expr.Name})
			value += toString(i, scope, val)
		case *ast.StringLiteralInterpExpr:
			val := i.evalExpr(scope, expr.Expr)
			value += toString(i, scope, val)
		default:
			i.throw("unexpected string literal part %T", expr)
		}
	}

	value += expr.StringEnd.Content

	return wrap(value)
}
