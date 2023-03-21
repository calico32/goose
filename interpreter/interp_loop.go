package interpreter

import (
	"github.com/calico32/goose/ast"
	"github.com/calico32/goose/token"
)

func (i *interp) runRepeatCountStmt(scope *Scope, stmt *ast.RepeatCountStmt) StmtResult {
	defer un(trace(i, "repeat count stmt"))
	count := int64(0)
	totalCount := i.evalExpr(scope, stmt.Count)

	if _, ok := totalCount.(Numeric); !ok {
		i.throw("repeat count must be numeric")
	}

	for count < totalCount.(Numeric).Int64() {
		repeatScope := scope.Fork(ScopeOwnerRepeat)
		result := i.runStmts(repeatScope, stmt.Body)
		switch result := result.(type) {
		case *Return:
			return result
		case *Break:
			return &Void{}
		case *Continue:
			// continue
		}
		count++
	}

	return &Void{}
}

func (i *interp) runBranchStmt(scope *Scope, stmt *ast.BranchStmt) StmtResult {
	defer un(trace(i, "branch stmt"))
	switch stmt.Tok {
	case token.Break:
		return &Break{}
	case token.Continue:
		return &Continue{}
	default:
		i.throw("unexpected branch type %v", stmt.Tok)
	}

	return nil
}

func (i *interp) runRepeatWhileStmt(scope *Scope, stmt *ast.RepeatWhileStmt) StmtResult {
	defer un(trace(i, "repeat while stmt"))
	for {
		cond := i.evalExpr(scope, stmt.Cond)
		if !isTruthy(cond) {
			break
		}

		repeatScope := scope.Fork(ScopeOwnerRepeat)
		result := i.runStmts(repeatScope, stmt.Body)
		switch result.(type) {
		case *Return:
			return result
		case *Break:
			return &Void{}
		case *Continue:
			// continue
		}
	}

	return &Void{}
}

func (i *interp) runRepeatForeverStmt(scope *Scope, stmt *ast.RepeatForeverStmt) StmtResult {
	defer un(trace(i, "repeat forever stmt"))
	for {
		repeatScope := scope.Fork(ScopeOwnerRepeat)
		result := i.runStmts(repeatScope, stmt.Body)
		switch result.(type) {
		case *Return:
			return result
		case *Break:
			return &Void{}
		case *Continue:
			// continue
		}
	}
}

func (interp *interp) runForStmt(scope *Scope, stmt *ast.ForStmt) StmtResult {
	defer un(trace(interp, "for stmt"))

	iterable := interp.evalExpr(scope, stmt.Iterable)

	var iterVal []Value

	switch iterable := iterable.(type) {
	case *String:
		for _, char := range iterable.Value {
			iterVal = append(iterVal, wrap(string(char)))
		}
	case *Array:
		iterVal = iterable.Elements
	default:
		interp.throw("for loop iterable must be... iterable")
	}

	name := stmt.Var.Name
	for i := 0; i < len(iterVal); i++ {
		forScope := scope.Fork(ScopeOwnerFor)
		forScope.Set(name, &Variable{
			Constant: false,
			Value:    iterVal[i],
		})

		result := interp.runStmts(forScope, stmt.Body)
		switch result.(type) {
		case *Return:
			return result
		case *Break:
			return &Void{}
		case *Continue:
			// continue
		}
	}

	return &Void{}
}
