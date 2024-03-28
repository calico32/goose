package std_fs

import (
	"os"

	. "github.com/calico32/goose/interpreter/lib"
	"github.com/calico32/goose/lib/types"
)

var Doc = types.StdlibDoc{
	Name:        "fs",
	Description: "Utilities for working with the file system.",
}

var Index = map[string]Value{
	"F/readFile": &Func{Executor: func(ctx *FuncContext) *Return {
		if len(ctx.Args) < 1 {
			ctx.Interp.Throw("std:fs.readFile(file): expected 1 argument")
		}
		file := ToString(ctx.Interp, ctx.Scope, ctx.Args[0])
		if file == "" {
			ctx.Interp.Throw("std:fs.readFile(file): expected string")
		}

		f, err := os.ReadFile(file)
		if err != nil {
			ctx.Interp.Throw(err.Error())
		}

		return NewReturn(string(f))
	}},
	"F/writeFile": &Func{Executor: func(ctx *FuncContext) *Return {
		if len(ctx.Args) < 2 {
			ctx.Interp.Throw("std:fs.writeFile(file): expected 2 arguments")
		}
		file := ToString(ctx.Interp, ctx.Scope, ctx.Args[0])
		content := ToString(ctx.Interp, ctx.Scope, ctx.Args[1])
		if file == "" {
			ctx.Interp.Throw("std:fs.writeFile(file): expected string")
		}

		err := os.WriteFile(file, []byte(content), 0644)
		if err != nil {
			ctx.Interp.Throw(err.Error())
		}

		return &Return{}
	}},
	"F/appendFile": &Func{Executor: func(ctx *FuncContext) *Return {
		ctx.Interp.Throw("TODO: implement std:fs/index.goose#appendFile")
		return NewReturn("TODO: implement std:fs/index.goose#appendFile")
	}},
	"F/deleteFile": &Func{Executor: func(ctx *FuncContext) *Return {
		ctx.Interp.Throw("TODO: implement std:fs/index.goose#deleteFile")
		return NewReturn("TODO: implement std:fs/index.goose#deleteFile")
	}},
}
