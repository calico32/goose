package parser

import (
	"github.com/calico32/goose/ast"
	"github.com/calico32/goose/token"
)

func (p *Parser) parseMatchExpr() *ast.MatchExpr {
	defer un(trace(p, "MatchExpr"))

	match := p.expect(token.Match)
	expr := p.ParseExpr()
	clauses := []ast.MatchArm{}
	for p.tok != token.End && p.tok != token.EOF {
		clauses = append(clauses, p.parseMatchArm())
	}
	end := p.expect(token.End)
	return &ast.MatchExpr{
		Match:    match,
		Expr:     expr,
		Clauses:  clauses,
		BlockEnd: end,
	}
}

func (p *Parser) parseMatchArm() ast.MatchArm {
	defer un(trace(p, "MatchArm"))

	if p.tok == token.Else {
		elseTok := p.expect(token.Else)
		arrow := p.expect(token.Arrow)
		expr := p.ParseExpr()
		return &ast.MatchElse{
			Else:  elseTok,
			Arrow: arrow,
			Expr:  expr,
		}
	}
	pattern := p.parsePatternExpr()
	arrow := p.expect(token.Arrow)
	expr := p.ParseExpr()
	return &ast.MatchPattern{
		Pattern: pattern,
		Arrow:   arrow,
		Expr:    expr,
	}
}

func (p *Parser) parsePatternExpr() ast.PatternExpr {
	defer un(trace(p, "PatternExpr"))

	switch p.tok {
	case token.LParen:
		lparen := p.expect(token.LParen)
		pattern := p.parsePatternExpr()
		rparen := p.expect(token.RParen)
		return &ast.PatternParen{X: pattern, LParen: lparen, RParen: rparen}
	case token.LBracket:
		return p.parsePatternTuple()
	case token.LBrace:
		return p.parsePatternComposite()
	case token.MatchBind:
		bind := p.expect(token.MatchBind)
		ident := p.parseIdent()
		return &ast.PatternBinding{Bind: bind, Ident: ident}
	default:
		return &ast.PatternNormal{X: p.parseBinaryExpr(nil, token.Arrow.Precedence()+1)}
	}
}

func (p *Parser) parsePatternTuple() *ast.PatternTuple {
	defer un(trace(p, "PatternTuple"))

	lbracket := p.expect(token.LBracket)
	patterns := make([]ast.PatternExpr, 0)
	for p.tok != token.RBracket && p.tok != token.EOF {
		patterns = append(patterns, p.parsePatternExpr())
		if p.tok != token.RBracket {
			p.expect(token.Comma)
		}
	}
	rbracket := p.expect(token.RBracket)
	return &ast.PatternTuple{List: patterns, Opening: lbracket, Closing: rbracket}
}

func (p *Parser) parsePatternComposite() (s *ast.PatternComposite) {
	defer un(trace(p, "PatternComposite"))

	if p.trace {
		defer un(trace(p, "CompositeLiteral"))
	}

	lbrace := p.expect(token.LBrace)

	var list []*ast.PatternCompositeField

	for p.tok != token.RBrace && p.tok != token.EOF {
		list = append(list, p.parsePatternCompositeField())
		if p.tok != token.RBrace {
			p.expect(token.Comma)
		}
	}

	rbrace := p.expect(token.RBrace)

	return &ast.PatternComposite{
		Opening: lbrace,
		Fields:  list,
		Closing: rbrace,
	}
}

func (p *Parser) parsePatternCompositeField() (s *ast.PatternCompositeField) {
	defer un(trace(p, "PatternCompositeField"))

	if p.trace {
		defer un(trace(p, "CompositeField"))
	}

	var key ast.PatternExpr
	var value ast.PatternExpr

	if p.tok == token.LBracket {
		p.next()
		key = p.parsePatternExpr()
		p.expect(token.RBracket)
	} else if p.tok == token.StringStart {
		key = &ast.PatternNormal{X: p.parseString()}
	} else if p.tok == token.Int || p.tok == token.Float {
		key = &ast.PatternNormal{X: &ast.Literal{Value: p.lit, ValuePos: p.pos, Kind: p.tok}}
		p.next()
	} else {
		key = &ast.PatternNormal{X: p.parseIdent()}
	}

	p.expect(token.Colon)
	value = p.parsePatternExpr()

	return &ast.PatternCompositeField{
		Key:   key,
		Value: value,
	}
}
