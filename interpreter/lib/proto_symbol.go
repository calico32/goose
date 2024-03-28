package lib

var SymbolPrototype = &Composite{
	Proto:  Object,
	Frozen: true,
	Properties: Properties{
		PKString: {
			"toString": &Func{
				Executor: func(ctx *FuncContext) *Return {
					return NewReturn(NewString(ctx.This.(*Symbol).Name))
				},
			},
		},
	},
}
