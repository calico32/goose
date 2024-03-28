package interpreter

import (
	"math/rand"

	"github.com/calico32/goose/ast"
	. "github.com/calico32/goose/interpreter/lib"
)

func (i *interp) runLetStmt(scope *Scope, stmt *ast.LetStmt) StmtResult {
	defer un(trace(i, "let stmt"))

	if stmt.Ident.Name == "_" {
		i.Throw("cannot declare _")
	}

	if scope.IsDefinedInCurrentScope(stmt.Ident.Name) {
		i.Throw("cannot redefine variable %s", stmt.Ident.Name)
	}

	var value Value

	if stmt.Value != nil {
		value = i.evalExpr(scope, stmt.Value)
	}

	if value == nil {
		value = NullValue
	}

	scope.Set(stmt.Ident.Name, &Variable{
		Constant: false,
		Value:    value,
	})

	return &Decl{
		Name:  stmt.Ident.Name,
		Value: value,
	}
}

func (i *interp) runConstStmt(scope *Scope, stmt *ast.ConstStmt) StmtResult {
	defer un(trace(i, "const stmt"))

	if stmt.Ident.Name == "_" {
		i.Throw("cannot declare _")
	}

	if scope.IsDefinedInCurrentScope(stmt.Ident.Name) {
		i.Throw("cannot redefine variable %s", stmt.Ident.Name)
	}

	value := i.evalExpr(scope, stmt.Value)

	scope.Set(stmt.Ident.Name, &Variable{
		Constant: true,
		Value:    value,
	})

	return &Decl{
		Name:  stmt.Ident.Name,
		Value: value,
	}
}

func (i *interp) runAssignStmt(scope *Scope, stmt *ast.AssignStmt) StmtResult {
	defer un(trace(i, "assign stmt"))
	// evaluate value

	switch lhs := stmt.Lhs.(type) {
	case *ast.Ident:
		ident := lhs.Name

		if ident[0] == '#' {
			x := scope.Get("this")
			if x == nil {
				i.Throw("invalid property assignment: 'this' is not defined")
			}
			ident = ident[1:]

			canAssign := false
			if composite, ok := x.Value.(*Composite); ok {
				if composite.Frozen {
					i.Throw("cannot assign to frozen composite")
				}

				canAssign = true
			}

			if !canAssign {
				i.Throw("cannot assign to property %s of type %s", ident, x.Value.Type())
			}

			rhs := i.evalExpr(scope, stmt.Rhs)

			SetProperty(x.Value, NewString(ident), rhs)

			return &Void{}
		}

		existing := scope.Get(ident)
		if existing == nil {
			i.Throw("%s is not defined", ident)
		}
		if existing.Constant {
			i.Throw("cannot assign to constant %s", ident)
		}
		op := GetOperator(existing.Value, stmt.Tok)
		if op == nil {
			i.Throw("operator %s not defined for type %s", stmt.Tok, existing.Value.Type())
		}
		rhs := i.evalExpr(scope, stmt.Rhs)
		newValue := op.Executor(&FuncContext{
			Interp: i,
			Scope:  scope,
			This:   existing.Value,
			Args:   []Value{rhs},
		})

		scope.Update(ident, newValue.Value)
		return &Void{}

	case *ast.SelectorExpr:
		existing := i.evalExpr(scope, lhs.X)

		canAssign := false
		if composite, ok := existing.(*Composite); ok {
			if composite.Frozen {
				i.Throw("cannot assign to frozen composite")
			}

			canAssign = true
		}

		sel := lhs.Sel.Name

		if !canAssign {
			i.Throw("cannot assign to property %s of type %s", lhs.Sel.Name, existing.Type())
		}

		rhs := i.evalExpr(scope, stmt.Rhs)

		SetProperty(existing, NewString(sel), rhs)
	case *ast.BracketSelectorExpr:
		existing := i.evalExpr(scope, lhs.X)

		if composite, ok := existing.(*Composite); ok {
			if composite.Frozen {
				i.Throw("cannot assign to frozen composite")
			}
		}
		sel := i.evalExpr(scope, lhs.Sel)

		if _, ok := existing.(*Array); ok {
			if _, ok := sel.(*Integer); !ok {
				i.Throw("cannot index array with type %s", sel.Type())
			}
		}

		if _, ok := sel.(PropertyKey); !ok {
			i.Throw("cannot index with type %s", sel.Type())
		}

		rhs := i.evalExpr(scope, stmt.Rhs)

		SetProperty(existing, sel.(PropertyKey), rhs)
	}

	return &Void{}
}

func (i *interp) runIncDecStmt(scope *Scope, stmt *ast.IncDecStmt) StmtResult {
	defer un(trace(i, "inc/dec stmt"))
	switch lhs := stmt.X.(type) {
	case *ast.Ident:
		ident := lhs.Name
		existing := scope.Get(ident)
		if existing == nil {
			i.Throw("%s is not defined", ident)
		}
		if existing.Constant {
			i.Throw("cannot assign to constant %s", ident)
		}
		if _, ok := existing.Value.(Numeric); !ok {
			i.Throw("cannot increment/decrement non-number %s", ident)
		}

		op := GetOperator(existing.Value, stmt.Tok)
		if op == nil {
			i.Throw("cannot increment/decrement %s of type %s", ident, existing.Value.Type())
		}
		newValue := op.Executor(&FuncContext{
			Interp: i,
			Scope:  scope,
			This:   existing.Value,
			Args:   []Value{Wrap(1)},
		})

		scope.Update(ident, newValue.Value)
	case *ast.BracketSelectorExpr:
		obj := i.evalExpr(scope, lhs.X)

		if composite, ok := obj.(*Composite); ok {
			if composite.Frozen {
				i.Throw("cannot assign to frozen composite")
			}
		}
		sel := i.evalExpr(scope, lhs.Sel)

		if _, ok := obj.(*Array); ok {
			if _, ok := sel.(*Integer); !ok {
				i.Throw("cannot index array with type %s", sel.Type())
			}
		}

		if _, ok := sel.(PropertyKey); !ok {
			i.Throw("cannot index with type %s", sel.Type())
		}

		existing := GetProperty(obj, sel.(PropertyKey))

		if _, ok := existing.(Numeric); !ok {
			i.Throw("cannot increment non-numeric value %s", existing)
		}

		op := GetOperator(existing, stmt.Tok)
		if op == nil {
			i.Throw("cannot increment %s with %s", existing, stmt.Tok)
		}
		newValue := op.Executor(&FuncContext{
			Interp: i,
			Scope:  scope,
			This:   existing,
			Args:   []Value{Wrap(1)},
		})

		SetProperty(obj, sel.(PropertyKey), newValue.Value)
	}
	return &Void{}
}

func (i *interp) runSymbolStmt(scope *Scope, stmt *ast.SymbolStmt) StmtResult {
	defer un(trace(i, "symbol stmt"))

	if scope.IsDefinedInCurrentScope(stmt.Ident.Name) {
		i.Throw("cannot redefine symbols %s", stmt.Ident.Name)
	}

	symbol := &Symbol{
		Name: stmt.Ident.Name,
		Id:   rand.Int63(),
	}

	scope.Set(stmt.Ident.Name, &Variable{
		Constant: true,
		Value:    symbol,
	})

	return &Decl{
		Name:  stmt.Ident.Name,
		Value: symbol,
	}
}
