package lib

import (
	"math"
	"strconv"

	"github.com/calico32/goose/token"
)

var FloatPrototype = &Composite{
	Proto:  Object,
	Frozen: true,
	Properties: Properties{
		PKString: {
			"toString": &Func{
				Executor: func(ctx *FuncContext) *Return {
					return NewReturn(NewString(strconv.FormatFloat(ctx.This.(*Float).Value, 'f', -1, 64)))
				},
			},
		},
	},
	Operators: Operators{
		token.Eq: OpFunc(func(c *OpContext[*Float, Value]) Value {
			if _, ok := c.Other.(*Float); ok {
				return BoolFrom[c.This.Value == c.Other.(*Float).Value]
			}
			return FalseValue
		}),
		token.Add: OpFunc(func(c *OpContext[*Float, Value]) Value {
			if c.Other == nil {
				return c.This
			}
			switch o := c.Other.(type) {
			case Numeric:
				return &Float{c.This.Value + o.Float64()}
			case *String:
				return NewString(strconv.FormatFloat(c.This.Value, 'f', -1, 64) + o.Value)
			default:
				c.Interp.Throw(operatorNotDefined("Float", token.Add, o.Type()))
				return nil
			}
		}),
		token.Gt: OpFunc(func(c *OpContext[*Float, Numeric]) Value {
			return BoolFrom[c.This.Value > c.Other.Float64()]
		}),
		token.Sub: OpFunc(func(c *OpContext[*Float, Numeric]) Value {
			if c.Other == nil {
				return &Float{-c.This.Value}
			}

			return &Float{c.This.Value - c.Other.Float64()}
		}),
		token.Mul: OpFunc(func(c *OpContext[*Float, Numeric]) Value {
			return &Float{c.This.Value * c.Other.Float64()}
		}),
		token.Quo: OpFunc(func(c *OpContext[*Float, Numeric]) Value {
			return &Float{c.This.Value / c.Other.Float64()}
		}),
		token.Rem: OpFunc(func(c *OpContext[*Float, Numeric]) Value {
			return &Float{math.Mod(c.This.Value, c.Other.Float64())}
		}),
		token.Pow: OpFunc(func(c *OpContext[*Float, Numeric]) Value {
			return &Float{math.Pow(c.This.Value, c.Other.Float64())}
		}),
	},
}
