package interpreter

import "fmt"

var TrueValue = &Bool{true}
var FalseValue = &Bool{false}
var NullValue = &Null{}

var BoolFrom = map[bool]*Bool{
	true:  TrueValue,
	false: FalseValue,
}

func wrap(value any) Value {
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
		return &Integer{int64(value)}
	case int64:
		return &Integer{value}
	case float64:
		return &Float{value}
	case string:
		return &String{value}
	case []Value:
		return &Array{value}
	case []int:
		vals := make([]Value, len(value))
		for i, v := range value {
			vals[i] = wrap(v)
		}
		return &Array{vals}
	case []int64:
		vals := make([]Value, len(value))
		for i, v := range value {
			vals[i] = wrap(v)
		}
		return &Array{vals}
	case []float64:
		vals := make([]Value, len(value))
		for i, v := range value {
			vals[i] = wrap(v)
		}
		return &Array{vals}
	case []string:
		vals := make([]Value, len(value))
		for i, v := range value {
			vals[i] = wrap(v)
		}
		return &Array{vals}
	case [][]Value:
		vals := make([]Value, len(value))
		for i, v := range value {
			vals[i] = wrap(v)
		}
		return &Array{vals}
	case FuncType:
		return &Func{Executor: value}
	default:
		panic(fmt.Errorf("unexpected type %T in wrap()", value))
	}
}
