package lib

import (
	"math"
	"math/big"

	"github.com/calico32/goose/token"
)

var IntegerPrototype = &Composite{
	Proto:  Object,
	Frozen: true,
	Properties: Properties{
		PKString: {
			"toString": &Func{
				Executor: func(ctx *FuncContext) *Return {
					base := 10
					if len(ctx.Args) >= 1 {
						base = int(ctx.Args[0].(*Integer).Value.Int64())
						if base < 2 || base > 62 {
							ctx.Interp.Throw("base must be between 2 and 62")
						}
					}
					return NewReturn(NewString(ctx.This.(*Integer).Value.Text(base)))
				},
			},
		},
	},
	Operators: Operators{
		token.Eq: OpFunc(func(c *OpContext[*Integer, Value]) Value {
			if i, ok := c.Other.(*Integer); ok {
				return BoolFrom[c.This.Value.Cmp(i.Value) == 0]
			}
			return FalseValue
		}),
		token.LogNot: OpFunc(func(c *OpContext[*Integer, Value]) Value {
			return BoolFrom[c.This.Value.Cmp(big.NewInt(0)) == 0]
		}),
		token.Add: OpFunc(func(c *OpContext[*Integer, Value]) Value {
			if c.Other == nil {
				return c.This
			}
			switch o := c.Other.(type) {
			case *Integer:
				res := new(big.Int)
				res.Add(c.This.Value, o.Value)
				return &Integer{res}
			case *Float:
				f, _ := c.This.Value.Float64()
				return &Float{f + o.Value}
			case *String:
				return NewString(c.This.Value.Text(10) + o.Value)
			default:
				c.Interp.Throw(operatorNotDefined(c.This.Type(), token.Add, o.Type()))
				return nil
			}
		}),
		token.Gt: OpFunc(func(c *OpContext[*Integer, Numeric]) Value {
			switch o := c.Other.(type) {
			case *Integer:
				return BoolFrom[c.This.Value.Cmp(o.Value) > 0]
			case *Float:
				f, _ := c.This.Value.Float64()
				return BoolFrom[f > o.Value]
			default:
				c.Interp.Throw(operatorNotDefined(c.This.Type(), token.Gt, o.Type()))
				return nil
			}
		}),
		token.Sub: OpFunc(func(c *OpContext[*Integer, Numeric]) Value {
			if c.Other == nil {
				n := new(big.Int)
				n.Neg(c.This.Value)
				return &Integer{n}
			}
			switch o := c.Other.(type) {
			case *Integer:
				res := new(big.Int)
				res.Sub(c.This.Value, o.Value)
				return &Integer{res}
			case *Float:
				f, _ := c.This.Value.Float64()
				return &Float{f - o.Value}
			default:
				c.Interp.Throw(operatorNotDefined(c.This.Type(), token.Sub, o.Type()))
				return nil
			}
		}),
		token.Mul: OpFunc(func(c *OpContext[*Integer, Numeric]) Value {
			switch o := c.Other.(type) {
			case *Integer:
				res := new(big.Int)
				res.Mul(c.This.Value, o.Value)
				return &Integer{res}
			case *Float:
				f, _ := c.This.Value.Float64()
				return &Float{f * o.Value}
			default:
				c.Interp.Throw(operatorNotDefined(c.This.Type(), token.Mul, o.Type()))
				return nil
			}
		}),
		token.Quo: OpFunc(func(c *OpContext[*Integer, Numeric]) Value {
			switch o := c.Other.(type) {
			case *Integer:
				if o.Value.Cmp(big.NewInt(0)) == 0 {
					c.Interp.Throw("division by zero")
				}
				res := new(big.Int)
				res.Quo(c.This.Value, o.Value)
				return &Integer{res}
			case *Float:
				if o.Value == 0 {
					c.Interp.Throw("division by zero")
				}
				f, _ := c.This.Value.Float64()
				return &Float{f / o.Value}
			default:
				c.Interp.Throw(operatorNotDefined(c.This.Type(), token.Quo, o.Type()))
				return nil
			}
		}),
		token.BitAnd: OpFunc(func(c *OpContext[*Integer, *Integer]) Value {
			res := new(big.Int)
			res.And(c.This.Value, c.Other.Value)
			return &Integer{res}
		}),
		token.BitOr: OpFunc(func(c *OpContext[*Integer, *Integer]) Value {
			res := new(big.Int)
			res.Or(c.This.Value, c.Other.Value)
			return &Integer{res}
		}),
		token.BitXor: OpFunc(func(c *OpContext[*Integer, *Integer]) Value {
			res := new(big.Int)
			res.Xor(c.This.Value, c.Other.Value)
			return &Integer{res}
		}),
		token.BitNot: OpFunc(func(c *OpContext[*Integer, *Integer]) Value {
			res := new(big.Int)
			res.Not(c.This.Value)
			return &Integer{res}
		}),
		token.Rem: OpFunc(func(c *OpContext[*Integer, Numeric]) Value {
			switch o := c.Other.(type) {
			case *Integer:
				if o.Value.Cmp(big.NewInt(0)) == 0 {
					c.Interp.Throw("division by zero")
				}
				res := new(big.Int)
				res.Rem(c.This.Value, o.Value)
				return &Integer{res}
			case *Float:
				f, _ := c.This.Value.Float64()
				return &Float{math.Mod(f, o.Value)}
			default:
				c.Interp.Throw(operatorNotDefined(c.This.Type(), token.Rem, o.Type()))
				return nil
			}
		}),
		token.Pow: OpFunc(func(c *OpContext[*Integer, Numeric]) Value {
			f, _ := c.This.Value.Float64()
			return &Float{math.Pow(f, c.Other.Float64())}
		}),
	},
}
