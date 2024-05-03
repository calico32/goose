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
	Flatten() []Node
}

type PosRange struct {
	From token.Pos
	To   token.Pos
}

func (p PosRange) Pos() token.Pos { return p.From }
func (p PosRange) End() token.Pos { return p.To }

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

func (s *BadStmt) Flatten() []Node     { return []Node{s} }
func (s *EmptyStmt) Flatten() []Node   { return []Node{s} }
func (s *LabeledStmt) Flatten() []Node { return []Node{s} }

type Module struct {
	Type      ModuleType
	Specifier string
	Scheme    string

	Size  int
	Stmts []Stmt
	Nodes []Node
}

func (m *Module) Pos() token.Pos {
	if len(m.Stmts) > 0 {
		return m.Stmts[0].Pos()
	}
	return token.NoPos
}

func (m *Module) End() token.Pos {
	if len(m.Stmts) > 0 {
		return m.Stmts[len(m.Stmts)-1].End()
	}
	return token.NoPos
}

func (m *Module) Flatten() []Node {
	if m.Nodes != nil {
		return m.Nodes
	}

	nodes := make([]Node, 0, len(m.Stmts))
	for _, stmt := range m.Stmts {
		nodes = append(nodes, stmt.Flatten()...)
	}
	m.Nodes = nodes
	return nodes
}

func (m *Module) FindNode(pos token.Pos) Node {
	return nil
	nodes := m.Flatten()
	i := 0
	j := len(nodes)
	for i < j {
		h := i + (j-i)/2
		if pos < nodes[h].Pos() {
			j = h
		} else {
			i = h + 1
		}
	}
	if i < len(nodes) && pos >= nodes[i].Pos() && pos < nodes[i].End() {
		return nodes[i]
	}
	return nil
}

type ModuleType int

const (
	ModuleTypeFilesystem ModuleType = iota
	ModuleTypeMemory
	ModuleTypeNetwork
)
