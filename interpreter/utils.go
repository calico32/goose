package interpreter

import (
	"fmt"
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
	case GooseComposite:
		return len(value) > 0
	default:
		panic(fmt.Errorf("unexpected type %T in isTruthy()", value))
	}
}

func valueOf(value any) any {
	switch value := value.(type) {
	case *GooseValue:
		return valueOf(value.Value)
	case int: // convert ints to int64s
		return int64(value)
	case rune: // convert runes to strings
		return string(value)
	default:
		return value
	}
}

func typeOf(value any) GooseType {
	switch value := value.(type) {
	case *GooseValue:
		return value.Type
	case int, int64:
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
	case GooseArray:
		return GooseTypeArray
	case GooseComposite:
		return GooseTypeComposite
	default:
		panic(fmt.Errorf("unexpected type %T in typeOf()", value))
	}
}

func isPowerOfTwo(v int) bool {
	return v != 0 && (v&(v-1)) == 0
}
