package interpreter

import (
	"math"
	"strconv"

	"github.com/calico32/goose/token"
)

var FloatPrototype = &Composite{
	Proto:  Object,
	Frozen: true,
	Properties: Properties{
		PropertyKeyString: {
			"toString": &Func{
				Executor: func(ctx *FuncContext) *Return {
					return &Return{&String{strconv.FormatFloat(ctx.This.(*Float).Value, 'f', -1, 64)}}
				},
			},
		},
	},
	Operators: Operators{
		token.Eq: OpFunc(func(c *opctx[*Float, Value]) Value {
			if _, ok := c.other.(*Float); ok {
				return BoolFrom[c.this.Value == c.other.(*Float).Value]
			}
			return FalseValue
		}),
		token.Add: OpFunc(func(c *opctx[*Float, Value]) Value {
			if c.other == nil {
				return c.this
			}
			switch o := c.other.(type) {
			case Numeric:
				return &Float{c.this.Value + o.Float64()}
			case *String:
				return &String{strconv.FormatFloat(c.this.Value, 'f', -1, 64) + o.Value}
			default:
				c.interp.throw(operatorNotDefined("Float", token.Add, o.Type()))
				return nil
			}
		}),
		token.Gt: OpFunc(func(c *opctx[*Float, Numeric]) Value {
			return BoolFrom[c.this.Value > c.other.Float64()]
		}),
		token.Sub: OpFunc(func(c *opctx[*Float, Numeric]) Value {
			if c.other == nil {
				return &Float{-c.this.Value}
			}

			return &Float{c.this.Value - c.other.Float64()}
		}),
		token.Mul: OpFunc(func(c *opctx[*Float, Numeric]) Value {
			return &Float{c.this.Value * c.other.Float64()}
		}),
		token.Quo: OpFunc(func(c *opctx[*Float, Numeric]) Value {
			return &Float{c.this.Value / c.other.Float64()}
		}),
		token.Rem: OpFunc(func(c *opctx[*Float, Numeric]) Value {
			return &Float{math.Mod(c.this.Value, c.other.Float64())}
		}),
		token.Pow: OpFunc(func(c *opctx[*Float, Numeric]) Value {
			return &Float{math.Pow(c.this.Value, c.other.Float64())}
		}),
	},
}
