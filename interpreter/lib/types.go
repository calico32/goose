package lib

import (
	"fmt"
	"math/big"
	"strconv"
	"unsafe"

	"github.com/calico32/goose/token"
)

type Value interface {
	gooseValue()
	Type() string
	Prototype() *Composite
	Clone() Value
	Unwrap() any
	Freeze()
	Unfreeze()
	Hash() string
}

type FuncContext struct {
	Interp Interpreter
	Scope  *Scope
	This   Value
	Args   []Value
}

type OpContext[T Value, U Value] struct {
	Interp Interpreter
	Scope  *Scope
	This   T
	Other  U
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
	BigInt() *big.Int
	Float64() float64
}

type PropertyKey interface {
	Value
	kind() PropertyKeyKind
	CanonicalValue() string
}

type PropertyKeyKind int

const (
	PKString PropertyKeyKind = iota
	PKInteger
	PKSymbol
	NumPKKinds
)

type (
	Null    struct{}
	Integer struct{ Value *big.Int }
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
		Frozen   bool
	}
	Composite struct {
		Name       string
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
		Frozen       bool
	}
	Generator struct {
		Async   bool
		Channel chan GeneratorMessage
		Frozen  bool
	}
	IntRange struct {
		Start *big.Int
		Stop  *big.Int
		Step  *big.Int
	}
	FloatRange struct {
		Start float64
		Stop  float64
		Step  float64
	}
)

func NewString(s string) *String      { return &String{Value: s} }
func NewInteger(i *big.Int) *Integer  { return &Integer{Value: i} }
func NewFloat(f float64) *Float       { return &Float{Value: f} }
func NewArray(values ...Value) *Array { return &Array{Elements: values} }

type ValueType interface {
	string | float64 | bool | []any | []Value |
		[]string | []int | []int64 | []float64 | []bool |
		Null | Bool | String | Func | Integer | Float | Array | Composite | Generator | IntRange | FloatRange |
		*Null | *Bool | *String | *Func | *Integer | *Float | *Array | *Composite | *Generator | *IntRange | *FloatRange | *Value
}

func (*Null) gooseValue()       {}
func (*Integer) gooseValue()    {}
func (*Float) gooseValue()      {}
func (*Symbol) gooseValue()     {}
func (*Bool) gooseValue()       {}
func (*String) gooseValue()     {}
func (*Array) gooseValue()      {}
func (*Composite) gooseValue()  {}
func (*Func) gooseValue()       {}
func (*Generator) gooseValue()  {}
func (*IntRange) gooseValue()   {}
func (*FloatRange) gooseValue() {}

func (*Null) Type() string       { return "Null" }
func (*Integer) Type() string    { return "Integer" }
func (*Float) Type() string      { return "Float" }
func (*Symbol) Type() string     { return "Symbol" }
func (*Bool) Type() string       { return "Bool" }
func (*String) Type() string     { return "String" }
func (*Array) Type() string      { return "Array" }
func (*Composite) Type() string  { return "Composite" }
func (*Func) Type() string       { return "Func" }
func (*Generator) Type() string  { return "Generator" }
func (*IntRange) Type() string   { return "IntRange" }
func (*FloatRange) Type() string { return "FloatRange" }

func (n *Null) Unwrap() any    { return nil }
func (i *Integer) Unwrap() any { return i.Value }
func (f *Float) Unwrap() any   { return f.Value }
func (s *Symbol) Unwrap() any  { return s.Name }
func (b *Bool) Unwrap() any    { return b.Value }
func (s *String) Unwrap() any  { return s.Value }
func (a *Array) Unwrap() any   { return a.Elements }
func (c *Composite) Unwrap() any {
	result := make(map[string]any)
	for kind := PropertyKeyKind(0); kind < NumPKKinds; kind++ {
		if c.Properties[kind] == nil {
			continue
		}

		for key, value := range c.Properties[kind] {
			result[key] = value.Unwrap()
		}
	}
	return result
}
func (f *Func) Unwrap() any       { return f.Executor }
func (g *Generator) Unwrap() any  { return g }
func (r *IntRange) Unwrap() any   { return r }
func (r *FloatRange) Unwrap() any { return r }

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
func (r *IntRange) Clone() Value {
	return &IntRange{
		Start: r.Start,
		Stop:  r.Stop,
		Step:  r.Step,
	}
}
func (r *FloatRange) Clone() Value {
	return &FloatRange{
		Start: r.Start,
		Stop:  r.Stop,
		Step:  r.Step,
	}
}

