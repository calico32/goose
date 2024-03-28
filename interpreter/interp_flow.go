package interpreter

import (
	"github.com/calico32/goose/ast"
	. "github.com/calico32/goose/interpreter/lib"
)

func (i *interp) evalIfExpr(scope *Scope, expr *ast.IfExpr) Value {
	val := i.evalExpr(scope, expr.Cond)

	if IsTruthy(val) {
		return i.evalExpr(scope, expr.Then)
	} else {
		return i.evalExpr(scope, expr.Else)
	}
}

func (i *interp) runIfStmt(scope *Scope, stmt *ast.IfStmt) StmtResult {
	defer un(trace(i, "if stmt"))
	cond := i.evalExpr(scope, stmt.Cond)

	if IsTruthy(cond) {
		result := i.runStmts(scope, stmt.Body)
		switch result.(type) {
		case *Return, *Break, *Continue:
			return result
		}
	} else if stmt.Else != nil && len(stmt.Else) > 0 {
		result := i.runStmts(scope, stmt.Else)
		switch result.(type) {
		case *Return, *Break, *Continue:
			return result
		}
	}

	return &Void{}
}

func (i *interp) evalDoExpr(scope *Scope, expr *ast.DoExpr) (ret Value) {
	defer un(trace(i, "do expr"))

	ret = NullValue

	newScope := scope.Fork(ScopeOwnerDo)

	result := i.runStmts(newScope, expr.Body)
	switch result := result.(type) {
	case *Break, *Continue:
		i.Throw("break/continue not allowed in do expression") // TODO: support
	case *Return:
		i.Throw("return not allowed in do expression") // TODO: support
	case *LoneValue:
		ret = result.Value
	}

	return
}
