package lib

import (
	"fmt"
	"strings"

	"github.com/calico32/goose/token"
)

var StringPrototype = &Composite{
	Proto:  Object,
	Frozen: true,
	Properties: Properties{
		PKString: {
			"toString": &Func{Executor: func(ctx *FuncContext) *Return {
				return NewReturn(&ctx.This)
			}},
			"split": &Func{Executor: func(ctx *FuncContext) *Return {
				if len(ctx.Args) == 0 {
					return NewReturn([]Value{ctx.This})
				}
				s := ToString(ctx.Interp, ctx.Scope, ctx.Args[0])
				return NewReturn(strings.Split(ctx.This.(*String).Value, s))
			}},
			"padLeft": &Func{Executor: func(ctx *FuncContext) *Return {
				if len(ctx.Args) == 0 {
					ctx.Interp.Throw("padLeft(x): expected 1 argument")
				}
				if _, ok := ctx.Args[0].(Numeric); !ok {
					ctx.Interp.Throw("padLeft(x): expected integer as first argument")
				}
				pad := " "
				if len(ctx.Args) == 2 {
					if _, ok := ctx.Args[1].(*String); !ok {
						ctx.Interp.Throw("padLeft(x, y): expected string as second argument")
					}
					pad = ctx.Args[1].(*String).Value
				}

				n := ctx.Args[0].(Numeric).Int()
				needed := n - len(ctx.This.(*String).Value)
				out := ctx.This.(*String).Value
				for needed > 0 {
					out = pad + out
					needed -= len(pad)
				}

				return NewReturn(NewString(out))
			}},
			"padRight": &Func{Executor: func(ctx *FuncContext) *Return {
				if len(ctx.Args) == 0 {
					ctx.Interp.Throw("padRight(x): expected 1 argument")
				}
				if _, ok := ctx.Args[0].(Numeric); !ok {
					ctx.Interp.Throw("padRight(x): expected integer as first argument")
				}
				n := ctx.Args[0].(Numeric).Int()
				return NewReturn(NewString(fmt.Sprintf("%-*s", n, ctx.This.(*String).Value)))
			}},
			"slice": &Func{Executor: func(ctx *FuncContext) *Return {
				if len(ctx.Args) != 1 && len(ctx.Args) != 2 {
					ctx.Interp.Throw("slice(x, y): expected 2 arguments")
				}
				if _, ok := ctx.Args[0].(Numeric); !ok {
					ctx.Interp.Throw("slice(x, y): expected integer as first argument")
				}
				if len(ctx.Args) == 2 && ctx.Args[1] != nil {
					if _, ok := ctx.Args[1].(Numeric); !ok {
						ctx.Interp.Throw("slice(x, y): expected integer as second argument")
					}
				}
				x := ctx.Args[0].(Numeric).Int()
				y := len(ctx.This.(*String).Value)
				if len(ctx.Args) == 2 && ctx.Args[1] != nil {
					y = ctx.Args[1].(Numeric).Int()
				}
				if x < 0 {
					x = len(ctx.This.(*String).Value) + x
				}
				if y < 0 {
					y = len(ctx.This.(*String).Value) + y
				}
				return NewReturn(NewString(ctx.This.(*String).Value[x:y]))
			}},
			"trim": &Func{Executor: func(ctx *FuncContext) *Return {
				return NewReturn(NewString(strings.TrimSpace(ctx.This.(*String).Value)))
			}},
			"endsWith": &Func{Executor: func(ctx *FuncContext) *Return {
				if len(ctx.Args) == 0 {
					ctx.Interp.Throw("endsWith(x): expected 1 argument")
				}
				if _, ok := ctx.Args[0].(*String); !ok {
					ctx.Interp.Throw("endsWith(x): expected string as first argument")
				}
				return NewReturn(&Bool{strings.HasSuffix(ctx.This.(*String).Value, ctx.Args[0].(*String).Value)})
			}},
			"startsWith": &Func{Executor: func(ctx *FuncContext) *Return {
				if len(ctx.Args) == 0 {
					ctx.Interp.Throw("startsWith(x): expected 1 argument")
				}
				if _, ok := ctx.Args[0].(*String); !ok {
					ctx.Interp.Throw("startsWith(x): expected string as first argument")
				}
				return NewReturn(&Bool{strings.HasPrefix(ctx.This.(*String).Value, ctx.Args[0].(*String).Value)})
			}},
			"toUpperCase": &Func{Executor: func(ctx *FuncContext) *Return {
				return NewReturn(NewString(strings.ToUpper(ctx.This.(*String).Value)))
			}},
			"toLowerCase": &Func{Executor: func(ctx *FuncContext) *Return {
				return NewReturn(NewString(strings.ToLower(ctx.This.(*String).Value)))
			}},
		},
	},
	Operators: Operators{
		token.Eq: OpFunc(func(c *OpContext[*String, Value]) Value {
			switch o := c.Other.(type) {
			case *String:
				return &Bool{c.This.Value == o.Value}
			default:
				return FalseValue
			}
		}),
		token.Add: OpFunc(func(c *OpContext[*String, Value]) Value {
			return NewString(c.This.Value + ToString(c.Interp, c.Scope, c.Other))
		}),
		token.Mul: OpFunc(func(c *OpContext[*String, Value]) Value {
			if n, ok := c.Other.(Numeric); ok {
				return NewString(strings.Repeat(c.This.Value, n.Int()))
			}
			c.Interp.Throw("cannot multiply string by non-integer")
			return nil
		}),
	},
}
