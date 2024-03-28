package std_random

import (
	"math/big"
	"math/rand"

	"github.com/calico32/goose/lib/types"

	. "github.com/calico32/goose/interpreter/lib"
)

var Doc = types.StdlibDoc{
	Name:        "random",
	Description: "A library for generating random numbers.",
}

var Index = map[string]Value{
	"F/choice": &Func{Executor: func(ctx *FuncContext) *Return {
		if len(ctx.Args) != 1 {
			ctx.Interp.Throw("choice() expects exactly 1 argument")
		}
		if _, ok := ctx.Args[0].(*Array); !ok {
			ctx.Interp.Throw("choice() expects an array as its argument")
		}
		els := ctx.Args[0].(*Array).Elements
		return &Return{Value: els[rand.Intn(len(els))]}
	}},
	"F/int": &Func{Executor: func(ctx *FuncContext) *Return {
		if len(ctx.Args) != 2 {
			ctx.Interp.Throw("int() expects exactly 2 arguments")
		}
		if _, ok := ctx.Args[0].(*Integer); !ok {
			ctx.Interp.Throw("int() expects an integer as its first argument")
		}
		if _, ok := ctx.Args[1].(*Integer); !ok {
			ctx.Interp.Throw("int() expects an integer as its second argument")
		}
		min := ctx.Args[0].(*Integer).Value
		max := ctx.Args[1].(*Integer).Value
		if min.Cmp(max) > 0 {
			ctx.Interp.Throw("int() expects the first argument to be less than or equal to the second argument")
		}
		d := new(big.Int)
		d.Sub(max, min)
		v := rand.Int63n(d.Int64()) + min.Int64()
		return &Return{Value: &Integer{Value: big.NewInt(v)}}
	}},
	"F/sample": &Func{Executor: func(ctx *FuncContext) *Return {
		if len(ctx.Args) != 2 {
			ctx.Interp.Throw("sample() expects exactly 2 arguments")
		}
		if _, ok := ctx.Args[0].(*Array); !ok {
			ctx.Interp.Throw("sample() expects an array as its first argument")
		}
		if _, ok := ctx.Args[1].(*Integer); !ok {
			ctx.Interp.Throw("sample() expects an integer as its second argument")
		}
		els := ctx.Args[0].(*Array).Elements
		n := ctx.Args[1].(*Integer).Value
		if n.Int64() > int64(len(els)) {
			ctx.Interp.Throw("sample() expects the second argument to be less than or equal to the length of the first argument")
		}
		rand.Shuffle(len(els), func(i, j int) {
			els[i], els[j] = els[j], els[i]
		})
		return &Return{Value: &Array{Elements: els[:n.Int64()]}}
	}},
	"F/uniform": &Func{Executor: func(ctx *FuncContext) *Return {
		return &Return{Value: &Float{Value: rand.Float64()}}
	}},
}