func (n *Null) Freeze()    {}
func (i *Integer) Freeze() {}
func (f *Float) Freeze()   {}
func (s *Symbol) Freeze()  {}
func (b *Bool) Freeze()    {}
func (s *String) Freeze()  {}
func (a *Array) Freeze() {
	a.Frozen = true
}
func (c *Composite) Freeze() {
	c.Frozen = true
}
func (f *Func) Freeze() {
	f.Frozen = true
}
func (g *Generator) Freeze() {
	g.Frozen = true
}
func (r *IntRange) Freeze()   {}
func (r *FloatRange) Freeze() {}

func (n *Null) Unfreeze()    {}
func (i *Integer) Unfreeze() {}
func (f *Float) Unfreeze()   {}
func (s *Symbol) Unfreeze()  {}
func (b *Bool) Unfreeze()    {}
func (s *String) Unfreeze()  {}
func (a *Array) Unfreeze() {
	a.Frozen = false
}
func (c *Composite) Unfreeze() {
	c.Frozen = false
}
func (f *Func) Unfreeze() {
	f.Frozen = false
}
func (g *Generator) Unfreeze() {
	g.Frozen = false
}
func (r *IntRange) Unfreeze()   {}
func (r *FloatRange) Unfreeze() {}

func (*Integer) numeric() {}
func (*Float) numeric()   {}

func (i *Integer) Int() int         { return int(i.Value.Int64()) }
func (i *Integer) Int64() int64     { return i.Value.Int64() }
func (i *Integer) BigInt() *big.Int { return i.Value }
func (i *Integer) Float64() float64 {
	f, _ := i.Value.Float64()
	return f
}
func (f *Float) Int() int         { return int(f.Value) }
func (f *Float) Int64() int64     { return int64(f.Value) }
func (f *Float) BigInt() *big.Int { return big.NewInt(int64(f.Value)) }
func (f *Float) Float64() float64 { return f.Value }

func (*String) kind() PropertyKeyKind  { return PKString }
func (*Integer) kind() PropertyKeyKind { return PKInteger }
func (*Symbol) kind() PropertyKeyKind  { return PKSymbol }

func (s *String) CanonicalValue() string  { return s.Value }
func (i *Integer) CanonicalValue() string { return i.Value.Text(10) }
func (s *Symbol) CanonicalValue() string  { return "@" + s.Name + "#" + strconv.FormatInt(s.Id, 10) }

func (n *Null) Prototype() *Composite       { return NullPrototype }
func (i *Integer) Prototype() *Composite    { return IntegerPrototype }
func (f *Float) Prototype() *Composite      { return FloatPrototype }
func (s *Symbol) Prototype() *Composite     { return SymbolPrototype }
func (b *Bool) Prototype() *Composite       { return BoolPrototype }
func (s *String) Prototype() *Composite     { return StringPrototype }
func (a *Array) Prototype() *Composite      { return ArrayPrototype }
func (c *Composite) Prototype() *Composite  { return c.Proto }
func (f *Func) Prototype() *Composite       { return FuncPrototype }
func (g *Generator) Prototype() *Composite  { return Object }
func (r *IntRange) Prototype() *Composite   { return RangePrototype }
func (r *FloatRange) Prototype() *Composite { return RangePrototype }

func (n *Null) Hash() string      { return "null" }
func (i *Integer) Hash() string   { return i.Value.Text(10) }
func (f *Float) Hash() string     { return strconv.FormatFloat(f.Value, 'f', -1, 64) }
func (s *Symbol) Hash() string    { return strconv.FormatInt(s.Id, 10) }
func (b *Bool) Hash() string      { return strconv.FormatBool(b.Value) }
func (s *String) Hash() string    { return s.Value }
func (a *Array) Hash() string     { return strconv.FormatUint(uint64(uintptr(unsafe.Pointer(a))), 10) }
func (c *Composite) Hash() string { return strconv.FormatUint(uint64(uintptr(unsafe.Pointer(c))), 10) }
func (f *Func) Hash() string      { return strconv.FormatUint(uint64(uintptr(unsafe.Pointer(f))), 10) }
func (g *Generator) Hash() string { return strconv.FormatUint(uint64(uintptr(unsafe.Pointer(g))), 10) }
func (r *IntRange) Hash() string {
	return fmt.Sprintf("%s:%s:%s", r.Start.Text(10), r.Stop.Text(10), r.Step.Text(10))
}
func (r *FloatRange) Hash() string { return fmt.Sprintf("%f:%f:%f", r.Start, r.Stop, r.Step) }

