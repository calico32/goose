package ast

import (
	"github.com/calico32/goose/token"
)

func isWhitespace(ch byte) bool { return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' }

func trimRight(s string) string {
	i := len(s)
	for i > 0 && isWhitespace(s[i-1]) {
		i--
	}
	return s[0:i]
}

type Token struct {
	Type   token.Token
	Source string
}

type Node interface {
	Pos() token.Pos
	End() token.Pos
}

type Expr interface {
	Node
	exprNode()
}

type Stmt interface {
	Node
	stmtNode()
}

type (
	BadStmt struct {
		From, To token.Pos
	}

	EmptyStmt struct {
	}

	LabeledStmt struct {
		Label *Ident
		Colon token.Pos
		Stmt  Stmt
	}
)

func (s *BadStmt) Pos() token.Pos     { return s.From }
func (s *EmptyStmt) Pos() token.Pos   { return 0 }
func (s *LabeledStmt) Pos() token.Pos { return s.Label.Pos() }

func (s *BadStmt) End() token.Pos     { return s.To }
func (s *EmptyStmt) End() token.Pos   { return s.Pos() }
func (s *LabeledStmt) End() token.Pos { return s.Stmt.End() }

func (*BadStmt) stmtNode()     {}
func (*EmptyStmt) stmtNode()   {}
func (*LabeledStmt) stmtNode() {}

type File struct {
	Type ModuleType
	Name string

	Size  int
	Stmts []Stmt
}

type ModuleType int

const (
	Filesystem ModuleType = iota
	Memory
	Network
)
