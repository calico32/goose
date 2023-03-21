package interpreter

import (
	"strings"

	"github.com/calico32/goose/ast"
	"github.com/calico32/goose/token"
)

func toDebugString(x ast.Expr) string {
	switch expr := x.(type) {
	case *ast.Ident:
		return expr.Name
	case *ast.SelectorExpr:
		return toDebugString(expr.X) + "." + expr.Sel.Name
	case *ast.BracketSelectorExpr:
		return toDebugString(expr.X) + "[" + toDebugString(expr.Sel) + "]"
	case *ast.SliceExpr:
		return toDebugString(expr.X) + "[" + toDebugString(expr.Low) + ":" + toDebugString(expr.High) + "]"
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
				output.WriteString(toDebugString(part.Expr))
				output.WriteString("}")
			case *ast.StringLiteralInterpIdent:
				output.WriteString("$")
				output.WriteString(part.Name)
			case *ast.StringLiteralMiddle:
				output.WriteString(part.Content)
			}
		}
		output.WriteString(expr.StringEnd.Content)
		output.WriteString("\"")
		return output.String()
	case *ast.ParenExpr:
		return "(" + toDebugString(expr.X) + ")"
	case *ast.UnaryExpr:
		if token.IsPostfixOperator(expr.Op) {
			return toDebugString(expr.X) + expr.Op.String()
		} else {
			return expr.Op.String() + toDebugString(expr.X)
		}
	case *ast.BinaryExpr:
		return toDebugString(expr.X) + " " + expr.Op.String() + " " + toDebugString(expr.Y)
	case *ast.CallExpr:
		var output strings.Builder
		output.WriteString(toDebugString(expr.Func))
		output.WriteString("(")
		for i, arg := range expr.Args {
			if i > 0 {
				output.WriteString(", ")
			}
			output.WriteString(toDebugString(arg))
		}
		output.WriteString(")")
		return output.String()
	case *ast.ArrayInitializer:
		var output strings.Builder
		output.WriteString("[")
		output.WriteString(toDebugString(expr.Value))
		output.WriteString("; ")
		output.WriteString(toDebugString(expr.Count))
		output.WriteString("]")
		return output.String()
	}

	return "<unknown>"
}

func toString(i *interp, scope *Scope, v Value) string {
	prop := GetProperty(v, &String{"toString"})
	if prop == nil {
		return "<unknown>"
	}

	if _, ok := prop.(*Func); !ok {
		return "<unknown>"
	}

	ret := prop.(*Func).Executor(&FuncContext{
		interp: i,
		Scope:  scope,
		This:   v,
	})

	return ret.Value.(*String).Value
}
