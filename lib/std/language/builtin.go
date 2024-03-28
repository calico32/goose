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
		Description: "Parse an integer from a string.",
	},
	{
		Name:        "int.tryParse",
		Label:       "int.tryParse(str)",
		Description: "Try to parse an integer from a string.",
	},
	{
		Name:        "float.parse",
		Label:       "float.parse(str)",
		Description: "Parse a float from a string.",
	},
	{
		Name:        "float.tryParse",
		Label:       "float.tryParse(str)",
		Description: "Try to parse a float from a string.",
	},
	{
		Name:        "bool.parse",
		Label:       "bool.parse(str)",
		Description: "Parse a boolean from a string.",
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
