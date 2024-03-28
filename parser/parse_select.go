package parser

import (
	"github.com/calico32/goose/ast"
	"github.com/calico32/goose/token"
)

func (p *Parser) parseBindExpr(c ast.Expr) ast.Expr {
	if p.trace {
		defer un(trace(p, "BindExpr"))
	}

	p.expect(token.Bind)

	var sel ast.Expr
	if p.tok == token.LParen {
		p.next()
		sel = p.ParseExpr()
		p.expect(token.RParen)
	} else {
		sel = p.parseOperand()
		// TODO: is this safe?
		for p.tok == token.Period {
			sel = p.parseSelectorExpr(sel)
		}
	}

	return &ast.BindExpr{
		X:   c,
		Sel: sel,
	}
}

func (p *Parser) parseSelectorExpr(c ast.Expr) *ast.SelectorExpr {
	if p.trace {
		defer un(trace(p, "SelectorExpr"))
	}

	p.expect(token.Period)

	return &ast.SelectorExpr{
		X:   c,
		Sel: p.parseIdent(),
	}
}

func (p *Parser) parseIndexOrSlice(x ast.Expr) ast.Expr {
	if p.trace {
		defer un(trace(p, "IndexOrSlice"))
	}

	lbrack := p.expect(token.LBracket)
	if p.tok == token.RBracket {
		// empty index
		p.errorExpected(p.pos, "expression")
		rbrack := p.pos
		p.next()
		return &ast.BracketSelectorExpr{
			X:      x,
			LBrack: lbrack,
			Sel:    &ast.BadExpr{From: lbrack, To: rbrack},
			RBrack: rbrack,
		}
	}

	var left ast.Expr
	var right ast.Expr
	isSlicing := false

	if p.tok != token.Colon {
		left = p.ParseExpr()
		if p.tok == token.Colon {
			isSlicing = true
			p.next()
		}
	} else {
		isSlicing = true
	}

	if p.tok != token.RBracket {
		right = p.ParseExpr()
	}

	if left == nil && right == nil {
		p.errorExpected(p.pos, "slice expression")
		return &ast.SliceExpr{
			X:      x,
			LBrack: lbrack,
			Low:    &ast.BadExpr{From: lbrack, To: p.pos},
			High:   &ast.BadExpr{From: p.pos, To: p.pos},
			RBrack: p.safePos(p.pos),
		}
	}

	rbrack := p.expect(token.RBracket)

	if isSlicing {
		return &ast.SliceExpr{
			X:      x,
			LBrack: lbrack,
			Low:    left,
			High:   right,
			RBrack: rbrack,
		}
	}

	return &ast.BracketSelectorExpr{
		X:      x,
		LBrack: lbrack,
		Sel:    left,
		RBrack: rbrack,
	}
}
