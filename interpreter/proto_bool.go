package interpreter

var BoolPrototype = &Composite{
	Proto:  Object,
	Frozen: true,
	Properties: Properties{
		PropertyKeyString: {
			"toString": &Func{
				Executor: func(ctx *FuncContext) *Return {
					if ctx.This.(*Bool).Value {
						return &Return{&String{"true"}}
					} else {
						return &Return{&String{"false"}}
					}
				},
			},
		},
	},
}
