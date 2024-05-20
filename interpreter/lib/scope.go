package lib

import "fmt"

type ScopeOwner int

const (
	ScopeOwnerFunc ScopeOwner = iota
	ScopeOwnerClosure
	ScopeOwnerDo
	ScopeOwnerArrayInit
	ScopeOwnerRepeat
	ScopeOwnerFor
	ScopeOwnerIf
	ScopeOwnerGlobal
	ScopeOwnerBuiltin
	ScopeOwnerModule
	ScopeOwnerPipeline
	ScopeOwnerBlock
	ScopeOwnerStruct
	ScopeOwnerGenerator
	ScopeOwnerImport
	ScopeOwnerMatch
	ScopeOwnerOperator
)

var scopeOwnerNames = [...]string{
	ScopeOwnerFunc:      "function",
	ScopeOwnerClosure:   "closure",
	ScopeOwnerDo:        "do",
	ScopeOwnerArrayInit: "array initializer",
	ScopeOwnerRepeat:    "repeat",
	ScopeOwnerFor:       "for",
	ScopeOwnerIf:        "if",
	ScopeOwnerGlobal:    "global",
	ScopeOwnerBuiltin:   "builtin",
	ScopeOwnerModule:    "module",
	ScopeOwnerPipeline:  "pipeline",
	ScopeOwnerBlock:     "block",
	ScopeOwnerStruct:    "struct",
	ScopeOwnerGenerator: "generator",
	ScopeOwnerImport:    "import",
	ScopeOwnerMatch:     "match",
	ScopeOwnerOperator:  "operator",
}

// scope hierarchy:
// builtin -> global -> module -> any other scope

func (s ScopeOwner) String() string {
	return scopeOwnerNames[s]
}

type Scope struct {
	module *Module
	owner  ScopeOwner
	parent *Scope
	idents map[string]*Variable
}

type Variable struct {
	Constant bool
	Value    Value
	Source   VariableSource
}

type VariableSource int

const (
	VariableSourceDecl VariableSource = iota
	VariableSourceImport
)

func NewGlobalScope(builtins map[string]*Variable) *Scope {
	return &Scope{
		owner:  ScopeOwnerGlobal,
		idents: make(map[string]*Variable),
		parent: &Scope{
			owner:  ScopeOwnerBuiltin,
			parent: nil,
			idents: builtins,
		},
	}
}

func NewScope(owner ScopeOwner) *Scope {
	return &Scope{
		owner:  owner,
		idents: make(map[string]*Variable),
	}
}

func (s *Scope) Builtins() *Scope {
	if s.owner == ScopeOwnerBuiltin {
		return s
	}
	if s.parent == nil {
		panic("no builtin scope found")
	}
	return s.parent.Builtins()
}

func (s *Scope) Global() *Scope {
	if s.owner == ScopeOwnerGlobal {
		return s
	}
	if s.parent == nil {
		panic("no global scope found")
	}
	return s.parent.Global()
}

func (s *Scope) Owner() ScopeOwner {
	return s.owner
}

func (s *Scope) Parent() *Scope {
	return s.parent
}

func (s *Scope) ModuleScope() *Scope {
	if s.owner == ScopeOwnerModule {
		return s
	}
	if s.parent == nil {
		panic("no module scope found")
	}
	return s.parent.ModuleScope()
}

func (s *Scope) Module() *Module {
	if s.module != nil {
		return s.module
	}
	if s.parent == nil {
		panic("no module found")
	}
	return s.parent.Module()
}

func (s *Scope) SetModule(module *Module) {
	s.module = module
}

func (s *Scope) Fork(owner ScopeOwner) *Scope {
	return &Scope{
		owner:  owner,
		parent: s,
		idents: make(map[string]*Variable),
	}
}

func (s *Scope) Get(name string) *Variable {
	if v, ok := s.idents[name]; ok {
		return v
	}
	if s.parent != nil {
		return s.parent.Get(name)
	}
	return nil
}

func (s *Scope) Idents() map[string]*Variable {
	return s.idents
}

func (s *Scope) GetValue(name string) Value {
	return s.Get(name).Value
}

func (s *Scope) Set(name string, value *Variable) {
	if s.Builtins().IsDefined(name) {
		panic(fmt.Errorf("cannot redefine builtin %s", name))
	}

	if v, ok := s.idents[name]; ok {
		if v.Constant {
			panic(fmt.Errorf("cannot assign to constant %s", name))
		}
	}

	s.idents[name] = value
}

func (s *Scope) Update(name string, value Value) {
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

func (s *Scope) IsDefined(name string) bool {
	if _, ok := s.idents[name]; ok {
		return true
	}
	if s.parent != nil {
		return s.parent.IsDefined(name)
	}
	return false
}

func (s *Scope) IsDefinedInCurrentScope(name string) bool {
	if _, ok := s.idents[name]; ok {
		return true
	}
	return false
}

func (s *Scope) Clone() *Scope {
	parent := s.parent
	if parent != nil {
		parent = parent.Clone()
	}
	return &Scope{
		owner:  s.owner,
		idents: s.idents,
		module: s.module,
		parent: parent,
	}
}

func (s *Scope) Reparent(parent *Scope) *Scope {
	clone := s.Clone()
	// put parent at the top of the scope hierarchy
	current := clone
	for current.parent != nil {
		current = current.parent
	}
	current.parent = parent
	return clone
}
