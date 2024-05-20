package lib

import (
	"fmt"
	"strings"

	"github.com/calico32/goose/token"
)

var BuiltinSingletons = map[*Composite]Value{}

var Object = &Composite{
	Name:   "Object",
	Frozen: true,
	Properties: Properties{
		PKString: {
			"toString": &Func{
				Executor: func(ctx *FuncContext) *Return {
					proto := ctx.This.Prototype()
					name := ctx.This.Type()
					if proto != nil && proto.Name != "" {
						name = proto.Name
					}
					return &Return{&String{"<object " + name + ">"}}
				},
			},
			"toDebugString": &Func{
				Executor: func(ctx *FuncContext) *Return {
					if _, ok := ctx.This.(*Composite); !ok {
						toString := GetProperty(ctx.This, NewString("toString")).(*Func)

						return toString.Executor(&FuncContext{
							Interp: ctx.Interp,
							Scope:  ctx.Scope,
							This:   ctx.This,
						})
					}
					depth := 0
					if len(ctx.Args) > 0 {
						if _, ok := ctx.Args[0].(*Integer); !ok {
							ctx.Interp.Throw("toDebugString(depth): expected integer as first argument")
						} else {
							depth = int(ctx.Args[0].(*Integer).Value.Int64())
						}
					}
					var sb strings.Builder
					proto := ctx.This.Prototype()
					name := ctx.This.Type()
					if proto != nil {
						name = proto.Name
					}
					indent := strings.Repeat("  ", depth)
					sb.WriteString(name + " {\n")
					for k, v := range ctx.This.(*Composite).Properties[PKString] {
						sb.WriteString(indent + "  " + k + ": " + ToDebugString(ctx.Interp, ctx.Scope, v, depth+1) + ",\n")
					}
					for k, v := range ctx.This.(*Composite).Properties[PKInteger] {
						sb.WriteString(indent + "  [" + k + "]: " + ToDebugString(ctx.Interp, ctx.Scope, v, depth+1) + ",\n")
					}
					for k, v := range ctx.This.(*Composite).Properties[PKSymbol] {
						sb.WriteString(indent + "  [" + k + "]: " + ToDebugString(ctx.Interp, ctx.Scope, v, depth+1) + ",\n")
					}
					sb.WriteString(indent + "}")
					return &Return{&String{sb.String()}}
				},
			},
		},
	},
	Operators: Operators{
		token.Assign: OpFunc(func(c *OpContext[Value, Value]) Value {
			return c.Other
		}),
		token.Eq: OpFunc(func(c *OpContext[Value, Value]) Value {
			return BoolFrom[c.This == c.Other] // object identity
		}),
		token.LogNot: OpFunc(func(c *OpContext[Value, Value]) Value {
			return BoolFrom[!IsTruthy(c.This)]
		}),
		token.Question: OpFunc(func(c *OpContext[Value, Value]) Value {
			prop := GetProperty(c.This, NewString("toString"))
			if prop == nil {
				return &String{"<unknown>"}
			}

			if _, ok := prop.(*Func); !ok {
				return &String{"<object Func>"}
			}

			ret := prop.(*Func).Executor(&FuncContext{
				Interp: c.Interp,
				Scope:  c.Scope,
				This:   c.This,
			})

			return ret.Value
		}),
		token.LogAnd: OpFunc(func(c *OpContext[Value, Value]) Value {
			if IsTruthy(c.This) {
				return c.Other
			} else { // short-circuit
				return c.This
			}
		}),
		token.LogOr: OpFunc(func(c *OpContext[Value, Value]) Value {
			if IsTruthy(c.This) { // short-circuit
				return c.This
			} else {
				return c.Other
			}
		}),
		token.LogNull: OpFunc(func(c *OpContext[Value, Value]) Value {
			_, ok := c.This.(*Null)
			if ok { // if we are null, return other
				return c.Other
			} else { // short-circuit
				return c.This
			}
		}),
		token.Is: OpFunc(func(c *OpContext[Value, Value]) Value {
			p1 := c.This.Prototype()
			p2 := c.Other

			// for primitives, check against their builtin singletons
			if singleton, ok := BuiltinSingletons[p1]; ok && p2 == singleton {
				return TrueValue
			}

			for p2 != nil {
				if p1 == p2 {
					return TrueValue
				}
				switch x := p2.(type) {
				case *Func:
					p2 = x.NewableProto
				case *Composite:
					if x == nil {
						return FalseValue
					}
					p2 = x.Prototype()
				default:
					p2 = p2.Prototype()
				}

			}

			return FalseValue
		}),
		token.IsNot: OpFunc(func(c *OpContext[Value, Value]) Value {
			p1 := c.This.Prototype()
			p2 := c.Other
			for p2 != nil {
				if p1 == p2 {
					return FalseValue
				}
				switch x := p2.(type) {
				case *Func:
					p2 = x.NewableProto
				case *Composite:
					if x == nil {
						return TrueValue
					}
					p2 = x.Prototype()
				default:
					p2 = p2.Prototype()
				}
			}

			return TrueValue
		}),
	},
}

func operatorNotDefined(this string, op token.Token, other string) string {
	return fmt.Sprintf("operator %s not defined for types %s and %s", op, this, other)
}

func OpFunc[T Value, U Value](op func(*OpContext[T, U]) Value) *OperatorFunc {
	return &OperatorFunc{
		Builtin: true,
		Executor: func(ctx *FuncContext) *Return {
			var ret Value
			if len(ctx.Args) == 0 {
				ret = op(&OpContext[T, U]{
					Interp: ctx.Interp,
					Scope:  ctx.Scope,
					This:   ctx.This.(T),
				})
			} else {
				val, ok := ctx.Args[0].(U)
				if !ok {
					ctx.Interp.Throw("operator %s not defined for types %s and %s", token.Add, ctx.This.Type(), ctx.Args[0].Type())
				}
				ret = op(&OpContext[T, U]{
					Interp: ctx.Interp,
					Scope:  ctx.Scope,
					This:   ctx.This.(T),
					Other:  val,
				})
			}

			return &Return{ret}
		},
	}
}
