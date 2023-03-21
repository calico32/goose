package interpreter

import (
	"strings"

	"github.com/calico32/goose/token"
)

var StringPrototype = &Composite{
	Proto:  Object,
	Frozen: true,
	Properties: Properties{
		PropertyKeyString: {
			"toString": &Func{Executor: func(ctx *FuncContext) *Return {
				return &Return{ctx.This}
			}},
			"split": &Func{Executor: func(ctx *FuncContext) *Return {
				if len(ctx.Args) == 0 {
					return &Return{&Array{[]Value{ctx.This}}}
				}
				s := toString(ctx.interp, ctx.Scope, ctx.Args[0])
				return &Return{wrap(strings.Split(ctx.This.(*String).Value, s))}
			}},
		},
	},
	Operators: Operators{
		token.Eq: OpFunc(func(c *opctx[*String, Value]) Value {
			switch o := c.other.(type) {
			case *String:
				return &Bool{c.this.Value == o.Value}
			default:
				return FalseValue
			}
		}),
		token.Add: OpFunc(func(c *opctx[*String, Value]) Value {
			return &String{c.this.Value + toString(c.interp, c.scope, c.other)}
		}),
	},
}
