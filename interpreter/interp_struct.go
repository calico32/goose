package interpreter

import "github.com/calico32/goose/ast"

func (i *interp) runStructStmt(scope *Scope, stmt *ast.StructStmt) StmtResult {
	defer un(trace(i, "struct stmt"))

	if scope.IsDefinedInCurrentScope(stmt.Name.Name) {
		i.throw("cannot redefine variable %s", stmt.Name.Name)
	}

	// validate parameters
	fieldNames := map[string]bool{}
	for _, param := range stmt.Fields.List {
		if fieldNames[param.Ident.Name] {
			i.throw("duplicate field %s", param.Ident.Name)

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

	var executor FuncType = func(ctx *FuncContext) *Return {
		// create new scope
		// TODO: closures
		newScope := ctx.Scope.Fork(ScopeOwnerStruct)

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

			if set, ok := obj.Properties[PropertyKeyString]; ok {
				set[param.Ident.Name] = v
			} else {
				obj.Properties[PropertyKeyString] = map[string]Value{
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
				i.throw("cannot return or branch from struct initializer")
			}
		}

		return &Return{obj}
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
