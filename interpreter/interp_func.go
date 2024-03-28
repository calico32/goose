package interpreter

import (
	"fmt"

	"github.com/calico32/goose/ast"

	. "github.com/calico32/goose/interpreter/lib"
)

func (i *interp) evalFuncExpr(scope *Scope, expr *ast.FuncExpr) Value {
	defer un(trace(i, "func expr"))

	if expr.Name != nil && expr.Receiver == nil {
		if scope.IsDefinedInCurrentScope(expr.Name.Name) {
			i.Throw("name %s is already defined", expr.Name.Name)
		}
	}

	// validate parameters
	paramNames := map[string]bool{}
	for _, param := range expr.Params.List {
		if paramNames[param.Ident.Name] {
			i.Throw("duplicate parameter %s", param.Ident.Name)
		}
		paramNames[param.Ident.Name] = true
	}

	var memoCache map[string]*Return
	if expr.Memo.IsValid() {
		memoCache = make(map[string]*Return)
	}

	var paramDefaults []Value
	for _, param := range expr.Params.List {
		var v Value
		if param.Value != nil {
			v = i.evalExpr(scope, param.Value)
		} else {
			v = NullValue
		}
		paramDefaults = append(paramDefaults, v.Clone())
	}

	closure := scope.Fork(ScopeOwnerClosure)

	// TODO: async

	var executor FuncType = func(ctx *FuncContext) (ret *Return) {
		// create new scope
		funcScope := closure.Fork(ScopeOwnerFunc)

		// use memo cache if applicable
		if expr.Memo.IsValid() {
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
		for idx, param := range expr.Params.List {
			var v Value
			if idx < len(ctx.Args) {
				v = ctx.Args[idx].Clone()
			} else {
				v = paramDefaults[idx].Clone()
			}

			funcScope.Set(param.Ident.Name, &Variable{
				Constant: false,
				Value:    v,
			})
		}

		// TODO: better this
		funcScope.Set("this", &Variable{
			Constant: true,
			Value:    ctx.This,
		})

		if expr.Arrow.IsValid() {
			result := i.evalExpr(funcScope, expr.ArrowExpr)
			return NewReturn(&result)
		}

		result := i.runStmts(funcScope, expr.Body)
		switch result := result.(type) {
		case *Return:
			return result
		case *Break, *Continue, *Yield:
			i.Throw("cannot branch from function")
		}

		return NewReturn(NullValue)
	}

	value := &Func{
		Async:    expr.Async.IsValid(),
		Memoized: expr.Memo.IsValid(),
		Executor: executor,
	}

	if expr.Name != nil {
		if expr.Receiver != nil {
			// find proto
			// TODO: limit to current module
			constructor := scope.Get(expr.Receiver.Name)
			if constructor == nil {
				i.Throw("unknown type %s", expr.Receiver.Name)
			}

			if val, ok := constructor.Value.(*Func); !ok || val.NewableProto == nil {
				i.Throw("%s is not a type", expr.Receiver.Name)
			}

			proto := constructor.Value.(*Func).NewableProto

			if proto.Properties[PKString] == nil {
				proto.Properties[PKString] = make(map[string]Value)
			}

			if _, ok := proto.Properties[PKString][expr.Name.Name]; ok {
				i.Throw("duplicate receiver function %s", expr.Name.Name)
			}

			proto.Properties[PKString][expr.Name.Name] = value
		} else {
			scope.Set(expr.Name.Name, &Variable{
				Constant: true, // functions are constants
				Value:    value,
			})
		}
	}

	return value
}

func (i *interp) evalCallExpr(scope *Scope, expr *ast.CallExpr) Value {
	defer un(trace(i, "call expr"))

	var this Value
	var fn Value
	switch fexpr := expr.Func.(type) {
	case *ast.SelectorExpr:
		this = i.evalExpr(scope, fexpr.X)
		fn = GetProperty(this, NewString(fexpr.Sel.Name)) // TODO: check type
	case *ast.BracketSelectorExpr:
		this = i.evalExpr(scope, fexpr.X)
		sel := i.evalExpr(scope, fexpr.Sel)
		fn = GetProperty(this, sel.(PropertyKey)) // TODO: check type
	default:
		fn = i.evalExpr(scope, expr.Func)
		if f, ok := fn.(*Func); ok && f.This != nil {
			this = f.This
		}
	}

	if _, ok := fn.(*Func); !ok {
		i.Throw("expression of type %s is not callable", fn.Type())
	}

	args := make([]Value, len(expr.Args))
	for idx, arg := range expr.Args {
		val := i.evalExpr(scope, arg)

		args[idx] = val
	}

	result := fn.(*Func).Executor(&FuncContext{
		Interp: i,
		Scope:  scope,
		This:   this,
		Args:   args,
	})

	return result.Value
}

func (i *interp) evalBindExpr(scope *Scope, expr *ast.BindExpr) Value {
	defer un(trace(i, "bind expr"))

	x := i.evalExpr(scope, expr.X)

	rhs := i.evalExpr(scope, expr.Sel)

	if rhs, ok := rhs.(*Func); ok {
		rhs.This = x
		return rhs
	} else {
		i.Throw("cannot bind non-function value")
		return nil
	}
}

func (i *interp) runReturnStmt(scope *Scope, stmt *ast.ReturnStmt) StmtResult {
	defer un(trace(i, "return stmt"))

	var ret Value

	if stmt.Result != nil {
		ret = i.evalExpr(scope, stmt.Result)
	}

	return NewReturn(&ret)
}
