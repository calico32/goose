package interpreter

import (
	"github.com/calico32/goose/ast"
	"github.com/calico32/goose/token"

	. "github.com/calico32/goose/interpreter/lib"
)

func (i *interp) evalMatchExpr(scope *Scope, expr *ast.MatchExpr) Value {
	defer un(trace(i, "match expr"))

	x := i.evalExpr(scope, expr.Expr)

	for _, clause := range expr.Clauses {
		switch clause := clause.(type) {
		case *ast.MatchPattern:
			matched, scope := i.matchPattern(scope, clause.Pattern, x)
			if matched {
				return i.evalExpr(scope, clause.Expr)
			}
		case *ast.MatchElse:
			return i.evalExpr(scope, clause.Expr)
		default:
			i.Throw("unknown match clause type: %T", clause)
		}
	}

	return NullValue
}

func (i *interp) matchPattern(scope *Scope, pattern ast.PatternExpr, x Value) (bool, *Scope) {
	switch pattern := pattern.(type) {
	case *ast.PatternBinding:
		// always matches, bind the value to the name
		scope = scope.Fork(ScopeOwnerMatch)
		scope.Set(pattern.Ident.Name, &Variable{Value: x})
		return true, scope
	case *ast.PatternNormal:
		// evaluate the pattern
		y := i.evalExpr(scope, pattern.X)
		switch y := y.(type) {
		case *IntRange:
			if x, ok := x.(*Integer); ok {
				return y.Contains(x), scope
			}
		case *FloatRange:
			if x, ok := x.(*Float); ok {
				return y.Contains(x), scope
			}
		default:
			op := GetOperator(x, token.Eq)
			if op == nil {
				i.Throw("operator %s not defined for type %s", token.Eq, x.Type())
			}
			ret := op.Executor(&FuncContext{
				Interp: i,
				Scope:  scope,
				This:   x,
				Args:   []Value{y},
			})
			return IsTruthy(ret.Value), scope
		}
	case *ast.PatternTuple, *ast.PatternComposite, *ast.PatternParen, *ast.PatternType:
		i.Throw("pattern type %T not implemented", pattern)
	default:
		i.Throw("unknown pattern type: %T", pattern)
	}

	return false, nil
}
