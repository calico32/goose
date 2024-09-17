package compiler

import (
	"fmt"

	. "github.com/calico32/goose/interpreter/lib"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
)

type CScope struct {
	fn     *ir.Func
	module *CModule
	owner  ScopeOwner
	parent *CScope
	idents map[string]*CVariable
}

func NewGlobalCScope(builtins map[string]*CVariable) *CScope {
	return &CScope{
		owner:  ScopeOwnerGlobal,
		idents: make(map[string]*CVariable),
		parent: &CScope{
			owner:  ScopeOwnerBuiltin,
			parent: nil,
			idents: builtins,
		},
	}
}

func (s *CScope) Builtins() *CScope {
	if s.owner == ScopeOwnerBuiltin {
		return s
	}
	if s.parent == nil {
		panic("no builtin scope found")
	}
	return s.parent.Builtins()
}

func (s *CScope) Global() *CScope {
	if s.owner == ScopeOwnerGlobal {
		return s
	}
	if s.parent == nil {
		panic("no global scope found")
	}
	return s.parent.Global()
}

func (s *CScope) Owner() ScopeOwner {
	return s.owner
}

func (s *CScope) Parent() *CScope {
	return s.parent
}

func (s *CScope) ModuleScope() *CScope {
	if s.owner == ScopeOwnerModule {
		return s
	}
	if s.parent == nil {
		panic("no module scope found")
	}
	return s.parent.ModuleScope()
}

func (s *CScope) Module() *CModule {
	if s.module != nil {
		return s.module
	}
	if s.parent == nil {
		panic("no module found")
	}
	return s.parent.Module()
}

func (s *CScope) SetModule(module *CModule) {
	s.module = module
}

func (s *CScope) Fork(owner ScopeOwner) *CScope {
	return &CScope{
		owner:  owner,
		parent: s,
		idents: make(map[string]*CVariable),
	}
}

func (s *CScope) Get(name string) *CVariable {
	if v, ok := s.idents[name]; ok {
		return v
	}
	if s.parent != nil {
		return s.parent.Get(name)
	}
	return nil
}

func (s *CScope) Idents() map[string]*CVariable {
	return s.idents
}

func (s *CScope) GetValue(name string) value.Value {
	return s.Get(name).Value
}

func (s *CScope) Set(name string, value *CVariable) error {
	if s.Builtins().IsDefined(name) {
		return fmt.Errorf("cannot redefine builtin %s", name)
	}

	if _, ok := s.idents[name]; ok {
		return fmt.Errorf("%s is already defined", name)
	}

	s.idents[name] = value
	return nil
}

func (s *CScope) Update(name string, value value.Value) {
	// try to look up the name in the current scope, if it's not there, look up in the parent scope

	if v, ok := s.idents[name]; ok {
		if v.Constant {
			panic(fmt.Errorf("cannot assign to constant %s", name))
		}
		s.idents[name].Value = value
		return
	}

	if s.parent != nil {
		s.parent.Update(name, value)
	} else {
		panic(fmt.Errorf("%s is not defined", name))
	}
}

func (s *CScope) IsDefined(name string) bool {
	if _, ok := s.idents[name]; ok {
		return true
	}
	if s.parent != nil {
		return s.parent.IsDefined(name)
	}
	return false
}

func (s *CScope) IsDefinedInCurrentScope(name string) bool {
	if _, ok := s.idents[name]; ok {
		return true
	}
	return false
}

func (s *CScope) Clone() *CScope {
	parent := s.parent
	if parent != nil {
		parent = parent.Clone()
	}
	return &CScope{
		owner:  s.owner,
		idents: s.idents,
		module: s.module,
		parent: parent,
	}
}

func (s *CScope) Reparent(parent *CScope) *CScope {
	clone := s.Clone()
	// put parent at the top of the scope hierarchy
	current := clone
	for current.parent != nil {
		current = current.parent
	}
	current.parent = parent
	return clone
}
