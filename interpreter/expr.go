package interpreter

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/wiisportsresort/goose/ast"
	"github.com/wiisportsresort/goose/token"
)

func (i *interpreter) evalExpr(scope *GooseScope, expr ast.Expr) (result *GooseValue, err error) {
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
	case *ast.BracketSelectorExpr:
		return i.evalBracketSelectorExpr(scope, expr)
	case *ast.SliceExpr:
		return i.evalSliceExpr(scope, expr)
	case *ast.Literal:
		return i.evalLiteral(scope, expr)
	default:
		if badExpr, ok := expr.(*ast.BadExpr); ok {
			return nil, fmt.Errorf("unexpected bad expression %#v", badExpr)
		}
		return nil, fmt.Errorf("unexpected expression type %T", expr)
	}
}

func (i *interpreter) evalLiteral(scope *GooseScope, expr *ast.Literal) (*GooseValue, error) {
	switch expr.Kind {
	case token.Int:
		strVal := strings.Replace(expr.Value, "_", "", -1)
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
			return nil, err
		}

		return wrap(val), nil

	case token.Float:
		val, err := strconv.ParseFloat(expr.Value, 64)
		if err != nil {
			return nil, err
		}

		return wrap(val), nil

	case token.Null:
		return null, nil

	default:
		return nil, fmt.Errorf("unexpected literal kind %s", expr.Kind)
	}
}

func (i *interpreter) evalBinaryExpr(scope *GooseScope, expr *ast.BinaryExpr) (*GooseValue, error) {
	defer un(trace(i, "binary expr"))

	left, err := i.evalExpr(scope, expr.X)
	if err != nil {
		return nil, err
	}

	right, err := i.evalExpr(scope, expr.Y)
	if err != nil {
		return nil, err
	}

	return i.evalBinaryValues(left, expr.Op, right)
}

func (i *interpreter) evalBinaryValues(left *GooseValue, op token.Token, right *GooseValue) (*GooseValue, error) {
	defer un(trace(i, "binary values"))

	var boolResult bool
	switch op {
	case token.Assign:
		return right, nil
	case token.Add, token.AddAssign:
		if left.Type == GooseTypeArray {
			values := left.Value.([]*GooseValue)
			values = append(values, right)
			return wrap(values), nil
		} else if left.Type == GooseTypeString || right.Type == GooseTypeString {
			return wrap(fmt.Sprintf("%v", left.Value) + fmt.Sprintf("%v", right.Value)), nil
		}
		fallthrough
	case token.Lt, token.Lte, token.Gt, token.Gte,
		token.Sub, token.Mul, token.Quo, token.Rem, token.Pow,
		token.SubAssign, token.MulAssign, token.QuoAssign, token.RemAssign, token.PowAssign,
		token.Inc, token.Dec:
		return i.numericOperation(left, op, right)

	case token.Eq:
		boolResult = left.Type == right.Type && left.Value == right.Value
	case token.Neq:
		boolResult = left.Type != right.Type || left.Value != right.Value
	case token.LogAnd:
		boolResult = isTruthy(left.Value) && isTruthy(right.Value)
	case token.LogOr:
		boolResult = isTruthy(left.Value) || isTruthy(right.Value)
	default:
		return nil, fmt.Errorf("unexpected binary operator %s", op)
	}

	return wrap(boolResult), nil
}

func (i *interpreter) evalUnaryExpr(scope *GooseScope, expr *ast.UnaryExpr) (*GooseValue, error) {
	defer un(trace(i, "unary expr"))

	value, err := i.evalExpr(scope, expr.X)
	if err != nil {
		return nil, err
	}

	switch expr.Op {
	case token.LogNot:
		return wrap(!isTruthy(value.Value)), nil
	case token.Add:
		err = expectType(value, GooseTypeNumeric)
		if err != nil {
			return nil, err
		}
		return &GooseValue{
			Type:  value.Type,
			Value: value.Value,
		}, nil
	case token.Sub:
		err = expectType(value, GooseTypeNumeric)
		if err != nil {
			return nil, err
		}

		if value.Type == GooseTypeInt {
			return wrap(-value.Value.(int64)), nil
		}

		return wrap(-value.Value.(float64)), nil

	default:
		return nil, fmt.Errorf("unexpected unary operator %s", expr.Op)
	}
}

