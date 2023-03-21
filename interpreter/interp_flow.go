package interpreter

import "github.com/calico32/goose/ast"

func (i *interp) evalIfExpr(scope *Scope, expr *ast.IfExpr) Value {
	val := i.evalExpr(scope, expr.Cond)

	if isTruthy(val) {
		return i.evalExpr(scope, expr.Then)
	} else {
		return i.evalExpr(scope, expr.Else)
	}
}

func (i *interp) runIfStmt(scope *Scope, stmt *ast.IfStmt) StmtResult {
	defer un(trace(i, "if stmt"))
	cond := i.evalExpr(scope, stmt.Cond)

	if isTruthy(cond) {
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

	result := i.runStmts(scope, expr.Body)
	switch result := result.(type) {
	case *Break, *Continue:
		i.throw("cannot branch from inside do expr")
	case *Return:
		ret = result.Value
	}

	return
}
