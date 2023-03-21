package interpreter

import "github.com/calico32/goose/token"

var NullPrototype = &Composite{
	Proto:  Object,
	Frozen: true,
	Properties: Properties{
		PropertyKeyString: {
			"toString": &Func{
				Executor: func(ctx *FuncContext) *Return {
					return &Return{&String{"null"}}
				},
			},
		},
	},
	Operators: Operators{
		token.Eq: OpFunc(func(c *opctx[Value, Value]) Value {
			_, ok := c.other.(*Null)
			return BoolFrom[ok]
		}),
		token.LogNot: OpFunc(func(c *opctx[Value, Value]) Value {
			return TrueValue
		}),
	},
}
