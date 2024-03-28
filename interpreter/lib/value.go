package lib

import (
	"fmt"
	"math/big"
)

var TrueValue = &Bool{true}
var FalseValue = &Bool{false}
var NullValue = &Null{}

var BoolFrom = map[bool]*Bool{
	true:  TrueValue,
	false: FalseValue,
}

func Wrap(value any) Value {
	switch value := value.(type) {
	case nil:
		return NullValue
	case bool:
		if value {
			return TrueValue
		} else {
			return FalseValue
		}
	case int:
		return &Integer{big.NewInt(int64(value))}
	case int64:
		return &Integer{big.NewInt(value)}
	case *big.Int:
		return &Integer{value}
	case float64:
		return &Float{value}
	case rune:
		return &String{string(value)}
	case []rune:
		return &String{string(value)}
	case byte:
		return &String{string(value)}
	case string:
		return &String{value}
	case []Value:
		return &Array{Elements: value}
	case []int:
		vals := make([]Value, len(value))
		for i, v := range value {
			vals[i] = Wrap(v)
		}
		return &Array{Elements: vals}
	case []int64:
		vals := make([]Value, len(value))
		for i, v := range value {
			vals[i] = Wrap(v)
		}
		return &Array{Elements: vals}
	case []float64:
		vals := make([]Value, len(value))
		for i, v := range value {
			vals[i] = Wrap(v)
		}
		return &Array{Elements: vals}
	case []string:
		vals := make([]Value, len(value))
		for i, v := range value {
			vals[i] = Wrap(v)
		}
		return &Array{Elements: vals}
	case [][]Value:
		vals := make([]Value, len(value))
		for i, v := range value {
			vals[i] = Wrap(v)
		}
		return &Array{Elements: vals}
	case FuncType:
		return &Func{Executor: value}
	case *Value:
		return Wrap(*value)
	case Value:
		return value
	case map[string]any:
		obj := NewComposite()
		for k, v := range value {
			SetProperty(obj, NewString(k), Wrap(v))
		}
		return obj
	default:
		panic(fmt.Errorf("unexpected type %T in wrap()", value))
	}
}
