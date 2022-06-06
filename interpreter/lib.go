package interpreter

import (
	"fmt"
	"strings"
	"time"
)

func padCommon(scope *GooseScope, args []*GooseValue) (str string, pad string, err error) {
	if len(args) < 2 {
		err = fmt.Errorf("pad__(x, int, val): expected at least 2 arguments")
		return
	}

	err = scope.interp.expectType(args[0], GooseTypeString)
	if err != nil {
		return
	}

	str = args[0].Value.(string)

	err = scope.interp.expectType(args[1], GooseTypeNumeric)
	if err != nil {
		return
	}

	length := toInt(args[1].Value)

	if length < 0 {
		err = fmt.Errorf("pad__(x, int, val): expected length >= 0")
		return
	}

	padChar := " "

	if len(args) > 2 {
		err = scope.interp.expectType(args[2], GooseTypeString)
		if err != nil {
			return
		}

		padChar = args[2].Value.(string)
	}

	if len(str) >= length {
		return
	}

	pad = ""
	for len(pad) < length-len(str) {
		pad += padChar
	}

	return str, pad, nil
}

var stdlib = map[string]GooseFunc{
	"indices": func(scope *GooseScope, args []*GooseValue) (*ReturnResult, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("indices(x): expected 1 argument")
		}

		err := scope.interp.expectType(args[0], GooseTypeArray)
		if err != nil {
			return nil, err
		}

		values := args[0].Value.([]*GooseValue)
		result := make([]*GooseValue, len(values))

		for i := range values {
			result[i] = &GooseValue{
				Type:  GooseTypeInt,
				Value: int64(i),
			}
		}

		return &ReturnResult{result}, nil
	},
	"padLeft": func(scope *GooseScope, args []*GooseValue) (*ReturnResult, error) {
		str, pad, err := padCommon(scope, args)
		if err != nil {
			return nil, err
		}

		return &ReturnResult{pad + str}, nil
	},
	"padRight": func(scope *GooseScope, args []*GooseValue) (*ReturnResult, error) {
		str, pad, err := padCommon(scope, args)
		if err != nil {
			return nil, err
		}

		return &ReturnResult{str + pad}, nil
	},
	"string": func(_ *GooseScope, args []*GooseValue) (*ReturnResult, error) {
		if len(args) == 0 {
			return nil, fmt.Errorf("string(x): expected 1 argument")
		}
		return &ReturnResult{toString(args[0].Value)}, nil
	},
	"len": func(_ *GooseScope, args []*GooseValue) (*ReturnResult, error) {
		if len(args) == 0 {
			return nil, fmt.Errorf("len(x): expected 1 argument")
		}

		switch args[0].Type {
		case GooseTypeArray:
			return &ReturnResult{int64(len(args[0].Value.([]*GooseValue)))}, nil
		case GooseTypeString:
			return &ReturnResult{int64(len(args[0].Value.(string)))}, nil
		default:
			return nil, fmt.Errorf("len(x): expected an array or string, got %s", args[0].Type)
		}
	},
	"sleep": func(scope *GooseScope, args []*GooseValue) (*ReturnResult, error) {
		var ms int64
		if len(args) == 0 {
			return nil, fmt.Errorf("sleep(x): expected 1 argument")
		}
		err := scope.interp.expectType(args[0], GooseTypeNumeric)
		if err != nil {
			return nil, err
		}
		ms = toInt64(args[0].Value)
		time.Sleep(time.Duration(ms * int64(time.Millisecond)))
		return &ReturnResult{}, nil
	},
	"milli": func(_ *GooseScope, _ []*GooseValue) (*ReturnResult, error) {
		ms := time.Now().UnixNano() / int64(time.Millisecond)
		return &ReturnResult{ms}, nil
	},
	"nano": func(_ *GooseScope, _ []*GooseValue) (*ReturnResult, error) {
		ns := time.Now().UnixNano() / int64(time.Nanosecond)
		return &ReturnResult{ns}, nil
	},
	"print": func(scope *GooseScope, args []*GooseValue) (*ReturnResult, error) {
		for _, arg := range args {
			fmt.Fprint(scope.interp.stdout, toString(arg.Value))
		}
		fmt.Fprintln(scope.interp.stdout)
		return &ReturnResult{}, nil
	},
	"printf": func(scope *GooseScope, args []*GooseValue) (*ReturnResult, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("printf(format, ...): expected at least 1 argument")
		}

		format := args[0].Value.(string)
		args = args[1:]

		var values []any

		for _, arg := range args {
			values = append(values, arg.Value)
		}

		fmt.Fprintf(scope.interp.stdout, format, values...)
		return &ReturnResult{}, nil
	},
	"exit": func(scope *GooseScope, args []*GooseValue) (*ReturnResult, error) {
		exitCode := 0
		if len(args) != 0 {
			err := scope.interp.expectType(args[0], GooseTypeNumeric)
			if err != nil {
				return nil, err
			}
			exitCode = toInt(args[0].Value)
		}
		// TODO: tinygo doesn't let you recover panics, so any exit will cause a crash
		panic(gooseExit{exitCode})
	},
	"floor": func(scope *GooseScope, args []*GooseValue) (*ReturnResult, error) {
		if len(args) == 0 {
			return nil, fmt.Errorf("floor(x): expected 1 argument")
		}
		err := scope.interp.expectType(args[0], GooseTypeNumeric)
		if err != nil {
			return nil, err
		}

		return &ReturnResult{toInt64(args[0].Value)}, nil
	},
	"ceil": func(scope *GooseScope, args []*GooseValue) (*ReturnResult, error) {
		if len(args) == 0 {
			return nil, fmt.Errorf("ceil(x): expected 1 argument")
		}
		err := scope.interp.expectType(args[0], GooseTypeNumeric)
		if err != nil {
			return nil, err
		}
		if x, ok := args[0].Value.(int64); ok {
			return &ReturnResult{x}, nil
		} else {
			return &ReturnResult{int64(x + 1)}, nil
		}

	},
	"round": func(scope *GooseScope, args []*GooseValue) (*ReturnResult, error) {
		if len(args) == 0 {
			return nil, fmt.Errorf("round(x): expected 1 argument")
		}
		err := scope.interp.expectType(args[0], GooseTypeNumeric)
		if err != nil {
			return nil, err
		}
		if x, ok := args[0].Value.(int64); ok {
			return &ReturnResult{x}, nil
		} else {
			return &ReturnResult{int64(args[0].Value.(float64) + 0.5)}, nil
		}
	},
	"join": func(scope *GooseScope, args []*GooseValue) (*ReturnResult, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("join(list, sep): expected at least 1 argument")
		}

		array := args[0]
		if err := scope.interp.expectType(array, GooseTypeArray); err != nil {
			return nil, err
		}

		values := array.Value.([]*GooseValue)

		if len(values) == 0 {
			return &ReturnResult{}, nil
		}
		if len(values) == 1 {
			return &ReturnResult{toString(values[0].Value)}, nil
		}

		var sep string
		if len(args) > 1 {
			if err := scope.interp.expectType(args[1], GooseTypeString); err != nil {
				return nil, err
			}
			sep = args[1].Value.(string)
		} else {
			sep = ","
		}

		var out strings.Builder
		for i, value := range values {
			if i > 0 {
				out.WriteString(sep)
			}
			out.WriteString(toString(value.Value))
		}

		return &ReturnResult{out.String()}, nil
	},
}

var builtins = map[string]GooseValue{
	"true":  {true, GooseTypeBool, true},
	"false": {true, GooseTypeBool, false},
	"null":  {true, GooseTypeNull, nil},
}
