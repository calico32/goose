package parser

import (
	"github.com/calico32/goose/ast"
	"github.com/calico32/goose/token"
)

func (p *Parser) parseNativeStmt() ast.NativeStmt {
	native := p.expect(token.Native)

	var async token.Pos
	switch p.tok {
	case token.Const:
		stmt := &ast.NativeConst{
			Native: native,
			Const:  p.expect(token.Const),
		}

		stmt.Ident = p.parseIdent()
		return stmt
	case token.Async:
		async = p.expect(token.Async)
		if p.tok == token.Operator {
			stmt := &ast.NativeOperator{
				Native:   native,
				Async:    async,
				Operator: p.expect(token.Operator),
			}

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

			return stmt
		}

		fallthrough // function
	case token.Memo, token.Func:
		stmt := &ast.NativeFunc{
			Native: native,
			Async:  async,
		}

		if p.tok == token.Async {
			stmt.Async = p.expect(token.Async)
		}

		if p.tok == token.Memo {
			stmt.Memo = p.expect(token.Memo)
		}

		stmt.Fn = p.expect(token.Func)

		if p.tok == token.Ident {
			part := p.parseIdent()
			if p.tok == token.Period {
				stmt.Receiver = part
				p.next()
				stmt.Name = p.parseIdent()
			} else {
				stmt.Name = part
			}
		}

		stmt.Params = p.parseParameters()

		return stmt
	case token.Struct:
		stmt := &ast.NativeStruct{
			Native: native,
			Struct: p.expect(token.Struct),
		}

		stmt.Name = p.parseIdent()
		stmt.Fields = p.parseStructFields()

		return stmt
	default:
		p.errorExpected(native, "const/struct declaration or function/operator signature")
		return nil
	}
}
