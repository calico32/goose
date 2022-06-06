package interpreter

import (
	"errors"
	"fmt"

	"github.com/wiisportsresort/goose/ast"
	"github.com/wiisportsresort/goose/token"
)

func (i *interpreter) runStmt(scope *GooseScope, stmt ast.Stmt) (StmtResult, error) {
	defer un(trace(i, "stmt"))
	defer pop(push(i, stmt))

	switch stmt := stmt.(type) {
	case *ast.FuncStmt:
		return i.runFuncStmt(scope, stmt)
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
	case *ast.DeclStmt:
		return i.runDeclStmt(scope, stmt)
	case *ast.AssignStmt:
		return i.runAssignStmt(scope, stmt)
	case *ast.ExprStmt:
		return i.runExprStmt(scope, stmt)
	case *ast.BranchStmt:
		return i.runBranchStmt(scope, stmt)
	case *ast.IncDecStmt:
		return i.runIncDecStmt(scope, stmt)
	default:
		return nil, fmt.Errorf("unexpected statement type %T", stmt)
	}
}

func (i *interpreter) runStmts(scope *GooseScope, body []ast.Stmt) (result StmtResult, err error) {
	for _, stmt := range body {
		result, err := i.runStmt(scope, stmt)
		if err != nil {
			return nil, err
		}

		switch result.(type) {
		case *ReturnResult, *BreakResult, *ContinueResult:
			return result, nil
		}
	}

	return &VoidResult{}, nil
}

func (i *interpreter) runFuncStmt(scope *GooseScope, stmt *ast.FuncStmt) (result *VoidResult, err error) {
	defer un(trace(i, "func stmt"))
	if value := scope.get(stmt.Name.Name); value != nil {
		err = fmt.Errorf("duplicate function %s", stmt.Name.Name)
		return
	}

	// validate parameters
	paramNames := map[string]bool{}
	for _, param := range stmt.Params.List {
		if paramNames[param.Ident.Name] {
			err = fmt.Errorf("duplicate parameter %s", param.Ident.Name)
			return
		}
		paramNames[param.Ident.Name] = true
	}

	var memoCache map[string]*ReturnResult
	if stmt.Memo.IsValid() {
		memoCache = make(map[string]*ReturnResult)
	}

	var executor GooseFunc = func(scope *GooseScope, args []*GooseValue) (ret *ReturnResult, err error) {
		// create new scope
		funcScope := scope.new(ScopeOwnerFunc)

		// use memo cache if applicable
		if stmt.Memo.IsValid() {
			// hash the arguments
			hash := ""
			for _, arg := range args {
				hash += fmt.Sprintf("%d|%v,", arg.Type, arg.Value)
			}
			hash = hash[:len(hash)-1]

			// check cache
			if memoCache[hash] != nil {
				return memoCache[hash], nil
			}

			// cache miss, store the result later
			defer func() {
				memoCache[hash] = ret
			}()
		}

		// set parameters in scope
		for i, param := range stmt.Params.List {
			v := &GooseValue{Type: GooseTypeNull, Value: nil}
			if i <= len(args) {
				v = args[i].Copy()
			}

			funcScope.set(param.Ident.Name, *v)
		}

		result, err := i.runStmts(funcScope, stmt.Body)
		if err != nil {
			return nil, err
		}
		switch result := result.(type) {
		case *ReturnResult:
			return result, nil
		case *BreakResult, *ContinueResult:
			return nil, errors.New("Cannot branch from function")
		}

		return &ReturnResult{nil}, nil
	}

	scope.set(stmt.Name.Name, GooseValue{
		Constant: true, // functions are constants
		Type:     GooseTypeFunc,
		Value:    executor,
	})

	return &VoidResult{}, nil
}

func (i *interpreter) runRepeatCountStmt(scope *GooseScope, stmt *ast.RepeatCountStmt) (result StmtResult, err error) {
	defer un(trace(i, "repeat count stmt"))
	count := int64(0)
	totalCount, err := i.evalExpr(scope, stmt.Count)
	if err != nil {
		return nil, err
	}

	if totalCount.Type == GooseTypeFloat {
		totalCount = &GooseValue{
			Constant: false,
			Type:     GooseTypeInt,
			Value:    int64(totalCount.Value.(float64)),
		}
	}

	if totalCount.Type != GooseTypeInt {
		return nil, fmt.Errorf("repeat count must be an integer")
	}

	for count < totalCount.Value.(int64) {
		repeatScope := scope.new(ScopeOwnerRepeat)
		result, err = i.runStmts(repeatScope, stmt.Body)
		if err != nil {
			return nil, err
		}
		switch result := result.(type) {
		case *ReturnResult:
			return result, nil
		case *BreakResult:
			return &VoidResult{}, nil
		case *ContinueResult:
			// continue
		}
		count++
	}

	return &VoidResult{}, nil
}

