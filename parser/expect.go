package parser

import "github.com/calico32/goose/token"

func (p *Parser) error(pos token.Pos, msg string) {
	if p.trace {
		defer un(trace(p, "error: "+msg))
	}

	epos := p.file.Position(pos)

	// If AllErrors is not set, discard errors reported on the same line
	// as the last recorded error and stop parsing if there are more than
	// 10 errors.
	// if p.mode&AllErrors == 0 {
	// 	n := len(p.errors)
	// 	if n > 0 && p.errors[n-1].Pos.Line == epos.Line {
	// 		return // discard - likely a spurious error
	// 	}
	// 	if n > 10 {
	// 		panic(bailout{})
	// 	}
	// }

	p.errors.Add(epos, msg)
}

func (p *Parser) errorExpected(pos token.Pos, msg string) {
	msg = "expected " + msg
	if pos == p.pos {
		// the error happened at the current position;
		// make the error message more specific
		switch {
		// case p.tok == token.SEMICOLON && p.lit == "\n":
		// 	msg += ", found newline"
		case p.tok.IsLiteral():
			// print 123 rather than 'INT', etc.
			msg += ", found " + p.lit
		default:
			msg += ", found '" + p.tok.String() + "'"
		}
	}
	p.error(pos, msg)
}

func (p *Parser) expect(tok token.Token) token.Pos {
	pos := p.pos
	if p.tok != tok {
		p.errorExpected(pos, "'"+tok.String()+"'")
	}
	p.next() // make progress
	return pos
}

func (p *Parser) expectMsg(tok token.Token, msg string) token.Pos {
	pos := p.pos
	if p.tok != tok {
		p.errorExpected(pos, msg)
	}
	p.next() // make progress
	return pos
}
