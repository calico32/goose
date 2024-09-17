package compiler

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
)

type IRFuncSignature struct {
	name     string
	ret      types.Type
	args     []*ir.Param
	variadic bool
}

var IRFuncs = map[string]IRFuncSignature{
	"exit": {
		name: "exit",
		ret:  types.Void,
		args: []*ir.Param{ir.NewParam("code", types.I32)},
	},
	"malloc": {
		name: "malloc",
		ret:  types.I8Ptr,
		args: []*ir.Param{ir.NewParam("size", types.I64)},
	},
	"printf": {
		name:     "printf",
		ret:      types.I32,
		args:     []*ir.Param{ir.NewParam("format", types.I8Ptr)},
		variadic: true,
	},
}

type ExternalFunctionRegistry struct {
	c   *Compiler
	fns map[string]*ir.Func
}

func NewExternalFunctionRegistry(c *Compiler) *ExternalFunctionRegistry {
	return &ExternalFunctionRegistry{c: c, fns: make(map[string]*ir.Func)}
}

func (f *ExternalFunctionRegistry) Get(name string) *ir.Func {
	if fn, ok := f.fns[name]; ok {
		return fn
	}
	if irfn, ok := IRFuncs[name]; ok {
		fn := f.c.m.NewFunc(name, irfn.ret, irfn.args...)
		fn.Sig.Variadic = irfn.variadic
		f.fns[name] = fn
		return fn
	}
	return nil
}