func (i *interpreter) runBranchStmt(scope *GooseScope, stmt *ast.BranchStmt) (result StmtResult, err error) {
	defer un(trace(i, "branch stmt"))
	switch stmt.Tok {
	case token.Break:
		return &BreakResult{}, nil
	case token.Continue:
		return &ContinueResult{}, nil
	default:
		return nil, fmt.Errorf("unexpected branch type %v", stmt.Tok)
	}
}

func (i *interpreter) runRepeatWhileStmt(scope *GooseScope, stmt *ast.RepeatWhileStmt) (result StmtResult, err error) {
	defer un(trace(i, "repeat while stmt"))
	for {
		cond, err := i.evalExpr(scope, stmt.Cond)
		if err != nil {
			return nil, err
		}

		if !isTruthy(cond) {
			break
		}

		repeatScope := scope.new(ScopeOwnerRepeat)
		result, err = i.runStmts(repeatScope, stmt.Body)
		if err != nil {
			return nil, err
		}
		switch result.(type) {
		case *ReturnResult:
			return result, nil
		case *BreakResult:
			return &VoidResult{}, nil
		case *ContinueResult:
			// continue
		}
	}

	return &VoidResult{}, nil
}

func (i *interpreter) runRepeatForeverStmt(scope *GooseScope, stmt *ast.RepeatForeverStmt) (result StmtResult, err error) {
	defer un(trace(i, "repeat forever stmt"))
	for {
		repeatScope := scope.new(ScopeOwnerRepeat)
		result, err = i.runStmts(repeatScope, stmt.Body)
		if err != nil {
			return nil, err
		}
		switch result.(type) {
		case *ReturnResult:
			return result, nil
		case *BreakResult:
			return &VoidResult{}, nil
		case *ContinueResult:
			// continue
		}
	}
}

func (interp *interpreter) runForStmt(scope *GooseScope, stmt *ast.ForStmt) (result StmtResult, err error) {
	defer un(trace(interp, "for stmt"))

	iterable, err := interp.evalExpr(scope, stmt.Iterable)
	if err != nil {
		return nil, err
	}

	var iterVal []*GooseValue

	switch iterable.Type {
	case GooseTypeString:
		chars := []rune(iterable.Value.(string))
		for _, char := range chars {
			iterVal = append(iterVal, &GooseValue{
				Constant: false,
				Type:     GooseTypeString,
				Value:    string(char),
			})
		}
	case GooseTypeArray:
		var ok bool
		iterVal, ok = iterable.Value.([]*GooseValue)
		if !ok {
			return nil, fmt.Errorf("for loop iterable must be...iterable")
		}
	}

	name := stmt.Var.Name
	for i := 0; i < len(iterVal); i++ {
		forScope := scope.new(ScopeOwnerFor)
		forScope.set(name, GooseValue{
			Constant: false,
			Type:     typeOf(iterVal[i]),
			Value:    valueOf(iterVal[i]),
		})

		result, err = interp.runStmts(forScope, stmt.Body)
		if err != nil {
			return nil, err
		}
		switch result.(type) {
		case *ReturnResult:
			return result, nil
		case *BreakResult:
			return &VoidResult{}, nil
		case *ContinueResult:
			// continue
		}
	}

	return &VoidResult{}, nil
}

func (i *interpreter) runIfStmt(scope *GooseScope, stmt *ast.IfStmt) (result StmtResult, err error) {
	defer un(trace(i, "if stmt"))
	cond, err := i.evalExpr(scope, stmt.Cond)
	if err != nil {
		return nil, err
	}

	if isTruthy(cond) {
		result, err = i.runStmts(scope, stmt.Body)
		if err != nil {
			return nil, err
		}
		switch result.(type) {
		case *ReturnResult, *BreakResult, *ContinueResult:
			return result, nil
		}
	} else if stmt.Else != nil && len(stmt.Else) > 0 {
		result, err = i.runStmts(scope, stmt.Else)
		if err != nil {
			return nil, err
		}
		switch result.(type) {
		case *ReturnResult, *BreakResult, *ContinueResult:
			return result, nil
		}
	}

	return &VoidResult{}, nil
}

