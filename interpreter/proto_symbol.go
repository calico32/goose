package interpreter

var SymbolPrototype = &Composite{
	Proto:  Object,
	Frozen: true,
	Properties: Properties{
		PropertyKeyString: {
			"toString": &Func{
				Executor: func(ctx *FuncContext) *Return {
					return &Return{&String{"@" + ctx.This.(*Symbol).Name}}
				},
			},
		},
	},
}
