package interpreter

import (
	"math/big"
	"strconv"
	"strings"

	"github.com/calico32/goose/ast"
	. "github.com/calico32/goose/interpreter/lib"
	"github.com/calico32/goose/token"
)

func (i *interp) evalLiteral(_ *Scope, expr *ast.Literal) Value {
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
		val := new(big.Int)
		val, ok := val.SetString(strVal, base)
		if !ok {
			i.Throw("failed to parse integer")
		}

		return Wrap(val)

	case token.Float:
		val, err := strconv.ParseFloat(expr.Value, 64)
		if err != nil {
			i.Throw(err.Error())
		}

		return Wrap(val)

	case token.Null:
		return NullValue

	case token.True:
		return TrueValue

	case token.False:
		return FalseValue

	default:
		i.Throw("unexpected literal kind %s", expr.Kind)
		return nil
	}
}

func (i *interp) evalIdent(scope *Scope, expr *ast.Ident) Value {
	defer un(trace(i, "ident"))

	if expr.Name[0] == '#' {
		this := scope.Get("this")
		if this == nil {
			i.Throw("invalid property access: 'this' is not defined")
		}
		// property
		prop := GetProperty(this.Value, NewString(expr.Name[1:]))
		if fn, ok := prop.(*Func); ok {
			fn.This = this.Value
		}

		return prop
	}

	val := scope.Get(expr.Name)
	if val == nil {
		i.Throw("%s is not defined", expr.Name)
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
			val := i.evalIdent(scope, expr.Ident)
			value += ToString(i, scope, val)
		case *ast.StringLiteralInterpExpr:
			val := i.evalExpr(scope, expr.Expr)
			value += ToString(i, scope, val)
		default:
			i.Throw("unexpected string literal part %T", expr)
		}
	}

	value += expr.StringEnd.Content

	return Wrap(value)
}
