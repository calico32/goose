package interpreter

import (
	"github.com/calico32/goose/ast"
	. "github.com/calico32/goose/interpreter/lib"
)

func (i *interp) evalGeneratorExpr(scope *Scope, expr *ast.GeneratorExpr) Value {
	defer un(trace(i, "generator expr"))

	if expr.Name != nil && expr.Receiver == nil {
		if value := scope.Get(expr.Name.Name); value != nil {
			i.Throw("duplicate generator %s", expr.Name.Name)
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

	var executor = func(ctx *FuncContext, channel chan GeneratorMessage) {
		defer func() {
			if err := recover(); err != nil {
				channel <- &GeneratorError{Error: err.(error)}
			}
		}()

		genScope := scope.Fork(ScopeOwnerGenerator)

		// set parameters in scope
		for idx, param := range expr.Params.List {
			var v Value
			if idx < len(ctx.Args) {
				v = ctx.Args[idx].Clone()
			} else {
				v = paramDefaults[idx].Clone()
			}

			genScope.Set(param.Ident.Name, &Variable{
				Constant: false,
				Value:    v,
			})
		}

		// TODO: better this
		genScope.Set("this", &Variable{
			Constant: true,
			Value:    ctx.This,
		})

		for {
			message, ok := <-channel
			if !ok {
				return
			}

			ip := 0
			stmts := expr.Body

			switch message.(type) {
			// TODO
			case *GeneratorNext:
				// run stmts
				for ip < len(stmts) {

				}
			case *GeneratorDone:
				m := &GeneratorIsDone{Done: ip >= len(stmts)}
				channel <- m
			}
		}
	}

	var factoryFunc FuncType = func(ctx *FuncContext) *Return {
		channel := make(chan GeneratorMessage)
		generator := &Generator{
			Async:   expr.Async.IsValid(),
			Channel: channel,
		}

		go executor(ctx, channel)

		return NewReturn(generator)
	}

	factory := &Func{Executor: factoryFunc}

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

			proto.Properties[PKString][expr.Name.Name] = factory
		} else {
			scope.Set(expr.Name.Name, &Variable{
				Constant: true, // functions are constants
				Value:    factory,
			})
		}
	}

	return factory
}