func (i *interpreter) evalCallExpr(scope *GooseScope, expr *ast.CallExpr) (*GooseValue, error) {
	defer un(trace(i, "call expr"))

	fn, err := i.evalExpr(scope, expr.Fun)
	if err != nil {
		return nil, err
	}

	if fn.Type != GooseTypeFunc {
		return nil, fmt.Errorf("expression of type %s is not callable", fn.Type)
	}

	args := make([]*GooseValue, len(expr.Args))
	for idx, arg := range expr.Args {
		val, err := i.evalExpr(scope, arg)
		if err != nil {
			return nil, err
		}

		args[idx] = val
	}

	result, err := fn.Value.(GooseFunc)(scope, args)
	if err != nil {
		return nil, err
	}

	if v, ok := result.value.(*GooseValue); ok {
		return v, nil
	}

	return wrap(result.value), nil
}

func (i *interpreter) evalIdent(scope *GooseScope, expr *ast.Ident) (*GooseValue, error) {
	defer un(trace(i, "ident"))

	val := scope.get(expr.Name)
	if val == nil {
		return nil, fmt.Errorf("%s is not defined", expr.Name)
	}

	return val, nil
}

// alex was here
func (i *interpreter) evalString(scope *GooseScope, expr *ast.StringLiteral) (*GooseValue, error) {
	defer un(trace(i, "string"))

	value := expr.StringStart.Content

	for _, part := range expr.Parts {
		switch expr := part.(type) {
		case *ast.StringLiteralMiddle:
			value += expr.Content
		case *ast.StringLiteralInterpIdent:
			val, err := i.evalIdent(scope, &ast.Ident{Name: expr.Name})
			if err != nil {
				return nil, err
			}
			value += toString(val.Value)
		case *ast.StringLiteralInterpExpr:
			val, err := i.evalExpr(scope, expr.Expr)
			if err != nil {
				return nil, err
			}
			value += toString(val.Value)
		default:
			return nil, fmt.Errorf("unexpected string literal part %T", expr)
		}
	}

	value += expr.StringEnd.Content

	return wrap(value), nil
}

func (i *interpreter) evalCompositeLiteral(scope *GooseScope, expr *ast.CompositeLiteral) (*GooseValue, error) {
	defer un(trace(i, "composite literal"))

	composite := make(GooseComposite)
	for _, field := range expr.Fields {
		var keyValue interface{}
		switch key := field.Key.(type) {
		case *ast.Ident:
			keyValue = key.Name
		case *ast.StringLiteral:
			k, err := i.evalString(scope, key)
			if err != nil {
				return nil, err
			}

			keyValue = k.Value
		default:
			lit, err := i.evalExpr(scope, key)
			if err != nil {
				return nil, err
			}

			switch lit.Type {
			case GooseTypeInt, GooseTypeFloat, GooseTypeBool, GooseTypeString:
				keyValue = lit.Value
			default:
				return nil, fmt.Errorf("unexpected composite literal key type %s", lit.Type)
			}
		}

		val, err := i.evalExpr(scope, field.Value)
		if err != nil {
			return nil, err
		}

		composite[keyValue] = val
	}

	return wrap(composite), nil
}

func (i *interpreter) evalArrayLiteral(scope *GooseScope, expr *ast.ArrayLiteral) (*GooseValue, error) {
	defer un(trace(i, "array lit"))

	var values []*GooseValue
	for _, value := range expr.List {
		val, err := i.evalExpr(scope, value)
		if err != nil {
			return nil, err
		}

		values = append(values, val)
	}

	return wrap(values), nil
}

