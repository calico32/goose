package lib

import (
	"github.com/calico32/goose/token"
)

var NullPrototype = &Composite{
	Proto:  Object,
	Frozen: true,
	Properties: Properties{
		PKString: {
			"toString": &Func{
				Executor: func(ctx *FuncContext) *Return {
					return NewReturn(NewString("null"))
				},
			},
		},
	},
	Operators: Operators{
		token.Eq: OpFunc(func(c *OpContext[Value, Value]) Value {
			_, ok := c.Other.(*Null)
			return BoolFrom[ok]
		}),
		token.LogNot: OpFunc(func(c *OpContext[Value, Value]) Value {
			return TrueValue
		}),
	},
}
