package compiler

import (
	"github.com/calico32/goose/ast"
	"github.com/calico32/goose/compiler/cstring"
	. "github.com/calico32/goose/interpreter/lib"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type CModule struct {
	*ast.Module
	Scope   *CScope
	Exports map[string]*CVariable
}

type CVariable struct {
	Block    *ir.Block
	Constant bool
	Value    value.Value
	Source   VariableSource
}

var CGlobalConstants = map[string]*CVariable{
	"true":  {Constant: true, Value: constant.NewBool(true)},
	"false": {Constant: true, Value: constant.NewBool(false)},
	"null":  {Constant: true, Value: constant.NewNull(types.I32Ptr)},
}

type CFuncContext struct {
	C     *Compiler
	Scope *CScope
	This  value.Value
	Args  []value.Value
}

type CFuncType func(ctx *CFuncContext) value.Value
type FuncBuilder func(c *Compiler) *ir.Func

const GoosePrefix = "goose__"

var CGlobals = map[string]FuncBuilder{
	"print": func(c *Compiler) *ir.Func {
		format := c.m.NewGlobalDef(cstring.NextStringName(), cstring.Constant("%s"))
		format.Immutable = true

		fn := c.m.NewFunc(GoosePrefix+"print", types.Void, ir.NewParam("s", types.I8Ptr))
		b := fn.NewBlock("entry")

		printf := c.ext.Get("printf")
		b.NewCall(printf, format, fn.Params[0])
		b.NewRet(nil)

		return fn
	},
	"println": func(c *Compiler) *ir.Func {
		format := c.m.NewGlobalDef(cstring.NextStringName(), cstring.Constant("%s\n"))
		format.Immutable = true

		fn := c.m.NewFunc(GoosePrefix+"println", types.Void, ir.NewParam("s", types.I8Ptr))
		b := fn.NewBlock("entry")

		printf := c.ext.Get("printf")
		b.NewCall(printf, format, fn.Params[0])
		b.NewRet(nil)

		return fn
	},
	"printf": func(c *Compiler) *ir.Func {
		fn := c.m.NewFunc(GoosePrefix+"printf", types.I32, ir.NewParam("format", types.I8Ptr))
		fn.Sig.Variadic = true
		b := fn.NewBlock("entry")

		printf := c.ext.Get("printf")

		values := make([]value.Value, 0, len(fn.Params))
		for _, param := range fn.Params {
			values = append(values, param)
		}

		r := b.NewCall(printf, values...)
		b.NewRet(r)

		return fn
	},
}
