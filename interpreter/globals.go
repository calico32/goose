package interpreter

import (
	"fmt"
	"math/big"

	. "github.com/calico32/goose/interpreter/lib"
	"github.com/calico32/goose/lib/types"
)

var GlobalDocs = []types.BuiltinDoc{
	{
		Name:        "len",
		Label:       "len(x)",
		Description: "Get the length of an array or string.",
	},
	{
		Name:        "print",
		Label:       "print(...)",
		Description: "Print values to the standard output.",
	},
	{
		Name:        "println",
		Label:       "println(...)",
		Description: "Print values to the standard output with a newline.",
	},
	{
		Name:        "printf",
		Label:       "printf(format, ...)",
		Description: "Print values to the standard output with a format string.",
	},
	{
		Name:        "exit",
		Label:       "exit(code)",
		Description: "Exit the program with an exit code.",
	},
	{
		Name:        "typeof",
		Label:       "typeof(value)",
		Description: "Get the type of a value.",
	},
}

var Globals = map[string]FuncType{
	"len": func(ctx *FuncContext) *Return {
		if len(ctx.Args) == 0 {
			ctx.Interp.Throw("len(x): expected 1 argument")
		}

		switch v := ctx.Args[0].(type) {
		case *Array:
			return NewReturn(NewInteger(big.NewInt(int64(len(v.Elements)))))
		case *String:
			return NewReturn(NewInteger(big.NewInt(int64(len(v.Value)))))
		default:
			ctx.Interp.Throw("len(x): expected an array or string, got %s", ctx.Args[0].Type())
			return nil
		}
	},
	// "sleep": func(ctx *FuncContext) *Return {
	// 	var ms int64
	// 	if len(ctx.Args) == 0 {
	// 		ctx.Interp.Throw("sleep(x): expected 1 argument")
	// 	}
	// 	if _, ok := ctx.Args[0].(Numeric); !ok {
	// 		ctx.Interp.Throw("sleep(x): expected integer as first argument")
	// 	}
	// 	ms = ctx.Args[0].(Numeric).Int64()

	// 	time.Sleep(time.Duration(ms * int64(time.Millisecond)))
	// 	return NewReturn(NullValue)
	// },
	"print": func(ctx *FuncContext) *Return {
		for i, arg := range ctx.Args {
			fmt.Fprint(ctx.Interp.Stdout(), ToString(ctx.Interp, ctx.Scope, arg))
			if i < len(ctx.Args)-1 {
				fmt.Fprint(ctx.Interp.Stdout(), " ")
			}
		}
		return &Return{}
	},
	"println": func(ctx *FuncContext) *Return {
		for i, arg := range ctx.Args {
			fmt.Fprint(ctx.Interp.Stdout(), ToString(ctx.Interp, ctx.Scope, arg))
			if i < len(ctx.Args)-1 {
				fmt.Fprint(ctx.Interp.Stdout(), " ")
			}
		}
		fmt.Fprintln(ctx.Interp.Stdout())
		return &Return{}
	},
	"printf": func(ctx *FuncContext) *Return {
		if len(ctx.Args) < 1 {
			ctx.Interp.Throw("printf(format, ...): expected at least 1 argument")
		}

		if _, ok := ctx.Args[0].(*String); !ok {
			ctx.Interp.Throw("printf(format, ...): expected string as first argument")
		}

		format := ctx.Args[0].(*String).Value
		ctx.Args = ctx.Args[1:]

		var values []any

		for _, arg := range ctx.Args {
			values = append(values, arg.Unwrap())
		}

		fmt.Fprintf(ctx.Interp.Stdout(), format, values...)
		return &Return{}
	},
	"exit": func(ctx *FuncContext) *Return {
		exitCode := 0
		if len(ctx.Args) != 0 {
			if _, ok := ctx.Args[0].(Numeric); !ok {
				ctx.Interp.Throw("exit(code): expected integer as first argument")
			}
			exitCode = ctx.Args[0].(Numeric).Int()
		}
		// TODO: tinygo doesn't let you recover panics, so any exit will cause a crash
		panic(gooseExit{exitCode})
	},
	"typeof": func(ctx *FuncContext) *Return {
		if len(ctx.Args) < 1 {
			ctx.Interp.Throw("typeof(value): expected 1 argument")
		}

		return NewReturn(ctx.Args[0].Type())
	},
}

var GlobalConstants = map[string]*Variable{
	"true":  {Constant: true, Value: TrueValue},
	"false": {Constant: true, Value: FalseValue},
	"null":  {Constant: true, Value: NullValue},
}
