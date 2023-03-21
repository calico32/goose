package interpreter

import (
	"math"
	"strconv"

	"github.com/calico32/goose/token"
)

var IntegerPrototype = &Composite{
	Proto:  Object,
	Frozen: true,
	Properties: Properties{
		PropertyKeyString: {
			"toString": &Func{
				Executor: func(ctx *FuncContext) *Return {
					base := 10
					if len(ctx.Args) >= 1 {
						base = int(ctx.Args[0].(*Integer).Value)
						if base < 2 || base > 36 {
							ctx.interp.throw("base must be between 2 and 36")
						}
					}
					return &Return{&String{strconv.FormatInt(ctx.This.(*Integer).Value, base)}}
				},
			},
		},
	},
	Operators: Operators{
		token.Eq: OpFunc(func(c *opctx[*Integer, Value]) Value {
			if i, ok := c.other.(*Integer); ok {
				return BoolFrom[c.this.Value == i.Value]
			}
			return FalseValue
		}),
		token.LogNot: OpFunc(func(c *opctx[*Integer, Value]) Value {
			return BoolFrom[c.this.Value == 0]
		}),
		token.Add: OpFunc(func(c *opctx[*Integer, Value]) Value {
			if c.other == nil {
				return c.this
			}
			switch o := c.other.(type) {
			case *Integer:
				return &Integer{c.this.Value + o.Value}
			case *Float:
				return &Float{float64(c.this.Value) + o.Value}
			case *String:
				return &String{strconv.FormatInt(c.this.Value, 10) + o.Value}
			default:
				c.interp.throw(operatorNotDefined(c.this.Type(), token.Add, o.Type()))
				return nil
			}
		}),
		token.Gt: OpFunc(func(c *opctx[*Integer, Numeric]) Value {
			return BoolFrom[float64(c.this.Value) > c.other.Float64()]
		}),
		token.Sub: OpFunc(func(c *opctx[*Integer, Numeric]) Value {
			if c.other == nil {
				return &Integer{-c.this.Value}
			}
			switch o := c.other.(type) {
			case *Integer:
				return &Integer{c.this.Value - o.Value}
			case *Float:
				return &Float{float64(c.this.Value) - o.Value}
			default:
				c.interp.throw(operatorNotDefined(c.this.Type(), token.Sub, o.Type()))
				return nil
			}
		}),
		token.Mul: OpFunc(func(c *opctx[*Integer, Numeric]) Value {
			switch o := c.other.(type) {
			case *Integer:
				return &Integer{c.this.Value * o.Value}
			case *Float:
				return &Float{float64(c.this.Value) * o.Value}
			default:
				c.interp.throw(operatorNotDefined(c.this.Type(), token.Mul, o.Type()))
				return nil
			}
		}),
		token.Quo: OpFunc(func(c *opctx[*Integer, Numeric]) Value {
			switch o := c.other.(type) {
			case *Integer:
				return &Integer{c.this.Value / o.Value}
			case *Float:
				return &Float{float64(c.this.Value) / o.Value}
			default:
				c.interp.throw(operatorNotDefined(c.this.Type(), token.Quo, o.Type()))
				return nil
			}
		}),
		token.BitAnd: OpFunc(func(c *opctx[*Integer, *Integer]) Value {
			return &Integer{c.this.Value & c.other.Value}
		}),
		token.BitOr: OpFunc(func(c *opctx[*Integer, *Integer]) Value {
			return &Integer{c.this.Value | c.other.Value}
		}),
		token.BitXor: OpFunc(func(c *opctx[*Integer, *Integer]) Value {
			return &Integer{c.this.Value ^ c.other.Value}
		}),
		token.BitNot: OpFunc(func(c *opctx[*Integer, *Integer]) Value {
			return &Integer{^c.this.Value}
		}),
		token.Rem: OpFunc(func(c *opctx[*Integer, Numeric]) Value {
			switch o := c.other.(type) {
			case *Integer:
				return &Integer{c.this.Value % o.Value}
			case *Float:
				return &Float{math.Mod(float64(c.this.Value), o.Value)}
			default:
				c.interp.throw(operatorNotDefined(c.this.Type(), token.Rem, o.Type()))
				return nil
			}
		}),
		token.Pow: OpFunc(func(c *opctx[*Integer, Numeric]) Value {
			return &Float{math.Pow(float64(c.this.Value), c.other.Float64())}
		}),
	},
}
