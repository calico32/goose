package parser

import (
	"github.com/calico32/goose/ast"
	"github.com/calico32/goose/token"
)

func (p *Parser) parseString() (s *ast.StringLiteral) {
	if p.trace {
		defer un(trace(p, "String"))
	}

	start := p.lit
	quote := p.expect(token.StringStart)
	var parts []ast.StringLiteralPart

	str := &ast.StringLiteral{
		StringStart: &ast.StringLiteralStart{
			Quote:   quote,
			Content: start[1:],
		},
	}

loop:
	for {
		switch p.tok {
		case token.StringMid:
			parts = append(parts, &ast.StringLiteralMiddle{
				StartPos: p.pos,
				Content:  p.lit,
			})
			p.next()
		case token.StringInterpIdent:
			parts = append(parts, &ast.StringLiteralInterpIdent{
				InterpPos: p.pos,
				Name:      p.lit[1:],
			})
			p.next()
		case token.StringInterpExprStart:
			pos := p.pos
			p.next()
			interpExpr := &ast.StringLiteralInterpExpr{
				InterpPos: pos,
				Expr:      p.parseExpr(),
			}
			p.expect(token.StringInterpExprEnd)
			parts = append(parts, interpExpr)
		case token.StringEnd:
			str.StringEnd = &ast.StringLiteralEnd{
				StartPos: p.pos,
				Content:  p.lit[:len(p.lit)-1],
				Quote:    token.Pos(int(p.pos) + len(p.lit)),
			}
			p.next()
			break loop
		default:
			p.errorExpected(p.pos, "string part")
			p.next()
			break loop
		}
	}

	str.Parts = parts
	return str
}

func (p *Parser) parseSymbolStmt() (stmt *ast.SymbolStmt) {
	if p.trace {
		defer un(trace(p, "SymbolStmt"))
	}

	stmt = &ast.SymbolStmt{}
	stmt.Symbol = p.expect(token.Symbol)
	stmt.Ident = p.parseIdent()

	return stmt
}
