package interpreter

import (
	"fmt"
	"math"
	"strings"

	"github.com/wiisportsresort/goose/token"
)

func isTruthy(value any) bool {
	switch value := value.(type) {
	case *GooseValue:
		return isTruthy(value.Value)
	case string:
		return value != ""
	case int64, float64:
		return value != 0
	case bool:
		return value
	case nil:
		return false
	case GooseFunc:
		return true
	case []any:
		return len(value) > 0
	default:
		panic(fmt.Errorf("Unexpected type %T in isTruthy()", value))
	}
}

func valueOf(value any) any {
	switch value := value.(type) {
	case *GooseValue:
		return valueOf(value.Value)
	default:
		return value
	}
}

func typeOf(value any) GooseType {
	switch value := value.(type) {
	case *GooseValue:
		return value.Type
	case int64:
		return GooseTypeInt
	case float64:
		return GooseTypeFloat
	case bool:
		return GooseTypeBool
	case nil:
		return GooseTypeNull
	case GooseFunc:
		return GooseTypeFunc
	case string:
		return GooseTypeString
	case []*GooseValue:
		return GooseTypeArray
	default:
		panic(fmt.Errorf("Unexpected type %T in typeOf()", value))
	}
}

func toInt(value any) int {
	switch value := value.(type) {
	case *GooseValue:
		return toInt(value.Value)
	case int64:
		return int(value)
	case float64:
		return int(value)
	default:
		panic(fmt.Errorf("Unexpected type %T in toInt()", value))
	}
}

func toInt64(value any) int64 {
	switch value := value.(type) {
	case *GooseValue:
		return toInt64(value.Value)
	case int64:
		return int64(value)
	case float64:
		return int64(value)
	default:
		panic(fmt.Errorf("Unexpected type %T in toInt()", value))
	}
}

func toFloat(value any) float64 {
	switch value := value.(type) {
	case *GooseValue:
		return toFloat(value.Value)
	case int64:
		return float64(value)
	case float64:
		return value
	default:
		panic(fmt.Errorf("Unexpected type %T in toInt()", value))
	}
}

func toString(x any) string {
	switch x := x.(type) {
	case *GooseValue:
		return toString(x.Value)
	case []*GooseValue:
		var output strings.Builder
		output.WriteString("[")
		for i, v := range x {
			if i > 0 {
				output.WriteString(", ")
			}
			output.WriteString(toString(v))
		}
		output.WriteString("]")
		return output.String()
	case GooseFunc:
		return fmt.Sprintf("<function %v>", x)
	// case int64, float64, string, bool, nil:
	default:
		return fmt.Sprintf("%v", x)
	}
}

func isPowerOfTwo(v int) bool {
	return v != 0 && (v&(v-1)) == 0
}

var numericOpIntOutputType = map[token.Token]GooseType{
	token.Add:       GooseTypeInt,
	token.AddAssign: GooseTypeInt,
	token.Sub:       GooseTypeInt,
	token.SubAssign: GooseTypeInt,
	token.Mul:       GooseTypeInt,
	token.MulAssign: GooseTypeInt,
	token.Quo:       GooseTypeFloat,
	token.QuoAssign: GooseTypeFloat,
	token.Pow:       GooseTypeFloat,
	token.PowAssign: GooseTypeFloat,
	token.Rem:       GooseTypeInt,
	token.RemAssign: GooseTypeInt,
}

func numericOp(lhs any, op token.Token, rhs any) any {
	bothInts := true
	switch lhs.(type) {
	case int64:
	case float64:
		bothInts = false
	default:
		panic(fmt.Errorf("Unexpected lhs type %T in numericOpInt()", lhs))
	}
	switch rhs.(type) {
	case int64:
	case float64:
		bothInts = false
	default:
		panic(fmt.Errorf("Unexpected rhs type %T in numericOpInt()", rhs))
	}
	outputType := GooseTypeFloat
	if bothInts {
		outputType = numericOpIntOutputType[op]
	}

	if outputType == GooseTypeFloat {
		lhs := toFloat(lhs)
		rhs := toFloat(rhs)
		switch op {
		case token.Lt:
			return lhs < rhs
		case token.Lte:
			return lhs <= rhs
		case token.Gt:
			return lhs > rhs
		case token.Gte:
			return lhs >= rhs
		case token.Add, token.AddAssign, token.Inc:
			return lhs + rhs
		case token.Sub, token.SubAssign, token.Dec:
			return lhs - rhs
		case token.Mul, token.MulAssign:
			return lhs * rhs
		case token.Quo, token.QuoAssign:
			return lhs / rhs
		case token.Rem, token.RemAssign:
			return int64(lhs) % int64(rhs)
		case token.Pow, token.PowAssign:
			return math.Pow(lhs, rhs)
		default:
			panic(fmt.Errorf("Unexpected operator %s", op))
		}
	} else {
		lhs := toInt64(lhs)
		rhs := toInt64(rhs)
		switch op {
		case token.Lt:
			return lhs < rhs
		case token.Lte:
			return lhs <= rhs
		case token.Gt:
			return lhs > rhs
		case token.Gte:
			return lhs >= rhs
		case token.Add, token.AddAssign, token.Inc:
			return lhs + rhs
		case token.Sub, token.SubAssign, token.Dec:
			return lhs - rhs
		case token.Mul, token.MulAssign:
			return lhs * rhs
		case token.Quo, token.QuoAssign:
			return lhs / rhs
		case token.Rem, token.RemAssign:
			return lhs % rhs
		case token.Pow, token.PowAssign:
			return math.Pow(float64(lhs), float64(rhs))
		default:
			panic(fmt.Errorf("Unexpected operator %s", op))
		}
	}
}

func stringOp(lhs string, op token.Token, rhs string) string {
	switch op {
	case token.Add, token.AddAssign:
		return lhs + rhs
	default:
		panic(fmt.Errorf("Unexpected operator %s", op))
	}
}
