package ast

import "github.com/calico32/goose/token"

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
		Name      string
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
	return token.Pos(int(x.InterpPos)+len(x.Name)) + 1
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
