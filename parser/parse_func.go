package parser

import (
	"github.com/calico32/goose/ast"
	"github.com/calico32/goose/token"
)

func (p *Parser) parseFuncExpr(asyncPos token.Pos) *ast.FuncExpr {
	if p.trace {
		defer un(trace(p, "FuncExpr"))
	}

	expr := &ast.FuncExpr{}

	if asyncPos.IsValid() {
		expr.Async = asyncPos
	}
	if p.tok == token.Memo {
		expr.Memo = p.pos
		p.next()
	}

	expr.Func = p.expect(token.Func)

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

	if p.tok == token.Arrow {
		expr.Arrow = p.pos
		p.next()
		expr.ArrowExpr = p.ParseExpr()
	} else {
		for p.tok != token.EOF && p.tok != token.End {
			expr.Body = append(expr.Body, p.parseStmt())
		}

		expr.BlockEnd = p.expect(token.End)
	}

	return expr
}

func (p *Parser) parseParameters() (params *ast.FuncParamList) {
	if p.trace {
		defer un(trace(p, "Parameters"))
	}

	opening := p.expect(token.LParen)
	var fields []*ast.FuncParam
	for p.tok != token.RParen {
		f := &ast.FuncParam{}
		if p.tok == token.Ellipsis {
			f.Ellipsis = p.pos
			p.next()
		}
		f.Ident = p.parseIdent()
		if p.tok == token.Assign {
			p.next()
			f.Value = p.ParseExpr()
		}

		fields = append(fields, f)
		if p.tok != token.RParen {
			p.expect(token.Comma)
		}
	}

	rparen := p.expect(token.RParen)
	params = &ast.FuncParamList{Opening: opening, List: fields, Closing: rparen}

	return
}

func (p *Parser) parseCall(fun ast.Expr) *ast.CallExpr {
	if p.trace {
		defer un(trace(p, "CallExpr"))
	}

	lparen := p.expect(token.LParen)
	var list []ast.Expr
	for p.tok != token.RParen && p.tok != token.EOF {
		list = append(list, p.ParseExpr())
		if p.tok != token.RParen {
			p.expect(token.Comma)
		}
	}
	rparen := p.expect(token.RParen)

	return &ast.CallExpr{Func: fun, LParen: lparen, Args: list, RParen: rparen}
}

func (p *Parser) parseReturnStmt() *ast.ReturnStmt {
	if p.trace {
		defer un(trace(p, "ReturnStmt"))
	}

	pos := p.pos
	p.expect(token.Return)
	var x ast.Expr
	if token.IsKeyword(p.tok.String()) {
		// return null
		return &ast.ReturnStmt{Return: pos}
	} else {
		x = p.ParseExpr()
		return &ast.ReturnStmt{Return: pos, Result: x}
	}
}
