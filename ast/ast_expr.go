package ast

import (
	"github.com/calico32/goose/token"
)

type (
	ExprStmt struct {
		X Expr
	}

	BadExpr struct {
		From, To token.Pos
	}

	FrozenExpr struct {
		Frozen token.Pos
		X      Expr
	}

	Ident struct {
		NamePos token.Pos
		Name    string
	}

	SymbolExpr struct {
		Module *Ident
		Symbol *Ident
	}

	ParenExpr struct {
		Lparen token.Pos
		X      Expr
		Rparen token.Pos
	}

	SelectorExpr struct {
		X   Expr
		Sel *Ident
	}

	BracketSelectorExpr struct {
		X      Expr
		LBrack token.Pos
		Sel    Expr
		RBrack token.Pos
	}

	BindExpr struct {
		X   Expr
		Sel Expr
	}

	CallExpr struct {
		Func   Expr
		LParen token.Pos
		Args   []Expr
		RParen token.Pos
	}

	UnaryExpr struct {
		OpPos token.Pos
		Op    token.Token
		X     Expr
	}

	BinaryExpr struct {
		X     Expr
		OpPos token.Pos
		Op    token.Token
		Y     Expr
	}

	EllipsisExpr struct {
		Ellipsis token.Pos
		X        Expr
	}

	KeyValueExpr struct {
		Key   Expr
		Colon token.Pos
		Value Expr
	}

	AwaitExpr struct {
		Await token.Pos
		X     Expr
	}
)

func (s *ExprStmt) Pos() token.Pos            { return s.X.Pos() }
func (x *BadExpr) Pos() token.Pos             { return x.From }
func (x *FrozenExpr) Pos() token.Pos          { return x.Frozen }
func (x *Ident) Pos() token.Pos               { return x.NamePos }
func (x *ParenExpr) Pos() token.Pos           { return x.Lparen }
func (x *SelectorExpr) Pos() token.Pos        { return x.X.Pos() }
func (x *BracketSelectorExpr) Pos() token.Pos { return x.X.Pos() }
func (x *BindExpr) Pos() token.Pos            { return x.X.Pos() }
func (x *CallExpr) Pos() token.Pos            { return x.Func.Pos() }
func (x *UnaryExpr) Pos() token.Pos           { return x.OpPos }
func (x *BinaryExpr) Pos() token.Pos          { return x.X.Pos() }
func (x *EllipsisExpr) Pos() token.Pos        { return x.Ellipsis }
func (x *KeyValueExpr) Pos() token.Pos        { return x.Key.Pos() }
func (x *AwaitExpr) Pos() token.Pos           { return x.Await }

func (s *ExprStmt) End() token.Pos            { return s.X.End() }
func (x *BadExpr) End() token.Pos             { return x.To }
func (x *FrozenExpr) End() token.Pos          { return x.X.End() }
func (x *Ident) End() token.Pos               { return token.Pos(int(x.NamePos) + len(x.Name)) }
func (x *ParenExpr) End() token.Pos           { return x.Rparen + 1 }
func (x *SelectorExpr) End() token.Pos        { return x.Sel.End() }
func (x *BracketSelectorExpr) End() token.Pos { return x.RBrack + 1 }
func (x *BindExpr) End() token.Pos            { return x.Sel.End() }
func (x *CallExpr) End() token.Pos            { return x.RParen + 1 }
func (x *UnaryExpr) End() token.Pos           { return x.X.End() }
func (x *BinaryExpr) End() token.Pos          { return x.Y.End() }
func (x *EllipsisExpr) End() token.Pos        { return x.X.End() }
func (x *KeyValueExpr) End() token.Pos        { return x.Value.End() }
func (x *AwaitExpr) End() token.Pos           { return x.X.End() }

func (*ExprStmt) stmtNode()            {}
func (*BadExpr) exprNode()             {}
func (*FrozenExpr) exprNode()          {}
func (*Ident) exprNode()               {}
func (*ParenExpr) exprNode()           {}
func (*SelectorExpr) exprNode()        {}
func (*BracketSelectorExpr) exprNode() {}
func (*BindExpr) exprNode()            {}
func (*CallExpr) exprNode()            {}
func (*UnaryExpr) exprNode()           {}
func (*BinaryExpr) exprNode()          {}
func (*EllipsisExpr) exprNode()        {}
func (*KeyValueExpr) exprNode()        {}
func (*AwaitExpr) exprNode()           {}

func (id *Ident) String() string {
	if id != nil {
		return id.Name
	}
	return "<nil>"
}

func (s *ExprStmt) Flatten() []Node   { return s.X.Flatten() }
func (x *BadExpr) Flatten() []Node    { return nil }
func (x *FrozenExpr) Flatten() []Node { return x.X.Flatten() }
func (x *Ident) Flatten() []Node      { return []Node{x} }
func (x *ParenExpr) Flatten() []Node  { return x.X.Flatten() }
func (x *SelectorExpr) Flatten() []Node {
	nodes := make([]Node, 0, 2)
	nodes = append(nodes, x.X.Flatten()...)
	// nodes = append(nodes, x.Sel)
	return nodes
}
func (x *BracketSelectorExpr) Flatten() []Node {
	nodes := make([]Node, 0, 2)
	nodes = append(nodes, x.X.Flatten()...)
	nodes = append(nodes, x.Sel.Flatten()...)
	return nodes
}
func (x *BindExpr) Flatten() []Node {
	nodes := make([]Node, 0, 2)
	nodes = append(nodes, x.X.Flatten()...)
	nodes = append(nodes, x.Sel.Flatten()...)
	return nodes
}
func (x *CallExpr) Flatten() []Node {
	nodes := make([]Node, 0, len(x.Args)+1)
	nodes = append(nodes, x.Func.Flatten()...)
	for _, arg := range x.Args {
		nodes = append(nodes, arg.Flatten()...)
	}
	return nodes
}
func (x *UnaryExpr) Flatten() []Node    { return x.X.Flatten() }
func (x *BinaryExpr) Flatten() []Node   { return append(x.X.Flatten(), x.Y.Flatten()...) }
func (x *EllipsisExpr) Flatten() []Node { return x.X.Flatten() }
func (x *KeyValueExpr) Flatten() []Node {
	nodes := make([]Node, 0, 2)
	nodes = append(nodes, x.Key.Flatten()...)
	nodes = append(nodes, x.Value.Flatten()...)
	return nodes
}
func (x *AwaitExpr) Flatten() []Node { return x.X.Flatten() }
