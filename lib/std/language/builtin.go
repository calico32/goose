package std_language

import (
	"math/big"
	"strconv"
	"strings"

	. "github.com/calico32/goose/interpreter/lib"
	"github.com/calico32/goose/lib/types"
)

var Builtins = []types.BuiltinDoc{
	{
		Name:        "int.parse",
		Label:       "int.parse(str)",
		Signature:   "int.parse(s: str, base?: int) -> int",
		Desc:        "Parse an integer from a string.",
		Description: "Parse an integer from a string. The base can be specified as a second argument. If no base is provided, base 10 is used. If the string is not a valid integer, an error is thrown (use `int.tryParse` to avoid this).",
	},
	{
		Name:        "int.tryParse",
		Label:       "int.tryParse(str)",
		Signature:   "int.tryParse(s: str, base?: int) -> int | null",
		Desc:        "Try to parse an integer from a string.",
		Description: "Try to parse an integer from a string. The base can be specified as a second argument. If no base is provided, base 10 is used. If the string is not a valid integer, `null` is returned.",
	},
	{
		Name:        "float.parse",
		Label:       "float.parse(str)",
		Signature:   "float.parse(s: str) -> float",
		Desc:        "Parse a float from a string.",
		Description: "Parse a float from a string. If the string is not a valid float, an error is thrown (use `float.tryParse` to avoid this).",
	},
	{
		Name:        "float.tryParse",
		Label:       "float.tryParse(str)",
		Signature:   "float.tryParse(s: str) -> float | null",
		Desc:        "Try to parse a float from a string.",
		Description: "Try to parse a float from a string. If the string is not a valid float, `null` is returned.",
	},
	{
		Name:        "bool.parse",
		Label:       "bool.parse(str)",
		Signature:   "bool.parse(s: str) -> bool",
		Desc:        "Parse a boolean from a string.",
		Description: `Parse a boolean from a string. Valid values are "1", "t", "T", "TRUE", "true", "True", "0", "f", "F", "FALSE", "false", and "False". If the string is not a valid boolean, an error is thrown.`,
	},
}

var Builtin = map[string]Value{
	"F/int.parse": &Func{Executor: func(ctx *FuncContext) *Return {
		if len(ctx.Args) < 1 {
			ctx.Interp.Throw("int.parse(s): expected at least 1 argument")
			return &Return{}
		}

		s := ctx.Args[0]
		base := 10
		if len(ctx.Args) > 1 {
			if integer, ok := ctx.Args[1].(*Integer); ok {
				base = int(integer.Value.Int64())
			} else {
				ctx.Interp.Throw("int.parse(s, base): expected integer")
			}
		}

		if str, ok := s.(*String); ok {
			i := new(big.Int)
			_, ok := i.SetString(strings.TrimSpace(str.Value), base)

			if !ok {
				ctx.Interp.Throw("int.parse(s): failed to parse integer")
				return &Return{}
			}
			return NewReturn(NewInteger(i))
		} else {
			ctx.Interp.Throw("int.parse(s): expected string")
			return &Return{}
		}
	}},
	"F/int.tryParse": &Func{Executor: func(ctx *FuncContext) *Return {
		if len(ctx.Args) < 1 {
			return ReturnNull
		}

		s := ctx.Args[0]
		base := 10
		if len(ctx.Args) > 1 {
			if integer, ok := ctx.Args[1].(*Integer); ok {
				base = int(integer.Value.Int64())
			} else {
				return ReturnNull
			}
		}

		if str, ok := s.(*String); ok {
			i := new(big.Int)
			_, ok := i.SetString(strings.TrimSpace(str.Value), base)

			if !ok {
				return ReturnNull
			}
			return NewReturn(NewInteger(i))
		} else {
			return ReturnNull
		}
	}},
	"F/float.parse": &Func{Executor: func(ctx *FuncContext) *Return {
		if len(ctx.Args) < 1 {
			ctx.Interp.Throw("float.parse(s): expected at least 1 argument")
			return &Return{}
		}

		s := ctx.Args[0]
		if str, ok := s.(*String); ok {
			f, err := strconv.ParseFloat(string(str.Value), 64)
			if err != nil {
				ctx.Interp.Throw("float.parse(s): " + err.Error())
				return &Return{}
			}
			return NewReturn(f)
		} else {
			ctx.Interp.Throw("float.parse(s): expected string")
			return &Return{}
		}
	}},
	"F/float.tryParse": &Func{Executor: func(ctx *FuncContext) *Return {
		if len(ctx.Args) < 1 {
			return ReturnNull
		}

		s := ctx.Args[0]
		if str, ok := s.(*String); ok {
			f, err := strconv.ParseFloat(string(str.Value), 64)
			if err != nil {
				return ReturnNull
			}
			return NewReturn(f)
		} else {
			return ReturnNull
		}
	}},
	"F/bool.parse": &Func{Executor: func(ctx *FuncContext) *Return {
		if len(ctx.Args) < 1 {
			ctx.Interp.Throw("bool.parse(s): expected at least 1 argument")
			return &Return{}
		}

		s := ctx.Args[0]
		if str, ok := s.(*String); ok {
			b, err := strconv.ParseBool(string(str.Value))
			if err != nil {
				ctx.Interp.Throw("bool.parse(s): " + err.Error())
				return &Return{}
			}
			return NewReturn(b)
		} else {
			ctx.Interp.Throw("bool.parse(s): expected string")
			return &Return{}
		}
	}},
}
