package interpreter

import (
	"fmt"
)

type GooseValue struct {
	Constant bool
	Type     GooseType
	Value    any
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
	gooseTypeCount = iota

	GooseTypeError   = -1
	GooseTypeNumeric = GooseTypeInt | GooseTypeFloat
)

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

type GooseFunc func(*GooseScope, []*GooseValue) (*ReturnResult, error)

func (ReturnResult) isStmtResult()   {}
func (BreakResult) isStmtResult()    {}
func (ContinueResult) isStmtResult() {}
func (VoidResult) isStmtResult()     {}

type ScopeOwner int

const (
	ScopeOwnerFunc ScopeOwner = iota
	ScopeOwnerArrayInit
	ScopeOwnerRepeat
	ScopeOwnerFor
	ScopeOwnerIf
	ScopeOwnerGlobal
	ScopeOwnerBuiltin
)

func (s ScopeOwner) String() string {
	switch s {
	case ScopeOwnerFunc:
		return "ScopeOwnerFunc"
	case ScopeOwnerArrayInit:
		return "ScopeOwnerArrayInit"
	case ScopeOwnerRepeat:
		return "ScopeOwnerRepeat"
	case ScopeOwnerFor:
		return "ScopeOwnerFor"
	case ScopeOwnerIf:
		return "ScopeOwnerIf"
	case ScopeOwnerGlobal:
		return "ScopeOwnerGlobal"
	case ScopeOwnerBuiltin:
		return "ScopeOwnerBuiltin"
	default:
		return fmt.Sprintf("<unknown_scope_owner:%d>", s)
	}
}

type GooseScope struct {
	interp *interpreter
	owner  ScopeOwner
	parent *GooseScope
	idents map[string]GooseValue
}

func NewGlobalScope(i *interpreter) *GooseScope {
	return &GooseScope{
		interp: i,
		owner:  ScopeOwnerGlobal,
		parent: &GooseScope{
			interp: i,
			owner:  ScopeOwnerBuiltin,
			parent: nil,
			idents: builtins,
		},
	}
}

func (s *GooseScope) builtins() *GooseScope {
	if s.owner == ScopeOwnerBuiltin {
		return s
	}
	return s.parent.builtins()
}

func (s *GooseScope) global() *GooseScope {
	if s.owner == ScopeOwnerGlobal {
		return s
	}
	return s.parent.global()
}

func (s *GooseScope) new(owner ScopeOwner) *GooseScope {
	return &GooseScope{
		owner:  owner,
		parent: s,
		interp: s.interp,
	}
}

func (s *GooseScope) get(name string) *GooseValue {
	if v, ok := s.idents[name]; ok {
		return &v
	}
	if s.parent != nil {
		return s.parent.get(name)
	}
	return nil
}

func (s *GooseScope) set(name string, value GooseValue) {
	if s.idents == nil {
		s.idents = make(map[string]GooseValue)
	}

	if s.builtins().isDefined(name) {
		panic(fmt.Errorf("Cannot redefine builtin %s", name))
	}

	if v, ok := s.idents[name]; ok {
		if v.Constant {
			panic(fmt.Errorf("cannot assign to constant %s", name))
		}

	}

	s.idents[name] = value
}

func (s *GooseScope) update(name string, value GooseValue) {
	if s.idents == nil {
		s.idents = make(map[string]GooseValue)
	}

	// try to look up the name in the current scope, if it's not there, look up in the parent scope

	if v, ok := s.idents[name]; ok {
		if v.Constant {
			panic(fmt.Errorf("cannot assign to constant %s", name))
		}
		s.idents[name] = value
		return
	}

	if s.parent != nil {
		s.parent.update(name, value)
	} else {
		panic(fmt.Errorf("%s is not defined", name))
	}
}

func (s *GooseScope) isDefined(name string) bool {
	if _, ok := s.idents[name]; ok {
		return true
	}
	if s.parent != nil {
		return s.parent.isDefined(name)
	}
	return false
}

func (s *GooseScope) isDefinedInCurrentScope(name string) bool {
	if _, ok := s.idents[name]; ok {
		return true
	}
	return false
}
