package interpreter

import (
	"fmt"

	"github.com/calico32/goose/token"
)

var Object = &Composite{
	Frozen: true,
	Properties: Properties{
		PropertyKeyString: {
			"toString": &Func{
				Executor: func(ctx *FuncContext) *Return {
					return &Return{&String{"<object " + ctx.This.Type() + ">"}}
				},
			},
		},
	},
	Operators: Operators{
		token.Assign: OpFunc(func(c *opctx[Value, Value]) Value {
			return c.other
		}),
		token.Eq: OpFunc(func(c *opctx[Value, Value]) Value {
			return BoolFrom[c.this == c.other] // object identity
		}),
		token.LogNot: OpFunc(func(c *opctx[Value, Value]) Value {
			return BoolFrom[!isTruthy(c.this)]
		}),
		token.Question: OpFunc(func(c *opctx[Value, Value]) Value {
			prop := GetProperty(c.this, &String{"toString"})
			if prop == nil {
				return wrap("<unknown>")
			}

			if _, ok := prop.(*Func); !ok {
				return wrap("<unknown>")
			}

			ret := prop.(*Func).Executor(&FuncContext{
				interp: c.interp,
				Scope:  c.scope,
				This:   c.this,
			})

			return ret.Value
		}),
		token.LogAnd: OpFunc(func(c *opctx[Value, Value]) Value {
			if isTruthy(c.this) {
				return c.other
			} else { // short-circuit
				return c.this
			}
		}),
		token.LogOr: OpFunc(func(c *opctx[Value, Value]) Value {
			if isTruthy(c.this) { // short-circuit
				return c.this
			} else {
				return c.other
			}
		}),
		token.LogNull: OpFunc(func(c *opctx[Value, Value]) Value {
			_, ok := c.this.(*Null)
			if ok { // if we are null, return other
				return c.other
			} else { // short-circuit
				return c.this
			}
		}),
	},
}

func padCommon(ctx *FuncContext) (str string, pad string) {
	if len(ctx.Args) < 2 {
		ctx.interp.throw("pad__(x, int, val): expected at least 2 arguments")
		return
	}

	if _, ok := ctx.Args[0].(*String); !ok {
		ctx.interp.throw("pad__(x, int, val): expected string as first argument")
	}

	if _, ok := ctx.Args[1].(Numeric); !ok {
		ctx.interp.throw("pad__(x, int, val): expected integer as second argument")
	}

	length := ctx.Args[1].(Numeric).Int()

	if length < 0 {
		ctx.interp.throw("pad__(x, int, val): expected length >= 0")
		return
	}

	padChar := " "

	if len(ctx.Args) > 2 {
		if _, ok := ctx.Args[2].(*String); !ok {
			ctx.interp.throw("pad__(x, int, val): expected string as third argument")
		}

		padChar = ctx.Args[2].(*String).Value
	}

	if len(str) >= length {
		return
	}

	pad = ""
	for len(pad) < length-len(str) {
		pad += padChar
	}

	return str, pad
}

func operatorNotDefined(this string, op token.Token, other string) string {
	return fmt.Sprintf("operator %s not defined for types %s and %s", op, this, other)
}

func OpFunc[T Value, U Value](op func(*opctx[T, U]) Value) *OperatorFunc {
	return &OperatorFunc{
		Builtin: true,
		Executor: func(ctx *FuncContext) *Return {
			var other *U
			if len(ctx.Args) == 0 {
				other = nil
			} else {
				val, ok := ctx.Args[0].(U)
				if !ok {
					ctx.interp.throw("operator %s not defined for types %s and %s", token.Add, ctx.This.Type(), ctx.Args[0].Type())
				}
				other = &val
			}

			ret := op(&opctx[T, U]{
				interp: ctx.interp,
				scope:  ctx.Scope,
				this:   ctx.This.(T),
				other:  *other,
			})
			return &Return{ret}
		},
	}
}
