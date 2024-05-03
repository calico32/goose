package lib

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/calico32/goose/ast"
	"github.com/calico32/goose/token"
)

func PrintExpr(x ast.Expr) string {
	switch expr := x.(type) {
	case *ast.Ident:
		return expr.Name
	case *ast.SelectorExpr:
		return PrintExpr(expr.X) + "." + expr.Sel.Name
	case *ast.BracketSelectorExpr:
		return PrintExpr(expr.X) + "[" + PrintExpr(expr.Sel) + "]"
	case *ast.SliceExpr:
		return PrintExpr(expr.X) + "[" + PrintExpr(expr.Low) + ":" + PrintExpr(expr.High) + "]"
	case *ast.Literal:
		return expr.Value
	case *ast.StringLiteral:
		var output strings.Builder
		output.WriteString("\"")
		output.WriteString(expr.StringStart.Content)
		for _, part := range expr.Parts {
			switch part := part.(type) {
			case *ast.StringLiteralInterpExpr:
				output.WriteString("${")
				output.WriteString(PrintExpr(part.Expr))
				output.WriteString("}")
			case *ast.StringLiteralInterpIdent:
				output.WriteString("$")
				output.WriteString(part.Ident.Name)
			case *ast.StringLiteralMiddle:
				output.WriteString(part.Content)
			}
		}
		output.WriteString(expr.StringEnd.Content)
		output.WriteString("\"")
		return output.String()
	case *ast.ParenExpr:
		return "(" + PrintExpr(expr.X) + ")"
	case *ast.UnaryExpr:
		if token.IsPostfixOperator(expr.Op) {
			return PrintExpr(expr.X) + expr.Op.String()
		} else {
			return expr.Op.String() + PrintExpr(expr.X)
		}
	case *ast.BinaryExpr:
		return PrintExpr(expr.X) + " " + expr.Op.String() + " " + PrintExpr(expr.Y)
	case *ast.CallExpr:
		var output strings.Builder
		output.WriteString(PrintExpr(expr.Func))
		output.WriteString("(")
		for i, arg := range expr.Args {
			if i > 0 {
				output.WriteString(", ")
			}
			output.WriteString(PrintExpr(arg))
		}
		output.WriteString(")")
		return output.String()
	case *ast.ArrayInitializer:
		var output strings.Builder
		output.WriteString("[")
		output.WriteString(PrintExpr(expr.Value))
		output.WriteString("; ")
		output.WriteString(PrintExpr(expr.Count))
		output.WriteString("]")
		return output.String()
	}

	return fmt.Sprintf("<%T>", x)
}

func ToString(i Interpreter, scope *Scope, v Value) string {
	prop := GetProperty(v, NewString("toString"))
	if prop == nil {
		return "<unknown>"
	}

	if _, ok := prop.(*Func); !ok {
		return "<unknown>"
	}

	ret := prop.(*Func).Executor(&FuncContext{
		Interp: i,
		Scope:  scope,
		This:   v,
	})

	if _, ok := ret.Value.(*String); !ok {
		// in the case that a toString doesn't return a string, keep going
		// until we get a string
		return ToString(i, scope, ret.Value)
	}

	return ret.Value.(*String).Value
}

func ToDebugString(i Interpreter, scope *Scope, v Value, depth int) string {
	prop := GetProperty(v, NewString("toDebugString"))
	if prop == nil {
		return ToString(i, scope, v)
	}

	if _, ok := prop.(*Func); !ok {
		return "<unknown>"
	}

	ret := prop.(*Func).Executor(&FuncContext{
		Interp: i,
		Scope:  scope,
		This:   v,
		Args:   []Value{NewInteger(big.NewInt(int64(depth)))},
	})

	return ret.Value.(*String).Value
}
