package parser

import (
	"github.com/calico32/goose/ast"
	"github.com/calico32/goose/token"
)

func (p *Parser) parseConstStmt() *ast.ConstStmt {
	if p.trace {
		defer un(trace(p, "ConstStmt"))
	}

	pos := p.pos
	decl := p.tok
	p.next()
	if decl != token.Const {
		p.errorExpected(pos, "const")
		decl = token.Let
	}

	lhs := p.parseIdent()
	var rhs ast.Expr
	tokPos := p.pos
	if p.tok != token.Assign {
		p.error(p.pos, "const declaration must be followed by an assignment")
		rhs = &ast.BadExpr{From: p.pos, To: p.pos}
	} else {
		p.next()
		rhs = p.parseExpr()
	}

	return &ast.ConstStmt{
		ConstPos: pos,
		Ident:    lhs,
		TokPos:   tokPos,
		Value:    rhs,
	}
}

func (p *Parser) parseLetStmt() *ast.LetStmt {
	if p.trace {
		defer un(trace(p, "LetStmt"))
	}

	pos := p.pos
	decl := p.tok
	p.next()
	if decl != token.Let {
		p.errorExpected(pos, "let")
		decl = token.Let
	}

	lhs := p.parseIdent()

	var tokPos token.Pos
	var rhs ast.Expr

	if p.tok == token.Assign {
		tokPos = p.expect(token.Assign)
		rhs = p.parseExpr()
	}

	return &ast.LetStmt{
		LetPos: pos,
		Ident:  lhs,
		TokPos: tokPos,
		Value:  rhs,
	}
}
