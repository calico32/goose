package interpreter

import (
	"fmt"
	"strings"
)

func toInt(value any) int {
	switch value := value.(type) {
	case *GooseValue:
		return toInt(value.Value)
	case int64:
		return int(value)
	case float64:
		return int(value)
	default:
		panic(fmt.Errorf("unexpected type %T in toInt()", value))
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
		panic(fmt.Errorf("unexpected type %T in toInt64()", value))
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
		panic(fmt.Errorf("unexpected type %T in toInt()", value))
	}
}

func toString(x any) string {
	switch x := x.(type) {
	case *GooseValue:
		if x == nil {
			return "<nil>"
		}
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
	case GooseComposite:
		var output strings.Builder
		output.WriteString("{ ")
		i := 0
		for k, v := range x {
			if i > 0 {
				output.WriteString(", ")
			}
			output.WriteString(toString(k))
			output.WriteString(": ")
			output.WriteString(toString(v))
			i++
		}
		output.WriteString(" }")
		return output.String()
	// case int64, float64, string, bool, nil:
	default:
		return fmt.Sprintf("%v", x)
	}
}
