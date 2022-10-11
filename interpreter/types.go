package interpreter

import (
	"fmt"
)

type GooseValue struct {
	Constant bool
	Type     GooseType
	Value    any
}

func wrap(value any) *GooseValue {
	return &GooseValue{Value: valueOf(value), Type: typeOf(value)}
}

var null = &GooseValue{
	Type:  GooseTypeNull,
	Value: nil,
}

type gooseExit struct{ code int }

func (e gooseExit) Error() string {
	return fmt.Sprintf("exit(%d)", e.code)
}

func (v *GooseValue) Copy() *GooseValue {
	return &GooseValue{
		Constant: v.Constant,
		Type:     v.Type,
		Value:    v.Value,
	}
}

type GooseType int

const (
	GooseTypeNull GooseType = 1 << iota
	GooseTypeInt
	GooseTypeFloat
	GooseTypeBool
	GooseTypeFunc
	GooseTypeString
	GooseTypeArray
	GooseTypeComposite
	gooseTypeCount = iota

	GooseTypeError            = -1
	GooseTypeNumeric          = GooseTypeInt | GooseTypeFloat
	GooseTypeMutableIndexable = GooseTypeArray | GooseTypeComposite
)

type GooseComposite map[any]*GooseValue
type GooseFunc func(*GooseScope, []*GooseValue) (*ReturnResult, error)

func (t GooseType) String() string {
	switch t {
	case GooseTypeNull:
		return "null"
	case GooseTypeInt:
		return "int"
	case GooseTypeFloat:
		return "float"
	case GooseTypeBool:
		return "bool"
	case GooseTypeFunc:
		return "function"
	case GooseTypeString:
		return "string"
	case GooseTypeArray:
		return "array"
	case GooseTypeComposite:
		return "composite"
	case GooseTypeError:
		return "<error_type>"
	case GooseTypeNumeric:
		return "<numeric>"
	default:
		return fmt.Sprintf("<unknown_type:%d>", t)
	}
}

type (
	StmtResult interface {
		isStmtResult()
	}

	ReturnResult   struct{ value any }
	BreakResult    struct{}
	ContinueResult struct{}
	VoidResult     struct{}
)

func (ReturnResult) isStmtResult()   {}
func (BreakResult) isStmtResult()    {}
func (ContinueResult) isStmtResult() {}
func (VoidResult) isStmtResult()     {}
