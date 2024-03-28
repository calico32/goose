package parser

import (
	"github.com/calico32/goose/ast"
	"github.com/calico32/goose/token"
)

func (p *Parser) parseStructStmt() *ast.StructStmt {
	if p.trace {
		defer un(trace(p, "StructStmt"))
	}

	stmt := &ast.StructStmt{}
	stmt.Struct = p.expect(token.Struct)
	stmt.Name = p.parseIdent()
	stmt.Fields = p.parseStructFields()

	if p.tok == token.Init {
		stmt.Init = p.parseStructInit()
	}

	return stmt
}

func (p *Parser) parseStructFields() (params *ast.StructFieldList) {
	if p.trace {
		defer un(trace(p, "StructFieldList"))
	}

	opening := p.expect(token.LParen)
	var fields []*ast.StructField
	for p.tok != token.RParen {
		ident := p.parseIdent()
		f := &ast.StructField{Ident: ident}
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
	params = &ast.StructFieldList{Opening: opening, List: fields, Closing: rparen}

	return
}

func (p *Parser) parseStructInit() *ast.StructInit {
	if p.trace {
		defer un(trace(p, "StructInit"))
	}

	init := &ast.StructInit{}
	init.Init = p.expect(token.Init)

	for p.tok != token.End && p.tok != token.EOF {
		init.Body = append(init.Body, p.parseStmt())
	}

	init.BlockEnd = p.expect(token.End)

	return init
}

func (p *Parser) parseBracketPropertyExpr() *ast.BracketPropertyExpr {
	if p.trace {
		defer un(trace(p, "BracketProperty"))
	}

	prop := &ast.BracketPropertyExpr{}
	prop.HashLBracket = p.expect(token.HashLBracket)
	prop.X = p.ParseExpr()
	prop.RBracket = p.expect(token.RBracket)

	return prop
}

func (p *Parser) parseOperatorStmt() *ast.OperatorStmt {
	if p.trace {
		defer un(trace(p, "OperatorStmt"))
	}

	stmt := &ast.OperatorStmt{}
	if p.tok == token.Async {
		stmt.Async = p.pos
		p.next()
	}
	stmt.Operator = p.expect(token.Operator)
	stmt.Receiver = p.parseIdent()

	if p.tok < token.OverloadAllowedStart || p.tok > token.OverloadAllowedEnd {
		p.errorExpected(stmt.Pos(), "overloadable operator")
	}

	stmt.TokPos = p.pos
	stmt.Tok = p.tok
	// do not advance if the token is a parenthesis because it is part of the
	// parameters
	if p.tok != token.LParen {
		p.next()
	}

	stmt.Params = p.parseParameters()

	if p.tok == token.Arrow {
		stmt.Arrow = p.pos
		p.next()
		stmt.ArrowExpr = p.ParseExpr()
	} else {
		for p.tok != token.EOF && p.tok != token.End {
			stmt.Body = append(stmt.Body, p.parseStmt())
		}

		stmt.BlockEnd = p.expect(token.End)
	}

	return stmt
}
