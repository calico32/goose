package parser

import (
	"github.com/calico32/goose/ast"
	"github.com/calico32/goose/token"
)

func (p *Parser) parseForStmt() *ast.ForStmt {
	if p.trace {
		defer un(trace(p, "ForStmt"))
	}

	pos := p.expect(token.For)
	var awaitPos token.Pos
	if p.tok == token.Await {
		awaitPos = p.pos
		p.next()
	}
	ident := p.parseIdent()
	p.expect(token.In)
	expr := p.ParseExpr()

	var body []ast.Stmt
	for p.tok != token.End && p.tok != token.EOF {
		body = append(body, p.parseStmt())
	}
	p.expect(token.End)

	return &ast.ForStmt{
		For:      pos,
		Await:    awaitPos,
		Var:      ident,
		Iterable: expr,
		Body:     body,
		BlockEnd: p.pos,
	}
}

func (p *Parser) parseRepeatStmt() ast.Stmt {
	if p.trace {
		defer un(trace(p, "RepeatStmt"))
	}

	pos := p.expect(token.Repeat)
	switch p.tok {
	case token.While:
		p.next()
		cond := p.ParseExpr()
		var body []ast.Stmt
		for p.tok != token.End && p.tok != token.EOF {
			body = append(body, p.parseStmt())
		}
		p.expect(token.End)
		return &ast.RepeatWhileStmt{Repeat: pos, Cond: cond, Body: body, BlockEnd: p.pos}
	case token.Forever:
		p.next()
		var body []ast.Stmt
		for p.tok != token.End && p.tok != token.EOF {
			body = append(body, p.parseStmt())
		}
		p.expect(token.End)
		return &ast.RepeatForeverStmt{Repeat: pos, Body: body, BlockEnd: p.pos}
	default:
		count := p.ParseExpr()
		p.expect(token.Times)
		var body []ast.Stmt
		for p.tok != token.End && p.tok != token.EOF {
			body = append(body, p.parseStmt())
		}
		p.expect(token.End)
		return &ast.RepeatCountStmt{Repeat: pos, Count: count, Body: body, BlockEnd: p.pos}
	}
}

func (p *Parser) parseBranchStmt(tok token.Token) *ast.BranchStmt {
	if p.trace {
		defer un(trace(p, "BranchStmt"))
	}

	pos := p.expect(tok)
	var label *ast.Ident
	if p.tok == token.Ident {
		label = p.parseIdent()
	}

	return &ast.BranchStmt{TokPos: pos, Tok: tok, Label: label}
}
