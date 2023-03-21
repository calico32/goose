package interpreter

import (
	"fmt"
	"time"
)

var globals = map[string]FuncType{
	"indices": func(ctx *FuncContext) *Return {
		if len(ctx.Args) != 1 {
			ctx.interp.throw("indices(x): expected 1 argument")
		}

		if _, ok := ctx.Args[0].(*Array); !ok {
			ctx.interp.throw("indices(x): expected array as argument")
		}

		values := ctx.Args[0].(*Array).Elements
		result := make([]Value, len(values))

		for i := range values {
			result[i] = wrap(int64(i))
		}

		return &Return{wrap(values)}
	},
	"string": func(ctx *FuncContext) *Return {
		if len(ctx.Args) == 0 {
			ctx.interp.throw("string(x): expected 1 argument")
		}
		return &Return{wrap(toString(ctx.interp, ctx.Scope, ctx.Args[0]))}
	},
	"len": func(ctx *FuncContext) *Return {
		if len(ctx.Args) == 0 {
			ctx.interp.throw("len(x): expected 1 argument")
		}

		switch v := ctx.Args[0].(type) {
		case *Array:
			return &Return{wrap(int64(len(v.Elements)))}
		case *String:
			return &Return{wrap(int64(len(v.Value)))}
		default:
			ctx.interp.throw("len(x): expected an array or string, got %s", ctx.Args[0].Type())
			return nil
		}
	},
	"sleep": func(ctx *FuncContext) *Return {
		var ms int64
		if len(ctx.Args) == 0 {
			ctx.interp.throw("sleep(x): expected 1 argument")
		}
		if _, ok := ctx.Args[0].(Numeric); !ok {
			ctx.interp.throw("sleep(x): expected integer as first argument")
		}
		ms = ctx.Args[0].(Numeric).Int64()

		time.Sleep(time.Duration(ms * int64(time.Millisecond)))
		return &Return{NullValue}
	},
	"print": func(ctx *FuncContext) *Return {
		for i, arg := range ctx.Args {
			fmt.Fprint(ctx.interp.stdout, toString(ctx.interp, ctx.Scope, arg))
			if i < len(ctx.Args)-1 {
				fmt.Fprint(ctx.interp.stdout, " ")
			}
		}
		fmt.Fprintln(ctx.interp.stdout)
		return &Return{}
	},
	"printf": func(ctx *FuncContext) *Return {
		if len(ctx.Args) < 1 {
			ctx.interp.throw("printf(format, ...): expected at least 1 argument")
		}

		if _, ok := ctx.Args[0].(*String); !ok {
			ctx.interp.throw("printf(format, ...): expected string as first argument")
		}

		format := ctx.Args[0].(*String).Value
		ctx.Args = ctx.Args[1:]

		var values []any

		for _, arg := range ctx.Args {
			values = append(values, arg.Unwrap())
		}

		fmt.Fprintf(ctx.interp.stdout, format, values...)
		return &Return{}
	},
	"exit": func(ctx *FuncContext) *Return {
		exitCode := 0
		if len(ctx.Args) != 0 {
			if _, ok := ctx.Args[0].(Numeric); !ok {
				ctx.interp.throw("exit(code): expected integer as first argument")
			}
			exitCode = ctx.Args[0].(Numeric).Int()
		}
		// TODO: tinygo doesn't let you recover panics, so any exit will cause a crash
		panic(gooseExit{exitCode})
	},

	"keys": func(ctx *FuncContext) *Return {
		if len(ctx.Args) < 1 {
			ctx.interp.throw("keys(composite): expected 1 argument")
		}

		composite := ctx.Args[0]
		if _, ok := composite.(*Composite); !ok {
			ctx.interp.throw("keys(composite): expected composite as first argument")
		}

		keys := []Value{}
		for k := range composite.(*Composite).Properties[PropertyKeyString] {
			keys = append(keys, wrap(k))
		}
		for k := range composite.(*Composite).Properties[PropertyKeyInteger] {
			keys = append(keys, wrap(k))
		}

		return &Return{wrap(keys)}
	},
	"values": func(ctx *FuncContext) *Return {
		if len(ctx.Args) < 1 {
			ctx.interp.throw("values(composite): expected 1 argument")
		}

		composite := ctx.Args[0]
		if _, ok := composite.(*Composite); !ok {
			ctx.interp.throw("keys(composite): expected composite as first argument")
		}

		values := []Value{}
		for _, v := range composite.(*Composite).Properties[PropertyKeyString] {
			values = append(values, wrap(v))
		}
		for _, v := range composite.(*Composite).Properties[PropertyKeyInteger] {
			values = append(values, wrap(v))
		}

		return &Return{wrap(values)}
	},
	"typeof": func(ctx *FuncContext) *Return {
		if len(ctx.Args) < 1 {
			ctx.interp.throw("typeof(value): expected 1 argument")
		}

		return &Return{wrap(ctx.Args[0].Type())}
	},
}

var builtins = map[string]*Variable{
	"true":  {Constant: true, Value: TrueValue},
	"false": {Constant: true, Value: FalseValue},
	"null":  {Constant: true, Value: NullValue},
}
