package parser

import (
	"github.com/calico32/goose/ast"
	"github.com/calico32/goose/token"
)

func (p *Parser) parseComposite() (s *ast.CompositeLiteral) {
	if p.trace {
		defer un(trace(p, "CompositeLiteral"))
	}

	lbrace := p.expect(token.LBrace)

	var list []*ast.CompositeField

	for p.tok != token.RBrace && p.tok != token.EOF {
		list = append(list, p.parseCompositeField())
		if p.tok != token.RBrace {
			p.expect(token.Comma)
		}
	}

	rbrace := p.expect(token.RBrace)

	return &ast.CompositeLiteral{
		Lbrace: lbrace,
		Fields: list,
		Rbrace: rbrace,
	}
}

func (p *Parser) parseCompositeField() (s *ast.CompositeField) {
	if p.trace {
		defer un(trace(p, "CompositeField"))
	}

	var key ast.Expr
	var value ast.Expr

	if p.tok == token.LBracket {
		p.next()
		key = p.ParseExpr()
		p.expect(token.RBracket)
	} else if p.tok == token.StringStart {
		key = p.parseString()
	} else if p.tok == token.Int || p.tok == token.Float {
		key = &ast.Literal{Value: p.lit, ValuePos: p.pos, Kind: p.tok}
		p.next()
	} else {
		ident := p.parseIdent()
		key = &ast.StringLiteral{
			StringStart: &ast.StringLiteralStart{
				Quote:   token.Pos(int(ident.Pos()) - 1),
				Content: ident.Name,
			},
			StringEnd: &ast.StringLiteralEnd{
				StartPos: token.Pos(int(ident.End()) - 1),
				Quote:    token.Pos(int(ident.End()) + 1),
			},
		}
	}

	p.expect(token.Colon)
	value = p.ParseExpr()

	return &ast.CompositeField{
		Key:   key,
		Value: value,
	}
}
