package cstring

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func Constant(in string) *constant.CharArray {
	return constant.NewCharArray(append([]byte(in), 0))
}

func ToI8Ptr(bb *ir.Block, src value.Value) value.Value {
	return src
	// return bb.NewGetElementPtr(ElemType(src), src, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
}

func Len(block *ir.Block, src value.Value) value.Value {
	if _, ok := src.Type().(*types.PointerType); ok {
		l := block.NewGetElementPtr(ElemType(src), src, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
		return block.NewLoad(ElemType(l), l)
	}
	return block.NewExtractValue(src, 0)
}

func ElemType(src value.Value) types.Type {
	return src.Type().(*types.PointerType).ElemType
}

var globalStringCounter uint

func NextStringName() string {
	name := fmt.Sprintf("str.%d", globalStringCounter)
	globalStringCounter++
	return name
}
