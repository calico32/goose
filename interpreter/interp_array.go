package interpreter

import "github.com/calico32/goose/ast"

func (i *interp) evalArrayLiteral(scope *Scope, expr *ast.ArrayLiteral) Value {
	defer un(trace(i, "array lit"))

	var values []Value
	for _, value := range expr.List {
		val := i.evalExpr(scope, value)

		values = append(values, val)
	}

	return wrap(values)
}

func (i *interp) evalArrayInitializer(scope *Scope, expr *ast.ArrayInitializer) Value {
	defer un(trace(i, "array initializer"))

	var values []Value

	count := i.evalExpr(scope, expr.Count)

	if _, ok := count.(*Integer); !ok {
		i.throw("array initializer count must be an integer")
	}

	countVal := count.(*Integer).Value

	if lit, ok := expr.Value.(*ast.Literal); ok {
		// don't bother evaluating the value if it's a literal
		value := i.evalLiteral(scope, lit)

		for i := int64(0); i < countVal; i++ {
			values = append(values, value)
		}

		return wrap(values)
	}

	for idx := int64(0); idx < countVal; idx++ {
		newScope := scope.Fork(ScopeOwnerArrayInit)
		newScope.Set("_", &Variable{
			Constant: true,
			Value:    wrap(idx),
		})

		val := i.evalExpr(newScope, expr.Value)

		values = append(values, val)
	}

	return wrap(values)
}

func (i *interp) evalSliceExpr(scope *Scope, expr *ast.SliceExpr) Value {
	defer un(trace(i, "slice expression"))

	left := i.evalExpr(scope, expr.X)

	if _, ok := left.(*Array); !ok {
		i.throw("cannot slice non-array type %s", left.Type())
	}

	values := left.(*Array).Elements
	start := int64(0)
	end := int64(len(values))

	if expr.Low != nil {
		begin := i.evalExpr(scope, expr.Low)

		if _, ok := begin.(Numeric); !ok {
			i.throw("expected numeric type for slice start index, got %s", begin.Type())
		}

		idx := begin.(Numeric).Int64()
		if idx < 0 {
			idx = int64(len(values)) + idx
			if idx < 0 {
				idx = 0
			}
			if idx >= int64(len(values)) {
				i.throw("index %d out of bounds", begin.(Numeric).Int64())
			}
		}

		start = idx
	}

	if expr.High != nil {
		endVal := i.evalExpr(scope, expr.High)

		if _, ok := endVal.(Numeric); !ok {
			i.throw("expected numeric type for slice end index, got %s", endVal.Type())
		}

		idx := endVal.(Numeric).Int64()
		if idx < 0 {
			idx = int64(len(values)) + idx
			if idx < 0 {
				idx = 0
			}
			if idx >= int64(len(values)) {
				i.throw("index %d out of bounds", endVal)
			}
		}
		end = idx
	}

	if end < start {
		i.throw("computed end index %d is less than computed start index %d", end, start)
	}

	if start < 0 {
		start = 0
	}

	if end > int64(len(values)) {
		end = int64(len(values))
	}

	return wrap(values[start:end])
}
