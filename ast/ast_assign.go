package ast

import "github.com/calico32/goose/token"

type (
	ConstStmt struct {
		ConstPos token.Pos
		Ident    *Ident
		TokPos   token.Pos
		Value    Expr
	}

	LetStmt struct {
		LetPos token.Pos
		Ident  *Ident
		TokPos token.Pos
		Value  Expr
	}

	AssignStmt struct {
		Lhs    Expr
		TokPos token.Pos
		Tok    token.Token
		Rhs    Expr
	}

	IncDecStmt struct {
		X      Expr
		TokPos token.Pos
		Tok    token.Token
	}
)

func (s *ConstStmt) End() token.Pos { return s.Value.End() }
func (s *LetStmt) End() token.Pos {
	if s.Value != nil {
		return s.Value.End()
	}
	return s.Ident.End()
}
func (s *AssignStmt) End() token.Pos { return s.Rhs.End() }
func (s *IncDecStmt) End() token.Pos { return s.TokPos + 2 /* len("++") */ }

func (s *ConstStmt) Pos() token.Pos  { return s.ConstPos }
func (s *LetStmt) Pos() token.Pos    { return s.LetPos }
func (s *AssignStmt) Pos() token.Pos { return s.Lhs.Pos() }
func (s *IncDecStmt) Pos() token.Pos { return s.X.Pos() }

func (*ConstStmt) stmtNode()  {}
func (*LetStmt) stmtNode()    {}
func (*AssignStmt) stmtNode() {}
func (*IncDecStmt) stmtNode() {}

func (s *ConstStmt) Flatten() []Node {
	nodes := make([]Node, 0, 2)
	nodes = append(nodes, s.Ident.Flatten()...)
	nodes = append(nodes, s.Value.Flatten()...)
	return nodes
}
func (s *LetStmt) Flatten() []Node {
	nodes := make([]Node, 0, 2)
	nodes = append(nodes, s.Ident.Flatten()...)
	if s.Value != nil {
		nodes = append(nodes, s.Value.Flatten()...)
	}
	return nodes
}
func (s *AssignStmt) Flatten() []Node {
	nodes := make([]Node, 0, 2)
	nodes = append(nodes, s.Lhs.Flatten()...)
	nodes = append(nodes, s.Rhs.Flatten()...)
	return nodes
}
func (s *IncDecStmt) Flatten() []Node {
	return s.X.Flatten()
}
