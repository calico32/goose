package interpreter

import (
	"fmt"

	"github.com/calico32/goose/ast"
	. "github.com/calico32/goose/interpreter/lib"
)

func (i *interp) runStructStmt(scope *Scope, stmt *ast.StructStmt) StmtResult {
	defer un(trace(i, "struct stmt"))

	if scope.IsDefinedInCurrentScope(stmt.Name.Name) {
		i.Throw("cannot redefine variable %s", stmt.Name.Name)
	}

	// validate parameters
	fieldNames := map[string]bool{}
	for _, param := range stmt.Fields.List {
		if fieldNames[param.Ident.Name] {
			i.Throw("duplicate field %s", param.Ident.Name)

		}
		fieldNames[param.Ident.Name] = true
	}

	var fieldDefaults []Value
	for _, field := range stmt.Fields.List {
		var v Value
		if field.Value != nil {
			v = i.evalExpr(scope, field.Value)
		} else {
			v = NullValue
		}
		fieldDefaults = append(fieldDefaults, v.Clone())
	}

	proto := NewComposite()
	proto.Name = stmt.Name.Name

	closure := scope.Fork(ScopeOwnerClosure)

	var executor FuncType = func(ctx *FuncContext) *Return {
		// create new scope
		// TODO: closures
		newScope := closure.Fork(ScopeOwnerStruct)

		// create new composite
		obj := &Composite{
			Proto:      proto,
			Properties: make(Properties),
			Operators:  make(Operators),
		}

		// set parameters in scope
		for idx, param := range stmt.Fields.List {
			var v Value
			if idx < len(ctx.Args) {
				v = ctx.Args[idx].Clone()
			} else {
				v = fieldDefaults[idx].Clone()
			}

			if stmt.Init != nil {
				newScope.Set(param.Ident.Name, &Variable{
					Constant: false,
					Value:    v,
				})
			}

			if set, ok := obj.Properties[PKString]; ok {
				set[param.Ident.Name] = v
			} else {
				obj.Properties[PKString] = map[string]Value{
					param.Ident.Name: v,
				}
			}
		}

		if stmt.Init != nil {
			// set this
			// TODO: better this
			newScope.Set("this", &Variable{
				Constant: false,
				Value:    obj,
			})

			result := i.runStmts(newScope, stmt.Init.Body)
			switch result.(type) {
			case *Return, *Break, *Continue:
				i.Throw("cannot return or branch from struct initializer")
			}
		}

		return NewReturn(obj)
	}

	value := &Func{
		NewableProto: proto,
		Executor:     executor,
	}

	scope.Set(stmt.Name.Name, &Variable{
		Constant: false,
		Value:    value,
	})

	return &Decl{
		Name:  stmt.Name.Name,
		Value: value,
	}
}

func (i *interp) runOperatorStmt(scope *Scope, stmt *ast.OperatorStmt) StmtResult {
	defer un(trace(i, "operator stmt"))

	constructor := scope.Get(stmt.Receiver.Name)
	if constructor == nil {
		i.Throw("operator receiver %s not found", stmt.Receiver.Name)
	}

	if receiver, ok := constructor.Value.(*Func); !ok || receiver.NewableProto == nil {
		i.Throw("operator receiver %s is not a valid struct", stmt.Receiver.Name)
	}

	proto := constructor.Value.(*Func).NewableProto

	// check if operator already exists
	if _, ok := proto.Operators[stmt.Tok]; ok {
		i.Throw("operator %s already defined for %s", stmt.Tok, stmt.Receiver.Name)
	}

	// validate parameters
	paramNames := map[string]bool{}
	for _, param := range stmt.Params.List {
		if paramNames[param.Ident.Name] {
			i.Throw("duplicate parameter %s", param.Ident.Name)
		}
		paramNames[param.Ident.Name] = true
	}

	var memoCache map[string]*Return
	if stmt.Memo.IsValid() {
		memoCache = make(map[string]*Return)
	}

	// create new scope
	closure := scope.Fork(ScopeOwnerClosure)

	var executor FuncType = func(ctx *FuncContext) (ret *Return) {
		// create new scope
		opScope := closure.Fork(ScopeOwnerOperator)

		// use memo cache if applicable
		if stmt.Memo.IsValid() {
			// hash the arguments
			hash := ""
			for _, arg := range ctx.Args {
				hash += fmt.Sprintf("%s|%v,", arg.Type(), arg.Hash())
			}
			hash = hash[:len(hash)-1]

			// check cache
			if memoCache[hash] != nil {
				return memoCache[hash]
			}

			// cache miss, store the result later
			defer func() {
				memoCache[hash] = ret
			}()
		}

		// set parameters in scope
		for idx, param := range stmt.Params.List {
			var v Value
			if idx < len(ctx.Args) {
				v = ctx.Args[idx].Clone()
			} else {
				v = NullValue
			}

			opScope.Set(param.Ident.Name, &Variable{
				Constant: false,
				Value:    v,
			})
		}

		// set this
		opScope.Set("this", &Variable{
			Constant: false,
			Value:    ctx.This,
		})

		// run operator
		result := i.runStmts(opScope, stmt.Body)
		switch result := result.(type) {
		case *Return:
			return result
		case *Break, *Continue:
			i.Throw("cannot branch from operator")
		}

		return NewReturn(NullValue)
	}

	value := &OperatorFunc{
		Builtin:  false,
		Executor: executor,
	}

	proto.Operators[stmt.Tok] = value

	return &Void{}
}
