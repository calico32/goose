package ast

import "github.com/calico32/goose/token"

type (
	MatchStmt struct {
		Match    token.Pos
		Expr     Expr
		Clauses  []*MatchArm
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
		Expr
		patternExpr()
	}

	PatternBinding struct {
		Dollar token.Pos
		Ident  *Ident
	}

	PatternTuple struct {
		Opening token.Pos
		List    []PatternExpr
		Closing token.Pos
	}

	PatternRange struct {
		Low  Expr
		To   token.Pos
		High Expr
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

	PatternCompositeField interface {
		patternCompositeField()
	}

	PatternCompositeFieldIdent struct {
		Ident *Ident
		Colon token.Pos
		Value PatternExpr
	}

	PatternCompositeFieldExpr struct {
		Expr  Expr
		Colon token.Pos
		Value PatternExpr
	}
)

func (x *MatchStmt) Pos() token.Pos        { return x.Match }
func (x *MatchElse) Pos() token.Pos        { return x.Else }
func (x *MatchPattern) Pos() token.Pos     { return x.Pattern.Pos() }
func (x *PatternBinding) Pos() token.Pos   { return x.Dollar }
func (x *PatternTuple) Pos() token.Pos     { return x.Opening }
func (x *PatternRange) Pos() token.Pos     { return x.Low.Pos() }
func (x *PatternType) Pos() token.Pos      { return x.Ident.Pos() }
func (x *PatternComposite) Pos() token.Pos { return x.Opening }

func (x *MatchElse) End() token.Pos        { return x.Expr.End() }
func (x *MatchPattern) End() token.Pos     { return x.Expr.End() }
func (x *PatternBinding) End() token.Pos   { return x.Ident.End() }
func (x *PatternTuple) End() token.Pos     { return x.Closing }
func (x *PatternRange) End() token.Pos     { return x.High.End() }
func (x *PatternType) End() token.Pos      { return x.Closing }
func (x *PatternComposite) End() token.Pos { return x.Closing }

func (x *MatchElse) matchArm()           {}
func (x *MatchPattern) matchArm()        {}
func (x *PatternBinding) patternExpr()   {}
func (x *PatternTuple) patternExpr()     {}
func (x *PatternRange) patternExpr()     {}
func (x *PatternType) patternExpr()      {}
func (x *PatternComposite) patternExpr() {}

func (X *PatternCompositeFieldIdent) patternCompositeField() {}
func (X *PatternCompositeFieldExpr) patternCompositeField()  {}

// extensions
func (X *Literal) patternExpr()      {}
func (X *Ident) patternExpr()        {}
func (X *ParenExpr) patternExpr()    {}
func (X *SelectorExpr) patternExpr() {}
