package ast

import (
	"strings"

	"github.com/calico32/goose/token"
)

type (
	Literal struct {
		ValuePos token.Pos
		Kind     token.Token
		Value    string
	}

	StringLiteral struct {
		StringStart *StringLiteralStart
		Parts       []StringLiteralPart
		StringEnd   *StringLiteralEnd
	}

	StringLiteralStart struct {
		Quote   token.Pos
		Content string
	}

	StringLiteralMiddle struct {
		StartPos token.Pos
		Content  string
	}

	StringLiteralInterpIdent struct {
		InterpPos token.Pos
		Ident     *Ident
	}

	StringLiteralInterpExpr struct {
		InterpPos token.Pos
		Expr      Expr
	}

	StringLiteralEnd struct {
		StartPos token.Pos
		Content  string
		Quote    token.Pos
	}

	SymbolStmt struct {
		Symbol token.Pos
		Ident  *Ident
	}
)

type StringLiteralPart interface {
	Node
	stringLiteralPart()
}

func (x *Literal) Pos() token.Pos                  { return x.ValuePos }
func (x *StringLiteral) Pos() token.Pos            { return x.StringStart.Pos() }
func (x *StringLiteralStart) Pos() token.Pos       { return x.Quote }
func (x *StringLiteralMiddle) Pos() token.Pos      { return x.StartPos }
func (x *StringLiteralInterpIdent) Pos() token.Pos { return x.InterpPos }
func (x *StringLiteralInterpExpr) Pos() token.Pos  { return x.InterpPos }
func (x *StringLiteralEnd) Pos() token.Pos         { return x.StartPos }
func (x *SymbolStmt) Pos() token.Pos               { return x.Symbol }

func (x *Literal) End() token.Pos             { return token.Pos(int(x.ValuePos) + len(x.Value)) }
func (x *StringLiteral) End() token.Pos       { return x.StringEnd.Quote + 1 }
func (x *StringLiteralStart) End() token.Pos  { return token.Pos(int(x.Quote) + len(x.Content)) }
func (x *StringLiteralMiddle) End() token.Pos { return token.Pos(int(x.StartPos) + len(x.Content)) }
func (x *StringLiteralInterpIdent) End() token.Pos {
	return token.Pos(int(x.InterpPos)+len(x.Ident.Name)) + 1
}
func (x *StringLiteralInterpExpr) End() token.Pos { return x.Expr.End() }
func (x *StringLiteralEnd) End() token.Pos        { return x.Quote + 1 }
func (x *SymbolStmt) End() token.Pos              { return x.Ident.End() }

func (*Literal) exprNode()       {}
func (*StringLiteral) exprNode() {}

func (*StringLiteralMiddle) stringLiteralPart()      {}
func (*StringLiteralInterpIdent) stringLiteralPart() {}
func (*StringLiteralInterpExpr) stringLiteralPart()  {}

func (*SymbolStmt) stmtNode() {}

func (s *StringLiteral) String() string {
	var sb strings.Builder
	sb.WriteString(s.StringStart.Content)
	for _, part := range s.Parts {
		switch p := part.(type) {
		case *StringLiteralMiddle:
			sb.WriteString(p.Content)
		case *StringLiteralInterpIdent:
			sb.WriteString("$")
			sb.WriteString(p.Ident.Name)
		case *StringLiteralInterpExpr:
			sb.WriteString("${expr}")
		}
	}
	sb.WriteString(s.StringEnd.Content)
	return sb.String()
}

func (s *Literal) Flatten() []Node { return nil }
func (s *StringLiteral) Flatten() []Node {
	nodes := make([]Node, 0, len(s.Parts))
	for _, part := range s.Parts {
		nodes = append(nodes, part.Flatten()...)
	}
	return nodes
}
func (s *StringLiteralStart) Flatten() []Node       { return nil }
func (s *StringLiteralMiddle) Flatten() []Node      { return nil }
func (s *StringLiteralInterpIdent) Flatten() []Node { return s.Ident.Flatten() }
func (s *StringLiteralInterpExpr) Flatten() []Node  { return s.Expr.Flatten() }
func (s *StringLiteralEnd) Flatten() []Node         { return nil }
func (s *SymbolStmt) Flatten() []Node               { return nil }
