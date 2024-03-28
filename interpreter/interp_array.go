package interpreter

import (
	"math/big"
	"strings"

	"github.com/calico32/goose/ast"
	"github.com/calico32/goose/token"

	. "github.com/calico32/goose/interpreter/lib"
)

func (i *interp) evalArrayLiteral(scope *Scope, expr *ast.ArrayLiteral) Value {
	defer un(trace(i, "array lit"))

	var values []Value
	for _, ex := range expr.List {
		if ex, ok := ex.(*ast.UnaryExpr); ok {
			if ex.Op == token.Ellipsis {
				spread := i.evalExpr(scope, ex.X)
				switch v := spread.(type) {
				case *Array:
					values = append(values, v.Elements...)
				default:
					i.Throw("cannot spread non-array type %s", v.Type())
				}
				continue
			}
		}
		val := i.evalExpr(scope, ex)

		values = append(values, val)
	}

	return Wrap(values)
}

func (i *interp) evalArrayInitializer(scope *Scope, expr *ast.ArrayInitializer) Value {
	defer un(trace(i, "array initializer"))

	var values []Value

	count := i.evalExpr(scope, expr.Count)

	if _, ok := count.(*Integer); !ok {
		i.Throw("array initializer count must be an integer")
	}

	countVal := count.(*Integer).Value

	if lit, ok := expr.Value.(*ast.Literal); ok {
		// don't bother evaluating the expressiong every time, just use the literal value
		value := i.evalLiteral(scope, lit)

		for i := big.NewInt(0); i.Cmp(countVal) == -1; i.Add(i, big.NewInt(1)) {
			values = append(values, value)
		}

		return Wrap(values)
	}

	for idx := big.NewInt(0); idx.Cmp(countVal) == -1; idx.Add(idx, big.NewInt(1)) {
		newScope := scope.Fork(ScopeOwnerArrayInit)
		newScope.Set("_", &Variable{
			Constant: true,
			Value:    Wrap(idx),
		})

		val := i.evalExpr(newScope, expr.Value)

		values = append(values, val)
	}

	return Wrap(values)
}

// TODO: optimize for strings
func (i *interp) evalSliceExpr(scope *Scope, expr *ast.SliceExpr) Value {
	defer un(trace(i, "slice expression"))

	x := i.evalExpr(scope, expr.X)

	var values []Value
	if a, ok := x.(*Array); ok {
		values = a.Elements
	} else if a, ok := x.(*String); ok {
		values = []Value{}
		for _, r := range a.Value {
			values = append(values, Wrap(string(r)))
		}
	} else {
		i.Throw("cannot slice non-array and non-string type %s", x.Type())
	}

	start := int64(0)
	end := int64(len(values))

	if expr.Low != nil {
		begin := i.evalExpr(scope, expr.Low)

		if _, ok := begin.(Numeric); !ok {
			i.Throw("expected numeric type for slice start index, got %s", begin.Type())
		}

		idx := begin.(Numeric).Int64()
		if idx < 0 {
			idx = int64(len(values)) + idx
			if idx < 0 {
				idx = 0
			}
			if idx >= int64(len(values)) {
				i.Throw("index %d out of bounds", begin.(Numeric).Int64())
			}
		}

		start = idx
	}

	if expr.High != nil {
		endVal := i.evalExpr(scope, expr.High)

		if _, ok := endVal.(Numeric); !ok {
			i.Throw("expected numeric type for slice end index, got %s", endVal.Type())
		}

		idx := endVal.(Numeric).Int64()
		if idx < 0 {
			idx = int64(len(values)) + idx
			if idx < 0 {
				idx = 0
			}
			if idx >= int64(len(values)) {
				i.Throw("index %d out of bounds", endVal)
			}
		}
		end = idx
	}

	if end < start {
		i.Throw("computed end index %d is less than computed start index %d", end, start)
	}

	if start < 0 {
		start = 0
	}

	if end > int64(len(values)) {
		end = int64(len(values))
	}

	result := values[start:end]

	if _, ok := x.(*Array); ok {
		return Wrap(result)
	} else {
		var out strings.Builder
		for _, v := range result {
			out.WriteString(v.(*String).Value)
		}
		return Wrap(out.String())
	}
}
