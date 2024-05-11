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
		Signature:   "len(x: array | string) -> int",
		Desc:        "Get the length of an array or string.",
		Description: "Get the length of an array or string. If the argument is an array, the length is the number of elements in the array. If the argument is a string, the length is the number of characters in the string.",
		Examples: []types.CodeSnippet{
			{
				Content: `<pre>
				const a = [1, 2, 3]
				println(len(a)) // 3
				println(len("hello")) // 5
				</pre>`,
			},
		},
	},
	{
		Name:        "print",
		Label:       "print(...)",
		Signature:   "print(...obj: any[]) -> void",
		Desc:        "Print values to the standard output.",
		Description: "Print values to the standard output. If multiple arguments are provided, they are separated by spaces.\n\nArguments are converted to strings using their `toString` method. No newline is printed at the end.",
		Examples: []types.CodeSnippet{
			{
				Content: `<pre>
				print("hello world") // hello world
				print(123) // 123
				print("hello", "world") // hello world
				print(1, 2, 3) // 1 2 3
				</pre>`,
			},
		},
	},
	{
		Name:        "println",
		Label:       "println(...)",
		Signature:   "println(...obj: any[]) -> void",
		Desc:        "Print values to the standard output with a newline.",
		Description: "Print values to the standard output with a newline. If multiple arguments are provided, they are separated by spaces.\n\nArguments are converted to strings using their `toString` method.",
		Examples: []types.CodeSnippet{
			{
				Content: `<pre>
				println("hello world") // hello world
				println(123) // 123
				println("hello", "world") // hello world
				println(1, 2, 3) // 1 2 3
				</pre>`,
			},
		},
	},
	{
		Name:      "printf",
		Label:     "printf(format, ...)",
		Signature: "printf(format: string, ...obj: any[]) -> void",
		Desc:      "Print values to the standard output with a format string.",
		Examples: []types.CodeSnippet{
			{
				Content: `<pre>
				printf("Hello, %s!\n", "world") // Hello, world!
				printf("The answer is %d\n", 42) // The answer is 42
				</pre>`,
			},
		},
	},
	{
		Name:      "exit",
		Label:     "exit(code)",
		Signature: "exit(code: int[>=0,<=255]) -> never",
		Desc:      "Exit the program with an exit code.",
		Examples: []types.CodeSnippet{
			{
				Content: `<pre>
				exit(0) // Exit with code 0
				exit(1) // Exit with code 1
				</pre>`,
			},
		},
	},
	{
		Name:        "typeof",
		Label:       "typeof(value)",
		Signature:   "typeof(value: any) -> string",
		Desc:        "Get the type of a value.",
		Description: "Get the type of a value. The return value is a string representing the type of the value. The possible types are:\n\n- `int`\n- `float`\n- `string`\n- `bool`\n- `array`\n- `object`\n- `function`\n- `null`",
		Examples: []types.CodeSnippet{
			{
				Content: `<pre>
				println(typeof(123)) // int
				println(typeof("hello")) // string
				println(typeof([1, 2, 3])) // array
				</pre>`,
			},
		},
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
