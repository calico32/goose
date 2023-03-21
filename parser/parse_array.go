package parser

import (
	"github.com/calico32/goose/ast"
	"github.com/calico32/goose/token"
)

func (p *Parser) parseArrayLitOrInitializer() (s ast.Expr) {
	if p.trace {
		defer un(trace(p, "ArrayLiteralOrInitializer"))
	}

	isInitializer := false
	var initializer ast.ArrayInitializer

	lbracket := p.expect(token.LBracket)
	var list []ast.Expr
	for p.tok != token.RBracket && p.tok != token.EOF {
		if isInitializer {
			// already read [value;
			// read count (bracket coming next)
			initializer.Count = p.parseExpr()
			break
		} else {
			list = append(list, p.parseExpr())
		}

		if p.tok != token.RBracket {
			if p.tok == token.Semi {
				if len(list) != 1 {
					// initializer semicolon too late (mixed comma and semicolon)
					p.error(p.pos, "expected ']' or ','")
					p.next()
				} else {
					isInitializer = true
					initializer.Value = list[0]
					initializer.Semi = p.pos
					p.next()
				}
				continue
			}

			p.expect(token.Comma)
		}
	}
	rbracket := p.expect(token.RBracket)

	if isInitializer {
		return &initializer
	}

	return &ast.ArrayLiteral{
		Opening: lbracket,
		List:    list,
		Closing: rbracket,
	}
}
