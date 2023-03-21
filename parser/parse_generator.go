package parser

import (
	"github.com/calico32/goose/ast"
	"github.com/calico32/goose/token"
)

func (p *Parser) parseRangeExpr(low ast.Expr) *ast.RangeExpr {
	if p.trace {
		defer un(trace(p, "RangeExpr"))
	}

	toPos := p.expect(token.To)
	high := p.parseExpr()

	stepPos := token.NoPos
	var step ast.Expr

	if p.tok == token.Step {
		stepPos = p.expect(token.Step)
		step = p.parseExpr()
	}

	return &ast.RangeExpr{
		Low:     low,
		ToPos:   toPos,
		High:    high,
		StepPos: stepPos,
		Step:    step,
	}
}

func (p *Parser) parseGeneratorExpr(asyncPos token.Pos) *ast.GeneratorExpr {
	if p.trace {
		defer un(trace(p, "GeneratorExpr"))
	}

	expr := &ast.GeneratorExpr{}

	if asyncPos.IsValid() {
		expr.Async = asyncPos
	}

	expr.Generator = p.expect(token.Generator)

	if p.tok == token.Ident {
		part := p.parseIdent()
		if p.tok == token.Period {
			expr.Receiver = part
			p.next()
			expr.Name = p.parseIdent()
		} else {
			expr.Name = part
		}
	}

	expr.Params = p.parseParameters()

	for p.tok != token.EOF && p.tok != token.End {
		expr.Body = append(expr.Body, p.parseStmt())
	}

	expr.BlockEnd = p.expect(token.End)

	return expr
}

func (p *Parser) parseYieldStmt() *ast.ReturnStmt {
	if p.trace {
		defer un(trace(p, "YieldStmt"))
	}

	pos := p.pos
	p.expect(token.Yield)
	x := p.parseExpr()
	return &ast.ReturnStmt{Return: pos, Result: x}
}
