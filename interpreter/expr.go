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
	case *ast.ArrayLiteral:
		return i.evalArrayLiteral(scope, expr)
	case *ast.ArrayInitializer:
		return i.evalArrayInitializer(scope, expr)
	case *ast.IndexExpr:
		return i.evalIndexExpr(scope, expr)
	case *ast.SliceExpr:
		return i.evalSliceExpr(scope, expr)
	case *ast.Literal:
		return i.evalLiteral(scope, expr)
	default:
		if badExpr, ok := expr.(*ast.BadExpr); ok {
			return nil, fmt.Errorf("Unexpected bad expression %#v", badExpr)
		}
		return nil, fmt.Errorf("Unexpected expression type %T", expr)
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

		return &GooseValue{
			Type:  GooseTypeInt,
			Value: val,
		}, nil

	case token.Float:
		val, err := strconv.ParseFloat(expr.Value, 64)
		if err != nil {
			return nil, err
		}

		return &GooseValue{
			Type:  GooseTypeFloat,
			Value: val,
		}, nil

	case token.Null:
		return &GooseValue{
			Type:  GooseTypeNull,
			Value: nil,
		}, nil

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

	var result bool
	switch op {
	case token.Assign:
		return right, nil
	case token.Add, token.AddAssign:
		if left.Type == GooseTypeArray {
			values := left.Value.([]*GooseValue)
			values = append(values, right)
			return &GooseValue{
				Constant: false,
				Type:     GooseTypeArray,
				Value:    values,
			}, nil
		} else if left.Type == GooseTypeString || right.Type == GooseTypeString {
			return &GooseValue{
				Type:  GooseTypeString,
				Value: fmt.Sprintf("%v", left.Value) + fmt.Sprintf("%v", right.Value),
			}, nil
		}
		fallthrough
	case token.Lt, token.Lte, token.Gt, token.Gte,
		token.Sub, token.Mul, token.Quo, token.Rem, token.Pow,
		token.SubAssign, token.MulAssign, token.QuoAssign, token.RemAssign, token.PowAssign,
		token.Inc, token.Dec:
		return i.numericOperation(left, op, right)
	case token.Eq:
		result = left.Type == right.Type && left.Value == right.Value
	case token.Neq:
		result = left.Type != right.Type || left.Value != right.Value
	case token.LogAnd:
		result = isTruthy(left.Value) && isTruthy(right.Value)
	case token.LogOr:
		result = isTruthy(left.Value) || isTruthy(right.Value)
	default:
		return nil, fmt.Errorf("unexpected binary operator %s", op)
	}

	return &GooseValue{
		Constant: false,
		Type:     GooseTypeBool,
		Value:    result,
	}, nil
}

func (i *interpreter) evalUnaryExpr(scope *GooseScope, expr *ast.UnaryExpr) (*GooseValue, error) {
	defer un(trace(i, "unary expr"))

	value, err := i.evalExpr(scope, expr.X)
	if err != nil {
		return nil, err
	}

	switch expr.Op {
	case token.LogNot:
		return &GooseValue{
			Constant: false,
			Type:     GooseTypeBool,
			Value:    !isTruthy(value.Value),
		}, nil
	case token.Add:
		err = i.expectType(value, GooseTypeNumeric)
		if err != nil {
			return nil, err
		}
		return &GooseValue{
			Constant: false,
			Type:     value.Type,
			Value:    value.Value,
		}, nil
	case token.Sub:
		err = i.expectType(value, GooseTypeNumeric)
		if err != nil {
			return nil, err
		}

		if value.Type == GooseTypeInt {
			return &GooseValue{
				Constant: false,
				Type:     GooseTypeInt,
				Value:    -value.Value.(int64),
			}, nil
		}

		return &GooseValue{
			Constant: false,
			Type:     GooseTypeFloat,
			Value:    -value.Value.(float64),
		}, nil

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

	return &GooseValue{
		Constant: false,
		Type:     typeOf(result.value),
		Value:    result.value,
	}, nil
}

func (i *interpreter) evalIdent(scope *GooseScope, expr *ast.Ident) (*GooseValue, error) {
	defer un(trace(i, "ident"))

	val := scope.get(expr.Name)
	if val == nil {
		return nil, fmt.Errorf("%s is not defined", expr.Name)
	}

	return val, nil
}

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

	return &GooseValue{
		Constant: false,
		Type:     GooseTypeString,
		Value:    value,
	}, nil
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

	return &GooseValue{
		Constant: false,
		Type:     GooseTypeArray,
		Value:    values,
	}, nil
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

		return &GooseValue{
			Constant: false,
			Type:     GooseTypeArray,
			Value:    values,
		}, nil
	}

	for idx := int64(0); idx < countVal; idx++ {
		newScope := scope.new(ScopeOwnerArrayInit)
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

	return &GooseValue{
		Constant: false,
		Type:     GooseTypeArray,
		Value:    values,
	}, nil
}

func (i *interpreter) evalIndexExpr(scope *GooseScope, expr *ast.IndexExpr) (*GooseValue, error) {
	defer un(trace(i, "index expression"))

	left, err := i.evalExpr(scope, expr.X)
	if err != nil {
		return nil, err
	}

	right, err := i.evalExpr(scope, expr.Index)
	if err != nil {
		return nil, err
	}

	if left.Type != GooseTypeArray {
		return nil, fmt.Errorf("cannot index non-array type %s", left.Type)
	}

	if err := i.expectType(right, GooseTypeNumeric); err != nil {
		return nil, err
	}

	values := left.Value.([]*GooseValue)

	idx := right.Value.(int64)

	if idx >= int64(len(values)) {
		return nil, fmt.Errorf("index %d out of bounds", idx)
	}

	if idx < 0 {
		idx = int64(len(values)) + idx
		if idx < 0 {
			idx = 0
		}
		if idx >= int64(len(values)) {
			return nil, fmt.Errorf("index %d out of bounds", right.Value.(int64))
		}
	}

	return values[idx], nil
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

		if err := i.expectType(begin, GooseTypeNumeric); err != nil {
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

		if err := i.expectType(endVal, GooseTypeNumeric); err != nil {
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

	return &GooseValue{
		Constant: false,
		Type:     GooseTypeArray,
		Value:    values[start:end],
	}, nil
}
