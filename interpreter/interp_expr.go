package interpreter

import (
	"fmt"

	"github.com/calico32/goose/ast"
	"github.com/calico32/goose/token"
)

func (i *interp) evalExpr(scope *Scope, expr ast.Expr) Value {
	defer un(trace(i, "expr"))
	defer pop(push(i, expr))

	switch expr := expr.(type) {
	case *ast.BinaryExpr:
		return i.evalBinaryExpr(scope, expr)
	case *ast.UnaryExpr:
		return i.evalUnaryExpr(scope, expr)
	case *ast.ParenExpr:
		return i.evalExpr(scope, expr.X)
	case *ast.CallExpr:
		return i.evalCallExpr(scope, expr)
	case *ast.Ident:
		return i.evalIdent(scope, expr)
	case *ast.StringLiteral:
		return i.evalString(scope, expr)
	case *ast.CompositeLiteral:
		return i.evalCompositeLiteral(scope, expr)
	case *ast.ArrayLiteral:
		return i.evalArrayLiteral(scope, expr)
	case *ast.ArrayInitializer:
		return i.evalArrayInitializer(scope, expr)
	case *ast.SelectorExpr:
		return i.evalSelectorExpr(scope, expr)
	case *ast.BindExpr:
		return i.evalBindExpr(scope, expr)
	case *ast.BracketSelectorExpr:
		return i.evalBracketSelectorExpr(scope, expr)
	case *ast.SliceExpr:
		return i.evalSliceExpr(scope, expr)
	case *ast.Literal:
		return i.evalLiteral(scope, expr)
	case *ast.FuncExpr:
		return i.evalFuncExpr(scope, expr)
	case *ast.IfExpr:
		return i.evalIfExpr(scope, expr)
	case *ast.DoExpr:
		return i.evalDoExpr(scope, expr)
	case *ast.GeneratorExpr:
		return i.evalGeneratorExpr(scope, expr)
	default:
		if badExpr, ok := expr.(*ast.BadExpr); ok {
			i.throw("unexpected bad expression %#v", badExpr)
		}

		i.throw("unexpected expression type %T", expr)
	}

	return nil
}

func (i *interp) evalBinaryExpr(scope *Scope, expr *ast.BinaryExpr) Value {
	defer un(trace(i, "binary expr"))

	left := i.evalExpr(scope, expr.X)

	if expr.Op == token.Arrow {
		// pass the left side as _ to the right side
		newScope := scope.Fork(ScopeOwnerPipeline)
		newScope.Set("_", &Variable{
			Value: left,
		})
		return i.evalExpr(newScope, expr.Y)
	}

	right := i.evalExpr(scope, expr.Y)

	op := GetOperator(left, expr.Op)
	if op == nil {
		i.throw("operator %s not defined for type %s", expr.Op, left.Type())
	}

	ret := op.Executor(&FuncContext{
		interp: i,
		Scope:  scope,
		This:   left,
		Args:   []Value{right},
	})
	return ret.Value
}

func (i *interp) evalUnaryExpr(scope *Scope, expr *ast.UnaryExpr) Value {
	defer un(trace(i, "unary expr"))

	value := i.evalExpr(scope, expr.X)

	switch expr.Op {
	case token.LogNot:
		return wrap(!isTruthy(value))
	case token.Question:
		if paren, ok := expr.X.(*ast.ParenExpr); ok {
			fmt.Println(toDebugString(paren.X) + " = " + toString(i, scope, value))
			return value
		} else {
			fmt.Println(toDebugString(expr.X) + " = " + toString(i, scope, value))
			return value
		}
	case token.LogNull, token.Add, token.Sub, token.BitNot:
		op := GetOperator(value, expr.Op)
		if op == nil {
			i.throw("operator %s not defined for type %s", expr.Op, value.Type())
		}

		ret := op.Executor(&FuncContext{
			interp: i,
			Scope:  scope,
			This:   value,
		})

		return ret.Value
	default:
		i.throw("unexpected unary operator %s", expr.Op)
		return nil
	}
}
