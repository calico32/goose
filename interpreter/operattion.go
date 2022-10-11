package interpreter

import (
	"fmt"
	"math"

	"github.com/wiisportsresort/goose/token"
)

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
		panic(fmt.Errorf("unexpected lhs type %T in numericOpInt()", lhs))
	}
	switch rhs.(type) {
	case int64:
	case float64:
		bothInts = false
	default:
		panic(fmt.Errorf("unexpected rhs type %T in numericOpInt()", rhs))
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
			panic(fmt.Errorf("unexpected operator %s", op))
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
			panic(fmt.Errorf("unexpected operator %s", op))
		}
	}
}

func stringOp(lhs string, op token.Token, rhs string) string {
	switch op {
	case token.Add, token.AddAssign:
		return lhs + rhs
	default:
		panic(fmt.Errorf("unexpected operator %s", op))
	}
}