func (i *interpreter) runReturnStmt(scope *GooseScope, stmt *ast.ReturnStmt) (result StmtResult, err error) {
	defer un(trace(i, "return stmt"))
	ret, err := i.evalExpr(scope, stmt.Result)
	if err != nil {
		return nil, err
	}
	return &ReturnResult{ret}, nil
}

func (i *interpreter) runAssignStmt(scope *GooseScope, stmt *ast.AssignStmt) (result StmtResult, err error) {
	defer un(trace(i, "assign stmt"))
	// evaluate value
	rhs, err := i.evalExpr(scope, stmt.Rhs)
	if err != nil {
		return nil, err
	}

	isIdent, ident, existing, index, err := i.getAssignIdentOrIndex(stmt.Lhs, scope)
	if err != nil {
		return nil, err
	}

	newValue, err := i.evalBinaryValues(existing, stmt.Tok, rhs)
	if err != nil {
		return nil, err
	}

	if isIdent {
		scope.update(ident, *newValue)
	} else {
		existing.Value.([]*GooseValue)[index] = newValue
	}

	return &VoidResult{}, nil
}

func (i *interpreter) getAssignIdentOrIndex(expr ast.Expr, scope *GooseScope) (isIdent bool, ident string, existing *GooseValue, index int, err error) {
	switch expr := expr.(type) {
	case *ast.Ident:
		isIdent = true
		ident = expr.Name
		existing = scope.get(ident)
		if existing == nil {
			err = fmt.Errorf("%s is not defined", ident)
			return
		}
		if existing.Constant {
			err = fmt.Errorf("cannot assign to constant %s", ident)
			return
		}
	case *ast.IndexExpr:
		isIdent = false
		existing, err = i.evalExpr(scope, expr.X)
		if err != nil {
			return
		}
		err = i.expectType(existing, GooseTypeArray)
		if err != nil {
			return
		}
		if existing.Constant {
			err = fmt.Errorf("cannot assign to constant %s", ident)
			return
		}
		var idx *GooseValue
		idx, err = i.evalExpr(scope, expr.Index)
		if err != nil {
			return
		}
		err = i.expectType(idx, GooseTypeNumeric)
		if err != nil {
			return
		}
		index = toInt(idx.Value)
		if index >= len(existing.Value.([]*GooseValue)) {
			err = fmt.Errorf("index %d out of bounds for array of length %d", index, len(existing.Value.([]*GooseValue)))
			return
		}
	default:
		err = fmt.Errorf("left hand side of assignment must be an identifier or index expression")
		return
	}

	return
}

func (i *interpreter) runDeclStmt(scope *GooseScope, stmt *ast.DeclStmt) (result StmtResult, err error) {
	defer un(trace(i, "decl stmt"))

	if stmt.Ident.Name == "_" {
		return nil, fmt.Errorf("cannot declare _")
	}

	if scope.isDefinedInCurrentScope(stmt.Ident.Name) {
		return nil, fmt.Errorf("cannot redefine variable %s", stmt.Ident.Name)
	}

	value, err := i.evalExpr(scope, stmt.Value)
	if err != nil {
		return nil, err
	}

	scope.set(stmt.Ident.Name, GooseValue{
		Constant: stmt.Decl == token.Const,
		Type:     value.Type,
		Value:    value.Value,
	})

	return &VoidResult{}, nil
}

func (i *interpreter) runExprStmt(scope *GooseScope, stmt *ast.ExprStmt) (result *VoidResult, err error) {
	defer un(trace(i, "expr stmt"))
	_, err = i.evalExpr(scope, stmt.X)
	if err != nil {
		return nil, err
	}
	return &VoidResult{}, nil
}

func (i *interpreter) runIncDecStmt(scope *GooseScope, stmt *ast.IncDecStmt) (result *VoidResult, err error) {
	defer un(trace(i, "inc/dec stmt"))
	isIdent, ident, existing, index, err := i.getAssignIdentOrIndex(stmt.X, scope)
	if err != nil {
		return nil, err
	}

	var value *GooseValue
	if isIdent {
		value = existing
	} else {
		value = existing.Value.([]*GooseValue)[index]
	}

	err = i.expectType(value, GooseTypeNumeric)
	if err != nil {
		return
	}

	newValue, err := i.evalBinaryValues(value, stmt.Tok, &GooseValue{
		Constant: false,
		Type:     GooseTypeInt,
		Value:    int64(1),
	})

	if err != nil {
		return
	}

	if isIdent {
		scope.update(ident, *newValue)
	} else {
		existing.Value.([]*GooseValue)[index] = newValue
	}

	return &VoidResult{}, nil
}
