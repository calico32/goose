package interpreter

import (
	"strconv"

	"github.com/calico32/goose/token"
)

type Value interface {
	gooseValue()
	Type() string
	Prototype() *Composite
	Clone() Value
	Unwrap() any
}

type FuncContext struct {
	interp *interp
	Scope  *Scope
	This   Value
	Args   []Value
}

type opctx[T Value, U Value] struct {
	interp *interp
	scope  *Scope
	this   T
	other  U
}

type FuncType func(*FuncContext) *Return

type Properties map[PropertyKeyKind]map[string]Value
type Operators map[token.Token]*OperatorFunc

type OperatorFunc struct {
	Async    bool
	Builtin  bool
	Executor FuncType
}

type Numeric interface {
	Value
	numeric()
	Int() int
	Int64() int64
	Float64() float64
}

type PropertyKey interface {
	Value
	propertyKeyKind() PropertyKeyKind
	CanonicalValue() string
}

type PropertyKeyKind int

const (
	PropertyKeyString PropertyKeyKind = iota
	PropertyKeyInteger
	PropertyKeySymbol
	NumPropertyKeyKinds
)

type (
	Null    struct{}
	Integer struct{ Value int64 }
	Float   struct{ Value float64 }
	Symbol  struct {
		Name string
		Id   int64
	}
	Bool struct{ Value bool }

	String struct {
		Value string
	}
	Array struct {
		Elements []Value
	}
	Composite struct {
		Proto      *Composite
		Properties Properties
		Operators  Operators
		Frozen     bool
	}
	Func struct {
		Executor     FuncType
		Async        bool
		Memoized     bool
		This         Value
		NewableProto *Composite
	}
	Generator struct {
		Async   bool
		Channel chan GeneratorMessage
	}
)

func (*Null) gooseValue()      {}
func (*Integer) gooseValue()   {}
func (*Float) gooseValue()     {}
func (*Symbol) gooseValue()    {}
func (*Bool) gooseValue()      {}
func (*String) gooseValue()    {}
func (*Array) gooseValue()     {}
func (*Composite) gooseValue() {}
func (*Func) gooseValue()      {}
func (*Generator) gooseValue() {}

func (*Null) Type() string      { return "Null" }
func (*Integer) Type() string   { return "Integer" }
func (*Float) Type() string     { return "Float" }
func (*Symbol) Type() string    { return "Symbol" }
func (*Bool) Type() string      { return "Bool" }
func (*String) Type() string    { return "String" }
func (*Array) Type() string     { return "Array" }
func (*Composite) Type() string { return "Composite" }
func (*Func) Type() string      { return "Func" }
func (*Generator) Type() string { return "Generator" }

func (n *Null) Unwrap() any    { return nil }
func (i *Integer) Unwrap() any { return i.Value }
func (f *Float) Unwrap() any   { return f.Value }
func (s *Symbol) Unwrap() any  { return s.Name }
func (b *Bool) Unwrap() any    { return b.Value }
func (s *String) Unwrap() any  { return s.Value }
func (a *Array) Unwrap() any   { return a.Elements }
func (c *Composite) Unwrap() any {
	result := make(map[string]any)
	for kind := PropertyKeyKind(0); kind < NumPropertyKeyKinds; kind++ {
		if c.Properties[kind] == nil {
			continue
		}

		for key, value := range c.Properties[kind] {
			result[key] = value.Unwrap()
		}
	}
	return result
}
func (f *Func) Unwrap() any      { return f.Executor }
func (g *Generator) Unwrap() any { return g }

func (n *Null) Clone() Value    { return n }
func (i *Integer) Clone() Value { return i }
func (f *Float) Clone() Value   { return f }
func (s *Symbol) Clone() Value  { return s }
func (b *Bool) Clone() Value    { return b }
func (s *String) Clone() Value  { return s }
func (a *Array) Clone() Value   { return a }
func (c *Composite) Clone() Value {
	return &Composite{
		Proto:      c.Proto,
		Properties: c.Properties,
		Operators:  c.Operators,
		Frozen:     c.Frozen,
	}
}
func (f *Func) Clone() Value { return f }
func (g *Generator) Clone() Value {
	return &Generator{
		Async:   g.Async,
		Channel: g.Channel,
	}
}

func (*Integer) numeric() {}
func (*Float) numeric()   {}

func (i *Integer) Int() int         { return int(i.Value) }
func (i *Integer) Int64() int64     { return i.Value }
func (i *Integer) Float64() float64 { return float64(i.Value) }
func (f *Float) Int() int           { return int(f.Value) }
func (f *Float) Int64() int64       { return int64(f.Value) }
func (f *Float) Float64() float64   { return f.Value }

func (*String) propertyKeyKind() PropertyKeyKind  { return PropertyKeyString }
func (*Integer) propertyKeyKind() PropertyKeyKind { return PropertyKeyInteger }
func (*Symbol) propertyKeyKind() PropertyKeyKind  { return PropertyKeySymbol }

func (s *String) CanonicalValue() string  { return s.Value }
func (i *Integer) CanonicalValue() string { return "i#" + strconv.FormatInt(i.Value, 10) }
func (s *Symbol) CanonicalValue() string  { return "@" + s.Name + "#" + strconv.FormatInt(s.Id, 10) }

