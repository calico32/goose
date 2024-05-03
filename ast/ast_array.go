package ast

import "github.com/calico32/goose/token"

type (
	ArrayLiteral struct {
		Opening token.Pos
		List    []Expr
		Closing token.Pos
	}

	ArrayInitializer struct {
		Opening token.Pos
		Value   Expr
		Semi    token.Pos
		Count   Expr
		Closing token.Pos
	}

	SliceExpr struct {
		X      Expr
		LBrack token.Pos
		Low    Expr
		High   Expr
		RBrack token.Pos
	}
)

func (x *ArrayLiteral) Pos() token.Pos     { return x.Opening }
func (x *ArrayInitializer) Pos() token.Pos { return x.Opening }
func (x *SliceExpr) Pos() token.Pos        { return x.X.Pos() }

func (x *ArrayLiteral) End() token.Pos     { return x.Closing + 1 }
func (x *ArrayInitializer) End() token.Pos { return x.Closing + 1 }
func (x *SliceExpr) End() token.Pos        { return x.RBrack + 1 }

func (*ArrayLiteral) exprNode()     {}
func (*ArrayInitializer) exprNode() {}
func (*SliceExpr) exprNode()        {}

func (x *ArrayLiteral) Flatten() []Node {
	nodes := make([]Node, 0, len(x.List))
	for _, expr := range x.List {
		nodes = append(nodes, expr.Flatten()...)
	}
	return nodes
}

func (x *ArrayInitializer) Flatten() []Node {
	nodes := make([]Node, 0, 3)
	nodes = append(nodes, x.Value.Flatten()...)
	nodes = append(nodes, x.Count.Flatten()...)
	return nodes
}

func (x *SliceExpr) Flatten() []Node {
	nodes := make([]Node, 0, 3)
	nodes = append(nodes, x.X)
	nodes = append(nodes, x.Low)
	nodes = append(nodes, x.High)
	return nodes
}
