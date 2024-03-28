package interpreter

import (
	"fmt"

	"github.com/calico32/goose/ast"
	"github.com/calico32/goose/token"
)

func (interp *interp) printTrace(a ...any) {
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

func push(i *interp, node ast.Node) *interp {
	i.lastPos = node.Pos()
	i.posStack = append(i.posStack, node.Pos())
	return i
}

func pop(i *interp) {
	if len(i.posStack) == 0 {
		i.Throw("stack underflow")
	}

	i.posStack = i.posStack[:len(i.posStack)-1]
}

func trace(i *interp, msg string) *interp {
	if i.trace {
		i.printTrace(msg, "(")
		i.indent++
	}
	return i
}

// Usage pattern: defer un(trace(p, "..."))
func un(i *interp) {
	if i.trace {
		i.indent--
		i.printTrace(")")
	}
}

func (i *interp) currentPos() token.Pos {
	if len(i.posStack) == 0 {
		return i.lastPos
	}
	return i.posStack[len(i.posStack)-1]
}

func (i *interp) Throw(msg string, parts ...any) {
	panic(fmt.Errorf("%s: Goose error: %s", i.fset.Position(i.currentPos()), fmt.Sprintf(msg, parts...)))
}
