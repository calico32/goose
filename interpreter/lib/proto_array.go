package lib

import (
	"math/big"
	"sort"
	"strings"

	"github.com/calico32/goose/token"
)

var ArrayPrototype = &Composite{
	Proto:  Object,
	Frozen: true,
	Properties: Properties{
		PKString: {
			"toString": &Func{
				Executor: func(ctx *FuncContext) *Return {
					var b strings.Builder
					b.WriteString("[")
					for i, v := range ctx.This.(*Array).Elements {
						if i > 0 {
							b.WriteString(", ")
						}
						b.WriteString(ToString(ctx.Interp, ctx.Scope, v))
					}
					b.WriteString("]")
					return NewReturn(NewString(b.String()))
				},
			},
			// "join": &Func{Executor: func(ctx *FuncContext) *Return {
			// 	if len(ctx.Args) < 1 {
			// 		ctx.Interp.Throw("join(list, sep): expected at least 1 argument")
			// 	}

			// 	array := ctx.Args[0]
			// 	if _, ok := array.(*Array); !ok {
			// 		ctx.Interp.Throw("join(list, sep): expected list to be an array")
			// 	}

			// 	values := array.(*Array).Elements

			// 	if len(values) == 0 {
			// 		return &Return{}
			// 	}
			// 	if len(values) == 1 {
			// 		return NewReturn(values[0])
			// 	}

			// 	var sep string
			// 	if len(ctx.Args) > 1 {
			// 		if _, ok := ctx.Args[1].(*String); !ok {
			// 			ctx.Interp.Throw("join(list, sep): expected sep to be a string")
			// 		}
			// 		sep = ctx.Args[1].(*String).Value
			// 	} else {
			// 		sep = ","
			// 	}

			// 	var out strings.Builder
			// 	for i, value := range values {
			// 		if i > 0 {
			// 			out.WriteString(sep)
			// 		}
			// 		out.WriteString(toString(ctx.Interp, ctx.Scope, value))
			// 	}

			// 	return NewReturn(out.String())
			// }},
			"push": &Func{Executor: func(ctx *FuncContext) *Return {
				if len(ctx.Args) < 1 {
					ctx.Interp.Throw("push(value): expected at least 1 argument")
				}

				values := ctx.This.(*Array).Elements

				if len(ctx.Args) >= 1 {
					values = append(values, ctx.Args[0:]...)
				}

				ctx.This.(*Array).Elements = values
				return NewReturn(values)
			}},
			"sort": &Func{Executor: func(ctx *FuncContext) *Return {
				if len(ctx.Args) < 1 {
					ctx.Interp.Throw("sort(func): expected at least 1 argument")
				}

				array := ctx.This.(*Array).Elements
				comparator := ctx.Args[0]
				if _, ok := comparator.(*Func); !ok {
					ctx.Interp.Throw("sort(func): expected func to be a function")
				}

				c := make([]Value, len(array))
				copy(c, array)
				sort.Slice(c, func(i, j int) bool {
					a := c[i]
					b := c[j]
					ret := comparator.(*Func).Executor(&FuncContext{
						Interp: ctx.Interp,
						Scope:  ctx.Scope,
						This:   Wrap(nil),
						Args:   []Value{a, b},
					})

					if ret == nil {
						ctx.Interp.Throw("sort(func): expected function to return a value")
						panic(0)
					}

					if _, ok := ret.Value.(*Integer); !ok {
						ctx.Interp.Throw("sort(func): expected function to return a number")
					}

					return ret.Value.(*Integer).Value.Cmp(big.NewInt(0)) == -1
				})

				return NewReturn(c)
			}},
			"map": &Func{Executor: func(ctx *FuncContext) *Return {
				if len(ctx.Args) < 1 {
					ctx.Interp.Throw("map(func): expected at least 1 argument")
				}

				array := ctx.This.(*Array).Elements
				mapper := ctx.Args[0]
				if _, ok := mapper.(*Func); !ok {
					ctx.Interp.Throw("map(func): expected func to be a function")
				}

				c := make([]Value, len(array))
				for i, v := range array {
					ret := mapper.(*Func).Executor(&FuncContext{
						Interp: ctx.Interp,
						Scope:  ctx.Scope,
						This:   Wrap(nil),
						Args:   []Value{v, Wrap(i)},
					})

					c[i] = ret.Value
				}

				return NewReturn(c)
			}},
			"join": &Func{Executor: func(ctx *FuncContext) *Return {
				if len(ctx.Args) < 1 {
					ctx.Interp.Throw("join(sep): expected at least 1 argument")
				}

				array := ctx.This.(*Array).Elements
				sep := ctx.Args[0].(*String).Value

				var b strings.Builder
				for i, v := range array {
					if i > 0 {
						b.WriteString(sep)
					}
					b.WriteString(ToString(ctx.Interp, ctx.Scope, v))
				}

				return NewReturn(b.String())
			}},
		},
	},
	Operators: Operators{
		token.Eq: OpFunc(func(c *OpContext[*Array, Value]) Value {
			switch o := c.Other.(type) {
			case *Array:
				if len(c.This.Elements) != len(o.Elements) {
					return FalseValue
				}
				for i, e := range c.This.Elements {
					op := GetOperator(e, token.Eq)
					if op == nil {
						c.Interp.Throw("array element is not comparable")
						panic(0) // static analysis
					}
					ret := op.Executor(&FuncContext{
						Interp: c.Interp,
						Scope:  c.Scope,
						This:   e,
						Args:   []Value{o.Elements[i]},
					})
					if !ret.Value.(*Bool).Value {
						return FalseValue
					}
				}
				return TrueValue
			default:
				return FalseValue
			}
		}),
	},
}
