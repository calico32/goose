package lib

var BoolPrototype = &Composite{
	Proto:  Object,
	Frozen: true,
	Properties: Properties{
		PKString: {
			"toString": &Func{
				Executor: func(ctx *FuncContext) *Return {
					if ctx.This.(*Bool).Value {
						return NewReturn(NewString("true"))
					} else {
						return NewReturn(NewString("false"))
					}
				},
			},
		},
	},
}
