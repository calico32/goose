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

	if len(start) > 1 {
		start = start[1:]
	}

	str := &ast.StringLiteral{
		StringStart: &ast.StringLiteralStart{
			Quote:   quote,
			Content: start,
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
				Ident: &ast.Ident{
					NamePos: p.pos + 1,
					Name:    p.lit[1:],
				},
			})
			p.next()
		case token.StringInterpExprStart:
			pos := p.pos
			p.next()
			interpExpr := &ast.StringLiteralInterpExpr{
				InterpPos: pos,
				Expr:      p.ParseExpr(),
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