func (i *interpreter) evalArrayInitializer(scope *GooseScope, expr *ast.ArrayInitializer) (*GooseValue, error) {
	defer un(trace(i, "array initializer"))

	var values []*GooseValue

	count, err := i.evalExpr(scope, expr.Count)
	if err != nil {
		return nil, err
	}

	countVal := toInt64(count.Value)

	if lit, ok := expr.Value.(*ast.Literal); ok {
		// don't bother evaluating the value if it's a literal
		value, err := i.evalLiteral(scope, lit)
		if err != nil {
			return nil, err
		}

		for i := int64(0); i < countVal; i++ {
			values = append(values, value)
		}

		return wrap(values), nil
	}

	for idx := int64(0); idx < countVal; idx++ {
		newScope := scope.fork(ScopeOwnerArrayInit)
		newScope.set("_", GooseValue{
			Constant: true,
			Type:     GooseTypeInt,
			Value:    idx,
		})

		val, err := i.evalExpr(newScope, expr.Value)
		if err != nil {
			return nil, err
		}

		values = append(values, val)
	}

	return wrap(values), nil
}

func (i *interpreter) evalSelectorExpr(scope *GooseScope, expr *ast.SelectorExpr) (*GooseValue, error) {
	defer un(trace(i, "selector expr"))

	x, err := i.evalExpr(scope, expr.X)
	if err != nil {
		return nil, err
	}

	sel := wrap(expr.Sel.Name)

	return getProperty(x, sel)
}

func (i *interpreter) evalBracketSelectorExpr(scope *GooseScope, expr *ast.BracketSelectorExpr) (*GooseValue, error) {
	defer un(trace(i, "bracket selector expr"))

	x, err := i.evalExpr(scope, expr.X)
	if err != nil {
		return nil, err
	}

	sel, err := i.evalExpr(scope, expr.Sel)
	if err != nil {
		return nil, err
	}

	return getProperty(x, sel)
}

func (i *interpreter) evalSliceExpr(scope *GooseScope, expr *ast.SliceExpr) (*GooseValue, error) {
	defer un(trace(i, "slice expression"))

	left, err := i.evalExpr(scope, expr.X)
	if err != nil {
		return nil, err
	}

	if left.Type != GooseTypeArray {
		return nil, fmt.Errorf("cannot slice non-array type %s", left.Type)
	}

	values := left.Value.([]*GooseValue)
	start := int64(0)
	end := int64(len(values))

	if expr.Low != nil {
		begin, err := i.evalExpr(scope, expr.Low)
		if err != nil {
			return nil, err
		}

		if err := expectType(begin, GooseTypeNumeric); err != nil {
			return nil, err
		}

		idx := toInt64(begin.Value)
		if idx < 0 {
			idx = int64(len(values)) + idx
			if idx < 0 {
				idx = 0
			}
			if idx >= int64(len(values)) {
				return nil, fmt.Errorf("index %d out of bounds", begin.Value)
			}
		}

		start = idx
	}

	if expr.High != nil {
		endVal, err := i.evalExpr(scope, expr.High)
		if err != nil {
			return nil, err
		}

		if err := expectType(endVal, GooseTypeNumeric); err != nil {
			return nil, err
		}

		idx := toInt64(endVal.Value)
		if idx < 0 {
			idx = int64(len(values)) + idx
			if idx < 0 {
				idx = 0
			}
			if idx >= int64(len(values)) {
				return nil, fmt.Errorf("index %d out of bounds", endVal.Value)
			}
		}
		end = idx
	}

	if end < start {
		return nil, fmt.Errorf("computed end index %d is less than computed start index %d", end, start)
	}

	if start < 0 {
		start = 0
	}

	if end > int64(len(values)) {
		end = int64(len(values))
	}

	return wrap(values[start:end]), nil
}