func (r *IntRange) Contains(i *Integer) bool {
	return i.Value.Cmp(r.Start) >= 0 && i.Value.Cmp(r.Stop) < 0
}
func (r *FloatRange) Contains(f *Float) bool {
	return f.Value >= r.Start && f.Value < r.Stop
}

func GetProperty(v Value, key PropertyKey) Value {
	array, ok1 := v.(*Array)
	index, ok2 := key.(*Integer)
	if ok1 && ok2 {
		if index.Value.Cmp(big.NewInt(0)) == -1 || index.Value.Cmp(big.NewInt(int64(len(array.Elements)))) >= 0 {
			return NullValue
		}

		return array.Elements[index.Value.Int64()]
	}

	var c *Composite
	if comp, ok := v.(*Composite); ok {
		c = comp
	} else {
		c = v.Prototype()
	}

	val := c.Properties[key.kind()][key.CanonicalValue()]
	if val != nil {
		return val
	}

	if c.Prototype() != nil {
		return GetProperty(c.Prototype(), key)
	}

	return NullValue
}

func SetProperty(v Value, key PropertyKey, val Value) error {
	if a, ok := v.(*Array); ok {
		index, ok := key.(*Integer)
		if !ok {
			return fmt.Errorf("cannot use %s as an index", key.Type())
		}

		if index.Value.Cmp(big.NewInt(0)) == -1 || index.Value.Cmp(big.NewInt(int64(len(a.Elements)))) >= 0 {
			return fmt.Errorf("index %s out of range", index.Value.Text(10))
		}

		a.Elements[index.Value.Int64()] = val
		return nil
	}
	var c *Composite
	if comp, ok := v.(*Composite); ok {
		c = comp
	} else {
		c = v.Prototype()
	}
	if c.Properties[key.kind()] == nil {
		c.Properties[key.kind()] = make(map[string]Value)
	}
	c.Properties[key.kind()][key.CanonicalValue()] = val
	return nil
}

func GetOperator(v Value, tok token.Token) *OperatorFunc {
	var derivedOperators = map[token.Token]FuncType{
		token.Lt: func(ctx *FuncContext) *Return {
			gt := GetOperator(ctx.Args[0], token.Gt).Executor(ctx)
			eq := GetOperator(ctx.Args[0], token.Eq).Executor(ctx)
			return NewReturn(!gt.Value.(*Bool).Value && !eq.Value.(*Bool).Value)
		},
		token.Lte: func(ctx *FuncContext) *Return {
			gt := GetOperator(ctx.Args[0], token.Gt).Executor(ctx)
			return NewReturn(!gt.Value.(*Bool).Value)
		},
		token.Gte: func(ctx *FuncContext) *Return {
			gt := GetOperator(ctx.Args[0], token.Gt).Executor(ctx)
			eq := GetOperator(ctx.Args[0], token.Eq).Executor(ctx)
			return NewReturn(gt.Value.(*Bool).Value || eq.Value.(*Bool).Value)
		},
		token.Neq: func(ctx *FuncContext) *Return {
			eq := GetOperator(ctx.Args[0], token.Eq).Executor(ctx)
			return NewReturn(!eq.Value.(*Bool).Value)
		},
		token.Assign: func(ctx *FuncContext) *Return {
			return &Return{ctx.Args[0]}
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

func IsTruthy(v Value) bool {
	switch v := v.(type) {
	case *Null:
		return false
	case *Bool:
		return v.Value
	case *Integer:
		return v.Value.Cmp(big.NewInt(0)) != 0
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

	for kind := PropertyKeyKind(0); kind < NumPKKinds; kind++ {
		c.Properties[kind] = make(map[string]Value)
	}

	return c
}
