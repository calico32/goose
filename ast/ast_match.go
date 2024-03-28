package ast

import "github.com/calico32/goose/token"

type (
	MatchExpr struct {
		Match    token.Pos
		Expr     Expr
		Clauses  []MatchArm
		BlockEnd token.Pos
	}

	MatchArm interface {
		matchArm()
	}

	MatchElse struct {
		Else  token.Pos
		Arrow token.Pos
		Expr  Expr
	}

	MatchPattern struct {
		Pattern PatternExpr
		Arrow   token.Pos
		Expr    Expr
	}

	PatternExpr interface {
		Node
		patternExpr()
	}

	PatternNormal struct {
		X Expr
	}

	PatternBinding struct {
		Bind  token.Pos
		Ident *Ident
	}

	PatternParen struct {
		LParen token.Pos
		X      PatternExpr
		RParen token.Pos
	}

	PatternTuple struct {
		Opening token.Pos
		List    []PatternExpr
		Closing token.Pos
	}

	PatternRange struct {
		Start Expr
		To    token.Pos
		Stop  Expr
	}

	PatternType struct {
		Ident   *Ident
		Opening token.Pos
		Binding PatternBinding
		Closing token.Pos
	}

	PatternComposite struct {
		Opening token.Pos
		Fields  []*PatternCompositeField
		Closing token.Pos
	}

	PatternCompositeField struct {
		Key   PatternExpr
		Colon token.Pos
		Value PatternExpr
	}
)

func (x *MatchExpr) Pos() token.Pos        { return x.Match }
func (x *MatchElse) Pos() token.Pos        { return x.Else }
func (x *MatchPattern) Pos() token.Pos     { return x.Pattern.Pos() }
func (x *PatternNormal) Pos() token.Pos    { return x.X.Pos() }
func (x *PatternBinding) Pos() token.Pos   { return x.Bind }
func (x *PatternParen) Pos() token.Pos     { return x.LParen }
func (x *PatternTuple) Pos() token.Pos     { return x.Opening }
func (x *PatternRange) Pos() token.Pos     { return x.Start.Pos() }
func (x *PatternType) Pos() token.Pos      { return x.Ident.Pos() }
func (x *PatternComposite) Pos() token.Pos { return x.Opening }

func (x *MatchExpr) End() token.Pos        { return x.BlockEnd + 3 }
func (x *MatchElse) End() token.Pos        { return x.Expr.End() }
func (x *MatchPattern) End() token.Pos     { return x.Expr.End() }
func (x *PatternNormal) End() token.Pos    { return x.X.End() }
func (x *PatternBinding) End() token.Pos   { return x.Ident.End() }
func (x *PatternTuple) End() token.Pos     { return x.Closing + 1 }
func (x *PatternParen) End() token.Pos     { return x.RParen + 1 }
func (x *PatternRange) End() token.Pos     { return x.Stop.End() }
func (x *PatternType) End() token.Pos      { return x.Closing + 1 }
func (x *PatternComposite) End() token.Pos { return x.Closing + 1 }

func (x *MatchExpr) exprNode() {}

func (x *MatchElse) matchArm()    {}
func (x *MatchPattern) matchArm() {}

func (x *PatternNormal) patternExpr()    {}
func (x *PatternBinding) patternExpr()   {}
func (x *PatternParen) patternExpr()     {}
func (x *PatternTuple) patternExpr()     {}
func (x *PatternRange) patternExpr()     {}
func (x *PatternType) patternExpr()      {}
func (x *PatternComposite) patternExpr() {}
