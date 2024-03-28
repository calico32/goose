package interpreter

import (
	"math/big"

	"github.com/calico32/goose/ast"
	. "github.com/calico32/goose/interpreter/lib"
	"github.com/calico32/goose/token"
)

func (i *interp) evalRangeExpr(scope *Scope, expr *ast.RangeExpr) Value {
	defer un(trace(i, "range expr"))

	start := i.evalExpr(scope, expr.Start)
	stop := i.evalExpr(scope, expr.Stop)
	var step Value = nil
	if expr.Step != nil {
		step = i.evalExpr(scope, expr.Step)
	}

	if _, ok := start.(Numeric); !ok {
		i.Throw("range start must be numeric")
	} else if _, ok := stop.(Numeric); !ok {
		i.Throw("range stop must be numeric")
	} else if step != nil {
		if _, ok := step.(Numeric); !ok {
			i.Throw("range step must be numeric")
		}
	}

	_, f1 := start.(*Float)
	_, f2 := stop.(*Float)
	_, f3 := step.(*Float)
	if f1 || f2 || f3 {
		s := 1.0
		if step != nil {
			s = step.(Numeric).Float64()
		}
		return &FloatRange{
			Start: start.(Numeric).Float64(),
			Stop:  stop.(Numeric).Float64(),
			Step:  s,
		}
	} else {
		s := big.NewInt(1)
		if step != nil {
			s = step.(Numeric).BigInt()
		}
		return &IntRange{
			Start: start.(Numeric).BigInt(),
			Stop:  stop.(Numeric).BigInt(),
			Step:  s,
		}
	}

}

func (i *interp) runRepeatCountStmt(scope *Scope, stmt *ast.RepeatCountStmt) StmtResult {
	defer un(trace(i, "repeat count stmt"))
	count := int64(0)
	totalCount := i.evalExpr(scope, stmt.Count)

	if _, ok := totalCount.(Numeric); !ok {
		i.Throw("repeat count must be numeric")
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

func (i *interp) runBranchStmt(_ *Scope, stmt *ast.BranchStmt) StmtResult {
	defer un(trace(i, "branch stmt"))
	switch stmt.Tok {
	case token.Break:
		return &Break{}
	case token.Continue:
		return &Continue{}
	default:
		i.Throw("unexpected branch type %v", stmt.Tok)
	}

	return nil
}

func (i *interp) runRepeatWhileStmt(scope *Scope, stmt *ast.RepeatWhileStmt) StmtResult {
	defer un(trace(i, "repeat while stmt"))
	for {
		cond := i.evalExpr(scope, stmt.Cond)
		if !IsTruthy(cond) {
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

	iterable := interp.spawnIterator(interp.evalExpr(scope, stmt.Iterable))

	name := stmt.Var.Name
	for iterVal := range iterable.channel {
		forScope := scope.Fork(ScopeOwnerFor)
		forScope.Set(name, &Variable{
			Constant: false,
			Value:    iterVal,
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

type iterator struct {
	channel chan Value
}

func (interp *interp) spawnIterator(iterable Value) *iterator {
	ch := make(chan Value)
	go func() {
		switch iterable := iterable.(type) {
		case *String:
			for _, char := range iterable.Value {
				ch <- Wrap(string(char))
			}
		case *Array:
			for _, elem := range iterable.Elements {
				ch <- elem
			}
		case *IntRange:
			// optimized for int64 ranges
			if iterable.Start.IsInt64() && iterable.Stop.IsInt64() && iterable.Step.IsInt64() {
				start := iterable.Start.Int64()
				stop := iterable.Stop.Int64()
				step := iterable.Step.Int64()

				if start > stop {
					if step > 0 {
						step = -step
					}
					for i := start; i > stop; i += step {
						ch <- Wrap(i)
					}
				} else {
					for i := start; i < stop; i += step {
						ch <- Wrap(i)
					}
				}
			} else {
				for i := iterable.Start; i.Cmp(iterable.Stop) == -1; i.Add(i, iterable.Step) {
					ch <- Wrap(i)
				}
			}
		case *FloatRange:
			for i := iterable.Start; i < iterable.Stop; i += iterable.Step {
				ch <- Wrap(i)
			}
		default:
			interp.Throw("for loop iterable must be... iterable")
		}
		close(ch)
	}()

	return &iterator{channel: ch}
}
