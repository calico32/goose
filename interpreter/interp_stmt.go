package interpreter

import (
	"github.com/calico32/goose/ast"
	. "github.com/calico32/goose/interpreter/lib"
)

func (i *interp) runStmt(scope *Scope, stmt ast.Stmt) StmtResult {
	defer un(trace(i, "stmt"))
	defer pop(push(i, stmt))

	switch stmt := stmt.(type) {
	case *ast.RepeatCountStmt:
		return i.runRepeatCountStmt(scope, stmt)
	case *ast.RepeatWhileStmt:
		return i.runRepeatWhileStmt(scope, stmt)
	case *ast.RepeatForeverStmt:
		return i.runRepeatForeverStmt(scope, stmt)
	case *ast.ForStmt:
		return i.runForStmt(scope, stmt)
	case *ast.IfStmt:
		return i.runIfStmt(scope, stmt)
	case *ast.ReturnStmt:
		return i.runReturnStmt(scope, stmt)
	case *ast.ConstStmt:
		return i.runConstStmt(scope, stmt)
	case *ast.LetStmt:
		return i.runLetStmt(scope, stmt)
	case *ast.AssignStmt:
		return i.runAssignStmt(scope, stmt)
	case *ast.ExprStmt:
		return i.runExprStmt(scope, stmt)
	case *ast.BranchStmt:
		return i.runBranchStmt(scope, stmt)
	case *ast.IncDecStmt:
		return i.runIncDecStmt(scope, stmt)
	case *ast.StructStmt:
		return i.runStructStmt(scope, stmt)
	case *ast.ExportDeclStmt:
		return i.runExportDeclStmt(scope, stmt)
	case *ast.ExportListStmt:
		return i.runExportListStmt(scope, stmt)
	case *ast.ExportSpecStmt:
		return i.runExportSpecStmt(scope, stmt)
	case *ast.ImportStmt:
		return i.runImportStmt(scope, stmt)
	case ast.NativeStmt:
		return i.runNativeStmt(scope, stmt)
	case *ast.SymbolStmt:
		return i.runSymbolStmt(scope, stmt)
	default:
		i.Throw("unexpected statement type %T", stmt)
		return nil
	}
}

func (i *interp) runStmts(scope *Scope, body []ast.Stmt) StmtResult {
	var last StmtResult
	for _, stmt := range body {
		result := i.runStmt(scope, stmt)
		switch result.(type) {
		case *Return, *Break, *Continue:
			return result
		}
		last = result
	}

	if val, ok := last.(*LoneValue); ok {
		return val
	}

	return &Void{}
}

func (i *interp) runExprStmt(scope *Scope, stmt *ast.ExprStmt) StmtResult {
	defer un(trace(i, "expr stmt"))
	val := i.evalExpr(scope, stmt.X)

	if fn, ok := stmt.X.(*ast.FuncExpr); ok && fn.Name != nil && fn.Receiver == nil {
		// special case: if the expression is a function expression with a name, mark it as a declaration
		return &Decl{
			Name:  fn.Name.Name,
			Value: scope.Get(fn.Name.Name).Value,
		}
	}

	return &LoneValue{Value: val}
}
