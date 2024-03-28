package parser

import (
	"github.com/calico32/goose/ast"
	"github.com/calico32/goose/token"
)

func (p *Parser) ParseExpr() ast.Expr {
	if p.trace {
		defer un(trace(p, "Expression"))
	}

	return p.parseBinaryExpr(nil, token.PrecedenceLowest+1)
}

func (p *Parser) parseBinaryExpr(x ast.Expr, precedence int) ast.Expr {
	if p.trace {
		defer un(trace(p, "BinaryExpr"))
	}

	if x == nil {
		x = p.parseUnaryExpr()
	}

	for {
		operator := p.tok
		operatorPrec := p.tok.Precedence()
		if operatorPrec < precedence {
			return x
		}
		pos := p.expect(operator)
		y := p.parseBinaryExpr(nil, operatorPrec+1)

		x = &ast.BinaryExpr{X: x, OpPos: pos, Op: operator, Y: y}
	}

}

func (p *Parser) parseUnaryExpr() ast.Expr {
	switch p.tok {
	case token.Add, token.Sub, token.LogNot, token.Ellipsis, token.Await:
		pos := p.pos
		op := p.tok
		p.next()
		x := p.parseUnaryExpr()
		return &ast.UnaryExpr{OpPos: pos, Op: op, X: x}
	}

	return p.parsePrimaryExpr(nil)
}

func (p *Parser) parsePrimaryExpr(x ast.Expr) ast.Expr {
	if p.trace {
		defer un(trace(p, "PrimaryExpr"))
	}

	if x == nil {
		x = p.parseOperand()
	}

	for {
		switch p.tok {
		case token.Bind:
			x = p.parseBindExpr(x)
		case token.LParen:
			x = p.parseCall(x)
		case token.LBracket:
			x = p.parseIndexOrSlice(x)
		case token.Period:
			x = p.parseSelectorExpr(x)
		case token.To:
			x = p.parseRangeExpr(x)
		case token.Question:
			x = p.parseQuestionExpr(x)
		default:
			return x
		}
	}
}

func (p *Parser) parseOperand() (e ast.Expr) {
	if p.trace {
		defer un(trace(p, "Operand"))
	}

	switch p.tok {
	case token.Ident:
		e = p.parseIdent()
	case token.Int, token.Float, token.Null, token.True, token.False:
		e = &ast.Literal{Value: p.lit, ValuePos: p.pos, Kind: p.tok}
		p.next()
	case token.StringStart:
		e = p.parseString()
	case token.LBracket:
		e = p.parseArrayLitOrInitializer()
	case token.HashLBracket:
		e = p.parseBracketPropertyExpr()
	case token.LParen:
		lparen := p.expect(token.LParen)
		x := p.ParseExpr()
		rparen := p.expect(token.RParen)
		e = &ast.ParenExpr{Lparen: lparen, X: x, Rparen: rparen}
	case token.LBrace:
		e = p.parseComposite()
	case token.Func, token.Memo:
		e = p.parseFuncExpr(0)
	case token.Generator:
		e = p.parseGeneratorExpr(0)
	case token.Throw:
		e = p.parseThrowExpr()
	case token.Frozen:
		pos := p.pos
		p.next()
		x := p.ParseExpr()
		e = &ast.FrozenExpr{Frozen: pos, X: x}
	case token.Native:
		e = p.parseNativeExpr()
	case token.Async:
		pos := p.pos
		p.next()
		switch p.tok {
		case token.Func:
			e = p.parseFuncExpr(pos)
		case token.Generator:
			e = p.parseGeneratorExpr(pos)
		default:
			p.errorExpected(pos, "func or generator")
			p.next() // make progress
			e = &ast.BadExpr{From: pos, To: p.safePos(p.pos)}
		}
	case token.If:
		node := p.parseIfExprStmt(true)
		if stmt, isStmt := node.(*ast.IfStmt); isStmt {
			p.errorExpected(stmt.If, "expression")
			e = &ast.BadExpr{From: stmt.If, To: p.safePos(stmt.BlockEnd)}
		} else {
			e = node.(*ast.IfExpr)
		}
	case token.Do:
		e = p.parseDoExpr()
	case token.Match:
		e = p.parseMatchExpr()
	default:
		pos := p.pos
		p.errorExpected(pos, "operand")
		p.next() // make progress
		e = &ast.BadExpr{From: pos, To: p.safePos(p.pos)}
	}

	return
}

func (p *Parser) parseQuestionExpr(c ast.Expr) *ast.UnaryExpr {
	if p.trace {
		defer un(trace(p, "QuestionExpr"))
	}

	qmark := p.expect(token.Question)

	return &ast.UnaryExpr{
		X:     c,
		OpPos: qmark,
		Op:    token.Question,
	}
}
