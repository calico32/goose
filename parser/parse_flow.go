package parser

import (
	"github.com/calico32/goose/ast"
	"github.com/calico32/goose/token"
)

func (p *Parser) parseIfExprStmt(allowExpr bool) ast.Node {
	if p.trace {
		defer un(trace(p, "IfExprOrStmt"))
	}

	pos := p.expect(token.If)
	cond := p.parseExpr()

	if p.tok == token.Then {
		if !allowExpr {
			p.errorExpected(p.pos, "statement")
		}
		expr := &ast.IfExpr{If: pos, Cond: cond, ThenPos: p.expect(token.Then)}
		expr.Then = p.parseExpr()
		if p.tok == token.Else {
			expr.ElsePos = p.pos
			p.next()
			expr.Else = p.parseExpr()
		}
		if !allowExpr {
			return &ast.BadStmt{From: expr.If, To: expr.Else.End()}
		}
		return expr
	}

	var body []ast.Stmt
	for p.tok != token.End && p.tok != token.EOF && p.tok != token.Else {
		body = append(body, p.parseStmt())
	}

	var else_ []ast.Stmt
	if p.tok == token.Else {
		p.next()
		if p.tok == token.If {
			else_ = []ast.Stmt{p.parseIfExprStmt(false).(*ast.IfStmt)}
		} else {
			for p.tok != token.End && p.tok != token.EOF {
				else_ = append(else_, p.parseStmt())
			}
			p.expect(token.End)
		}
	} else {
		p.expect(token.End)
	}

	return &ast.IfStmt{If: pos, Cond: cond, Body: body, Else: else_, BlockEnd: p.pos}
}

func (p *Parser) parseTryStmt() *ast.TryStmt {
	if p.trace {
		defer un(trace(p, "TryStmt"))
	}

	pos := p.expect(token.Try)
	var body []ast.Stmt
	for p.tok != token.End && p.tok != token.EOF && p.tok != token.Catch && p.tok != token.Finally {
		body = append(body, p.parseStmt())
	}

	try := &ast.TryStmt{Try: pos, Body: body}

	if p.tok == token.Catch {

		try.BlockEnd = p.pos
		try.Catch = p.parseCatchStmt()
	}
	if p.tok == token.Finally {
		if !try.BlockEnd.IsValid() {
			try.BlockEnd = p.pos
		}
		try.Finally = p.parseFinallyStmt()
	}

	if try.Catch == nil && try.Finally == nil {
		p.errorExpected(p.pos, "catch or finally")
		p.next()
	} else {
		p.expect(token.End)
	}

	return try
}

func (p *Parser) parseCatchStmt() *ast.CatchStmt {
	if p.trace {
		defer un(trace(p, "CatchStmt"))
	}

	catch := &ast.CatchStmt{Catch: p.expect(token.Catch)}
	if p.tok == token.As {
		p.next()
		catch.Ident = p.parseIdent()
	}
	for p.tok != token.End && p.tok != token.EOF && p.tok != token.Finally {
		catch.Body = append(catch.Body, p.parseStmt())
	}

	catch.BlockEnd = p.pos

	return catch
}

func (p *Parser) parseFinallyStmt() *ast.FinallyStmt {
	if p.trace {
		defer un(trace(p, "FinallyStmt"))
	}

	pos := p.expect(token.Finally)
	var body []ast.Stmt
	for p.tok != token.End && p.tok != token.EOF {
		body = append(body, p.parseStmt())
	}

	finally := &ast.FinallyStmt{Finally: pos, Body: body, BlockEnd: p.pos}

	return finally
}

func (p *Parser) parseThrowExpr() ast.Expr {
	if p.trace {
		defer un(trace(p, "ThrowExpr"))
	}

	pos := p.expect(token.Throw)
	expr := p.parseExpr()

	return &ast.ThrowExpr{Throw: pos, X: expr}
}

func (p *Parser) parseDoExpr() ast.Expr {
	if p.trace {
		defer un(trace(p, "DoExpr"))
	}

	expr := &ast.DoExpr{
		Do:   p.expect(token.Do),
		Body: []ast.Stmt{},
	}

	for p.tok != token.End && p.tok != token.EOF {
		expr.Body = append(expr.Body, p.parseStmt())
	}

	expr.BlockEnd = p.expect(token.End)

	return expr
}
