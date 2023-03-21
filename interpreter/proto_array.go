package interpreter

import (
	"sort"
	"strings"

	"github.com/calico32/goose/token"
)

var ArrayPrototype = &Composite{
	Proto:  Object,
	Frozen: true,
	Properties: Properties{
		PropertyKeyString: {
			"toString": &Func{
				Executor: func(ctx *FuncContext) *Return {
					var b strings.Builder
					b.WriteString("[")
					for i, v := range ctx.This.(*Array).Elements {
						if i > 0 {
							b.WriteString(", ")
						}
						b.WriteString(toString(ctx.interp, ctx.Scope, v))
					}
					b.WriteString("]")
					return &Return{&String{b.String()}}
				},
			},
			// "join": &Func{Executor: func(ctx *FuncContext) *Return {
			// 	if len(ctx.Args) < 1 {
			// 		ctx.interp.throw("join(list, sep): expected at least 1 argument")
			// 	}

			// 	array := ctx.Args[0]
			// 	if _, ok := array.(*Array); !ok {
			// 		ctx.interp.throw("join(list, sep): expected list to be an array")
			// 	}

			// 	values := array.(*Array).Elements

			// 	if len(values) == 0 {
			// 		return &Return{}
			// 	}
			// 	if len(values) == 1 {
			// 		return &Return{values[0]}
			// 	}

			// 	var sep string
			// 	if len(ctx.Args) > 1 {
			// 		if _, ok := ctx.Args[1].(*String); !ok {
			// 			ctx.interp.throw("join(list, sep): expected sep to be a string")
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
			// 		out.WriteString(toString(ctx.interp, ctx.Scope, value))
			// 	}

			// 	return &Return{wrap(out.String())}
			// }},
			"push": &Func{Executor: func(ctx *FuncContext) *Return {
				if len(ctx.Args) < 1 {
					ctx.interp.throw("push(value): expected at least 1 argument")
				}

				values := ctx.This.(*Array).Elements

				if len(ctx.Args) >= 1 {
					values = append(values, ctx.Args[0:]...)
				}

				ctx.This.(*Array).Elements = values
				return &Return{wrap(values)}
			}},
			"sort": &Func{Executor: func(ctx *FuncContext) *Return {
				if len(ctx.Args) < 1 {
					ctx.interp.throw("sort(func): expected at least 1 argument")
				}

				array := ctx.This.(*Array).Elements
				comparator := ctx.Args[0]
				if _, ok := comparator.(*Func); !ok {
					ctx.interp.throw("sort(func): expected func to be a function")
				}

				c := make([]Value, len(array))
				copy(c, array)
				sort.Slice(c, func(i, j int) bool {
					a := c[i]
					b := c[j]
					ret := comparator.(*Func).Executor(&FuncContext{
						interp: ctx.interp,
						Scope:  ctx.Scope,
						This:   wrap(nil),
						Args:   []Value{a, b},
					})

					if ret == nil {
						ctx.interp.throw("sort(func): expected function to return a value")
					}

					if _, ok := ret.Value.(*Integer); !ok {
						ctx.interp.throw("sort(func): expected function to return a number")
					}

					return ret.Value.(*Integer).Value < 0
				})

				return &Return{wrap(c)}
			}},
			"map": &Func{Executor: func(ctx *FuncContext) *Return {
				if len(ctx.Args) < 1 {
					ctx.interp.throw("map(func): expected at least 1 argument")
				}

				array := ctx.This.(*Array).Elements
				mapper := ctx.Args[0]
				if _, ok := mapper.(*Func); !ok {
					ctx.interp.throw("map(func): expected func to be a function")
				}

				c := make([]Value, len(array))
				for i, v := range array {
					ret := mapper.(*Func).Executor(&FuncContext{
						interp: ctx.interp,
						Scope:  ctx.Scope,
						This:   wrap(nil),
						Args:   []Value{v, wrap(i)},
					})

					c[i] = ret.Value
				}

				return &Return{wrap(c)}
			}},
			"join": &Func{Executor: func(ctx *FuncContext) *Return {
				if len(ctx.Args) < 1 {
					ctx.interp.throw("join(sep): expected at least 1 argument")
				}

				array := ctx.This.(*Array).Elements
				sep := ctx.Args[0].(*String).Value

				var b strings.Builder
				for i, v := range array {
					if i > 0 {
						b.WriteString(sep)
					}
					b.WriteString(toString(ctx.interp, ctx.Scope, v))
				}

				return &Return{wrap(b.String())}
			}},
		},
	},
	Operators: Operators{
		token.Eq: OpFunc(func(c *opctx[*Array, Value]) Value {
			switch o := c.other.(type) {
			case *Array:
				if len(c.this.Elements) != len(o.Elements) {
					return FalseValue
				}
				for i, e := range c.this.Elements {
					op := GetOperator(e, token.Eq)
					if op == nil {
						c.interp.throw("array element is not comparable")
					}
					ret := op.Executor(&FuncContext{
						interp: c.interp,
						Scope:  c.scope,
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