func (n *Null) Prototype() *Composite      { return NullPrototype }
func (i *Integer) Prototype() *Composite   { return IntegerPrototype }
func (f *Float) Prototype() *Composite     { return FloatPrototype }
func (s *Symbol) Prototype() *Composite    { return SymbolPrototype }
func (b *Bool) Prototype() *Composite      { return BoolPrototype }
func (s *String) Prototype() *Composite    { return StringPrototype }
func (a *Array) Prototype() *Composite     { return ArrayPrototype }
func (c *Composite) Prototype() *Composite { return c.Proto }
func (f *Func) Prototype() *Composite      { return FuncPrototype }
func (g *Generator) Prototype() *Composite { return Object }

func GetProperty(v Value, key PropertyKey) Value {
	array, ok1 := v.(*Array)
	index, ok2 := key.(*Integer)
	if ok1 && ok2 {
		if index.Value < 0 || index.Value >= int64(len(array.Elements)) {
			return NullValue
		}

		return array.Elements[index.Value]
	}

	var c *Composite
	if comp, ok := v.(*Composite); ok {
		c = comp
	} else {
		c = v.Prototype()
	}

	val := c.Properties[key.propertyKeyKind()][key.CanonicalValue()]
	if val != nil {
		return val
	}

	if c.Prototype() != nil {
		return GetProperty(c.Prototype(), key)
	}

	return NullValue
}

func SetProperty(v Value, key PropertyKey, val Value) error {
	var c *Composite
	if comp, ok := v.(*Composite); ok {
		c = comp
	} else {
		c = v.Prototype()
	}
	c.Properties[key.propertyKeyKind()][key.CanonicalValue()] = val
	return nil
}

func GetOperator(v Value, tok token.Token) *OperatorFunc {
	var derivedOperators = map[token.Token]FuncType{
		token.Lt: func(ctx *FuncContext) *Return {
			gt := GetOperator(ctx.Args[0], token.Gt).Executor(ctx)
			eq := GetOperator(ctx.Args[0], token.Eq).Executor(ctx)
			return &Return{Value: wrap(!gt.Value.(*Bool).Value && !eq.Value.(*Bool).Value)}
		},
		token.Lte: func(ctx *FuncContext) *Return {
			gt := GetOperator(ctx.Args[0], token.Gt).Executor(ctx)
			return &Return{Value: wrap(!gt.Value.(*Bool).Value)}
		},
		token.Gte: func(ctx *FuncContext) *Return {
			gt := GetOperator(ctx.Args[0], token.Gt).Executor(ctx)
			eq := GetOperator(ctx.Args[0], token.Eq).Executor(ctx)
			return &Return{Value: wrap(gt.Value.(*Bool).Value || eq.Value.(*Bool).Value)}
		},
		token.Neq: func(ctx *FuncContext) *Return {
			eq := GetOperator(ctx.Args[0], token.Eq).Executor(ctx)
			return &Return{Value: wrap(!eq.Value.(*Bool).Value)}
		},
	}

	var tokenAliases = map[token.Token]token.Token{
		token.AddAssign:     token.Add,
		token.Inc:           token.Add,
		token.SubAssign:     token.Sub,
		token.Dec:           token.Sub,
		token.MulAssign:     token.Mul,
		token.QuoAssign:     token.Quo,
		token.PowAssign:     token.Pow,
		token.RemAssign:     token.Rem,
		token.LogAndAssign:  token.LogAnd,
		token.LogOrAssign:   token.LogOr,
		token.LogNullAssign: token.LogNull,
		token.BitAndAssign:  token.BitAnd,
		token.BitOrAssign:   token.BitOr,
		token.BitXorAssign:  token.BitXor,
		token.BitShlAssign:  token.BitShl,
		token.BitShrAssign:  token.BitShr,
	}

	var op *OperatorFunc
	proto := v.Prototype()
	for proto != nil {
		if proto.Operators != nil {
			if alias, ok := tokenAliases[tok]; ok {
				op = proto.Operators[alias]
			} else {
				op = proto.Operators[tok]
			}
		}
		if op != nil {
			break
		}

		if proto == Object && (proto.Prototype() == Object || proto.Prototype() == nil) {
			break
		}

		proto = proto.Prototype()
	}

	if op == nil {
		if op, ok := derivedOperators[tok]; ok {
			return &OperatorFunc{Builtin: true, Executor: op}
		}
	}

	return op
}

func isTruthy(v Value) bool {
	switch v := v.(type) {
	case *Null:
		return false
	case *Bool:
		return v.Value
	case *Integer:
		return v.Value != 0
	case *Float:
		return v.Value != 0
	case *String:
		return v.Value != ""
	case *Array:
		return len(v.Elements) != 0
	case *Composite:
		return len(v.Properties) != 0
	case *Func:
		return true
	case *Generator:
		return true
	}
	panic("unreachable")
}

func NewComposite() *Composite {
	c := &Composite{
		Proto:      Object,
		Properties: make(Properties),
		Operators:  make(Operators),
		Frozen:     false,
	}

	for kind := PropertyKeyKind(0); kind < NumPropertyKeyKinds; kind++ {
		c.Properties[kind] = make(map[string]Value)
	}

	return c
}
