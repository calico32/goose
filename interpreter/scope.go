package interpreter

import "fmt"

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

var scopeOwnerNames = [...]string{
	ScopeOwnerFunc:      "function",
	ScopeOwnerArrayInit: "array initializer",
	ScopeOwnerRepeat:    "repeat",
	ScopeOwnerFor:       "for",
	ScopeOwnerIf:        "if",
	ScopeOwnerGlobal:    "global",
	ScopeOwnerBuiltin:   "builtin",
}

func (s ScopeOwner) String() string {
	return scopeOwnerNames[s]
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

func (s *GooseScope) fork(owner ScopeOwner) *GooseScope {
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
		panic(fmt.Errorf("cannot redefine builtin %s", name))
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
