package interpreter

import (
	"fmt"
	"io"
	"os"

	"github.com/wiisportsresort/goose/ast"
	"github.com/wiisportsresort/goose/token"
)

func RunFile(file *ast.File, fset *token.FileSet, trace bool, doPanic bool, stdout io.Writer, stderr io.Writer) (exitCode int, err error) {
	i := &interpreter{
		file:   file,
		trace:  trace,
		fset:   fset,
		stdout: stdout,
		stderr: stderr,
	}

	return i.run(doPanic), nil
}

type interpreter struct {
	fset   *token.FileSet
	file   *ast.File
	global *GooseScope
	stdout io.Writer
	stderr io.Writer

	// internal state
	trace    bool
	indent   int
	lastPos  token.Pos
	posStack []token.Pos
}

func (interp *interpreter) printTrace(a ...any) {
	const dots = ". . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . "
	const n = len(dots)
	i := 2 * interp.indent
	for i > n {
		fmt.Fprint(interp.stderr, dots)
		i -= n
	}
	// i <= n
	fmt.Fprint(interp.stderr, dots[0:i])
	for _, arg := range a {
		fmt.Fprintf(interp.stderr, "%v ", arg)
	}
}

func push(i *interpreter, node ast.Node) *interpreter {
	i.lastPos = node.Pos()
	i.posStack = append(i.posStack, node.Pos())
	return i
}

func pop(i *interpreter) {
	if len(i.posStack) == 0 {
		i.throw("stack underflow")
	}

	i.posStack = i.posStack[:len(i.posStack)-1]
}

func trace(i *interpreter, msg string) *interpreter {
	if i.trace {
		i.printTrace(msg, "(")
		i.indent++
	}
	return i
}

// Usage pattern: defer un(trace(p, "..."))
func un(i *interpreter) {
	if i.trace {
		i.indent--
		i.printTrace(")")
	}
}

func (i *interpreter) currentPos() token.Pos {
	if len(i.posStack) == 0 {
		return i.lastPos
	}
	return i.posStack[len(i.posStack)-1]
}

func (i *interpreter) throw(msg string) {
	panic(fmt.Errorf("%s: Goose error: %s", i.fset.Position(i.currentPos()), msg))
}

const PANIC int = 128

func (i *interpreter) run(doPanic bool) (exitCode int) {
	defer func() {
		if r := recover(); r != nil {
			if exit, ok := r.(gooseExit); ok {
				exitCode = int(exit.code)
			} else {
				if doPanic {
					panic(r)
				}
				fmt.Fprintf(i.stderr, "panic: %#v+\n", r)
				exitCode = PANIC
			}
		}
	}()

	if i.stdout == nil {
		i.stdout = os.Stdout
	}
	if i.stderr == nil {
		i.stderr = os.Stderr
	}
	i.global = NewGlobalScope(i)
	for k, v := range stdlib {
		i.global.set(k, GooseValue{
			Constant: false, // allow overwriting of builtins
			Type:     GooseTypeFunc,
			Value:    v,
		})
	}

	for _, stmt := range i.file.Stmts {
		result, err := i.runStmt(i.global, stmt)
		if err != nil {
			i.throw(err.Error())
		}

		switch result.(type) {
		case *ReturnResult:
			i.throw("cannot return from top-level")
		case *BreakResult:
			i.throw("cannot break from top-level")
		case *ContinueResult:
			i.throw("cannot continue from top-level")
		}
	}

	return 0
}

func expectType(value *GooseValue, t GooseType) error {
	if t > GooseTypeError {
		if value.Type&t != 0 {
			return nil
		} else {
			return fmt.Errorf("expected type %s, got %s", t, value.Type)
		}
	}

	if !isPowerOfTwo(int(t)) {
		return fmt.Errorf("unexpected GooseType %d", t)
	}

	if value.Type != t {
		return fmt.Errorf("expected type %s, got type %s", t, value.Type)
	}

	return nil
}

func (i *interpreter) numericOperation(lhs *GooseValue, op token.Token, rhs *GooseValue) (*GooseValue, error) {
	defer un(trace(i, "numeric op"))
	err := expectType(lhs, GooseTypeNumeric)
	if err != nil {
		return nil, err
	}
	err = expectType(rhs, GooseTypeNumeric)
	if err != nil {
		return nil, err
	}

	lt := typeOf(lhs)
	rt := typeOf(rhs)

	var result any
	if lt == GooseTypeFloat {
		if rt == GooseTypeFloat {
			result = numericOp(lhs.Value.(float64), op, rhs.Value.(float64))
		} else {
			result = numericOp(lhs.Value.(float64), op, float64(rhs.Value.(int64)))
		}
	} else {
		if rt == GooseTypeFloat {
			result = numericOp(float64(lhs.Value.(int64)), op, rhs.Value.(float64))
		} else {
			result = numericOp(lhs.Value.(int64), op, rhs.Value.(int64))
		}
	}

	return wrap(result), nil
}
