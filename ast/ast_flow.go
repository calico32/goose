package ast

import "github.com/calico32/goose/token"

type (
	IfStmt struct {
		If       token.Pos
		Cond     Expr
		Body     []Stmt
		Else     []Stmt
		BlockEnd token.Pos
	}

	IfExpr struct {
		If      token.Pos
		Cond    Expr
		ThenPos token.Pos
		Then    Expr
		ElsePos token.Pos
		Else    Expr
	}

	TryStmt struct {
		Try      token.Pos
		Body     []Stmt
		Catch    *CatchStmt
		Finally  *FinallyStmt
		BlockEnd token.Pos
	}

	CatchStmt struct {
		Catch    token.Pos
		Ident    *Ident
		Body     []Stmt
		BlockEnd token.Pos
	}

	FinallyStmt struct {
		Finally  token.Pos
		Body     []Stmt
		BlockEnd token.Pos
	}

	ThrowExpr struct {
		Throw token.Pos
		X     Expr
	}

	DoExpr struct {
		Do       token.Pos
		Body     []Stmt
		BlockEnd token.Pos
	}
)

func (s *IfStmt) Pos() token.Pos      { return s.If }
func (e *IfExpr) Pos() token.Pos      { return e.If }
func (s *TryStmt) Pos() token.Pos     { return s.Try }
func (s *CatchStmt) Pos() token.Pos   { return s.Catch }
func (s *FinallyStmt) Pos() token.Pos { return s.Finally }
func (s *ThrowExpr) Pos() token.Pos   { return s.Throw }
func (s *DoExpr) Pos() token.Pos      { return s.Do }

func (s *IfStmt) End() token.Pos {
	if s.Else != nil {
		return s.Else[len(s.Else)-1].End()
	}
	return s.BlockEnd
}
func (e *IfExpr) End() token.Pos {
	if e.Else != nil {
		return e.Else.End()
	}
	return e.Then.End()
}
func (s *TryStmt) End() token.Pos {
	if s.Finally != nil {
		return s.Finally.End()
	}
	if s.Catch != nil {
		return s.Catch.End()
	}
	return s.BlockEnd
}
func (s *CatchStmt) End() token.Pos {
	if s.BlockEnd.IsValid() {
		return s.BlockEnd
	}
	return s.Body[len(s.Body)-1].End()
}
func (s *FinallyStmt) End() token.Pos {
	if s.BlockEnd.IsValid() {
		return s.BlockEnd
	}
	return s.Body[len(s.Body)-1].End()
}
func (s *ThrowExpr) End() token.Pos {
	return s.X.End()
}
func (s *DoExpr) End() token.Pos { return s.BlockEnd }

func (*IfStmt) stmtNode()      {}
func (*IfExpr) exprNode()      {}
func (*TryStmt) stmtNode()     {}
func (*CatchStmt) stmtNode()   {}
func (*FinallyStmt) stmtNode() {}
func (*ThrowExpr) exprNode()   {}
func (*DoExpr) exprNode()      {}

func (s *IfStmt) Flatten() []Node {
	nodes := make([]Node, 0, len(s.Body))
	nodes = append(nodes, s.Cond.Flatten()...)
	for _, stmt := range s.Body {
		nodes = append(nodes, stmt.Flatten()...)
	}
	for _, stmt := range s.Else {
		nodes = append(nodes, stmt.Flatten()...)
	}
	return nodes
}

func (e *IfExpr) Flatten() []Node {
	nodes := make([]Node, 0, 3)
	nodes = append(nodes, e.Cond.Flatten()...)
	nodes = append(nodes, e.Then.Flatten()...)
	nodes = append(nodes, e.Else.Flatten()...)
	return nodes
}

func (s *TryStmt) Flatten() []Node {
	nodes := make([]Node, 0, len(s.Body))
	for _, stmt := range s.Body {
		nodes = append(nodes, stmt.Flatten()...)
	}
	if s.Catch != nil {
		nodes = append(nodes, s.Catch.Flatten()...)
	}
	if s.Finally != nil {
		nodes = append(nodes, s.Finally.Flatten()...)
	}
	return nodes
}

func (s *CatchStmt) Flatten() []Node {
	nodes := make([]Node, 0, len(s.Body))
	for _, stmt := range s.Body {
		nodes = append(nodes, stmt.Flatten()...)
	}
	return nodes
}

func (s *FinallyStmt) Flatten() []Node {
	nodes := make([]Node, 0, len(s.Body))
	for _, stmt := range s.Body {
		nodes = append(nodes, stmt.Flatten()...)
	}
	return nodes
}

func (s *ThrowExpr) Flatten() []Node {
	return s.X.Flatten()
}

func (s *DoExpr) Flatten() []Node {
	nodes := make([]Node, 0, len(s.Body))
	for _, stmt := range s.Body {
		nodes = append(nodes, stmt.Flatten()...)
	}
	return nodes
}
