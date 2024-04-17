package parser

import (
	"fmt"
	"io"

	"github.com/calico32/goose/ast"
	"github.com/calico32/goose/scanner"
	"github.com/calico32/goose/token"
)

type Parser struct {
	file    *token.File
	errors  scanner.ErrorList
	scanner scanner.Scanner

	trace       bool
	traceWriter io.Writer
	indent      int

	pos token.Pos
	tok token.Token
	lit string

	nextPos token.Pos
	nextTok token.Token
	nextLit string

	// syncPos   token.Pos
	// syncCount int

	// exprLev int // < 0: in control clause, >= 0: in expression
	// inRhs   bool
}

func (p *Parser) Init(fset *token.FileSet, specifier string, src []byte, trace io.Writer) {
	p.file = fset.AddFile(specifier, -1, len(src))
	p.trace = trace != nil
	if p.trace {
		p.traceWriter = trace
	}

	errorHandler := func(pos token.Position, msg string) { p.errors.Add(pos, msg) }
	p.scanner.Init(p.file, src, errorHandler)

	p.next()
	p.next()
}

func (p *Parser) printTrace(a ...any) {
	const dots = ". . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . "
	const n = len(dots)
	pos := p.file.Position(p.pos)
	fmt.Fprintf(p.traceWriter, "%5d:%3d: ", pos.Line, pos.Column)
	i := 2 * p.indent
	for i > n {
		fmt.Fprint(p.traceWriter, dots)
		i -= n
	}
	// i <= n
	fmt.Fprint(p.traceWriter, dots[0:i])
	fmt.Fprintln(p.traceWriter, a...)
}

func trace(p *Parser, msg string) *Parser {
	if p.traceWriter == nil {
		return p
	}
	p.printTrace(msg, "(")
	p.indent++
	return p
}

// Usage pattern: defer un(trace(p, "..."))
func un(p *Parser) {
	if p.traceWriter == nil {
		return
	}
	p.indent--
	p.printTrace(")")
}

// Advance to the next token.
func (p *Parser) _next() {
	// Because of one-token look-ahead, print the previous token
	// when tracing as it provides a more readable output. The
	// very first token (!p.pos.IsValid()) is not initialized
	// (it is token.ILLEGAL), so don't print it.
	if p.trace && p.pos.IsValid() {
		s := p.tok.String()
		switch {
		case p.tok.IsLiteral():
			p.printTrace(s, p.lit)
		case p.tok.IsOperator(), p.tok.IsKeyword():
			p.printTrace("\"" + s + "\"")
		default:
			p.printTrace(s)
		}
	}

	p.pos, p.tok, p.lit = p.nextPos, p.nextTok, p.nextLit
	p.nextPos, p.nextTok, p.nextLit = p.scanner.Scan()
}

func (p *Parser) next() {
	p._next()

	// skip comments
	for p.tok == token.Comment {
		p._next()
	}
}

// If x is of the form (T), unparen returns unparen(T), otherwise it returns x.
func unparen(x ast.Expr) ast.Expr {
	if p, isParen := x.(*ast.ParenExpr); isParen {
		x = unparen(p.X)
	}
	return x
}

func (p *Parser) safePos(pos token.Pos) (res token.Pos) {
	defer func() {
		if recover() != nil {
			res = token.Pos(p.file.Base() + p.file.Size()) // EOF position
		}
	}()
	_ = p.file.Offset(pos) // trigger a panic if position is out-of-range
	return pos
}

func (p *Parser) parseStmt() (s ast.Stmt) {
	if p.trace {
		defer un(trace(p, "Statement"))
	}

	switch p.tok {
	case
		// tokens that may start an expression
		token.Ident, token.Int, token.Float, token.LParen, token.LBracket, token.LBrace, token.StringStart, token.Null, // operands
		token.Add, token.Sub, token.Mul, token.LogNot, // unary operators
		token.Memo, token.Func, token.Generator, // function declarations
		token.Throw,                                      // throw statement
		token.Await, token.Do, token.Frozen, token.Match: // expression blocks
		s = p.parseSimpleStmt()
	case token.Const:
		s = p.parseConstStmt()
	case token.Let:
		s = p.parseLetStmt()
	case token.Symbol:
		s = p.parseSymbolStmt()
	case token.Return:
		s = p.parseReturnStmt()
	case token.Yield:
		s = p.parseYieldStmt()
	case token.Break, token.Continue:
		s = p.parseBranchStmt(p.tok)
	case token.If:
		n := p.parseIfExprStmt(true)
		if expr, isExpr := n.(*ast.IfExpr); isExpr {
			s = &ast.ExprStmt{X: expr}
		} else {
			s = n.(*ast.IfStmt)
		}
	case token.For:
		s = p.parseForStmt()
	case token.Repeat:
		s = p.parseRepeatStmt()
	case token.Try:
		s = p.parseTryStmt()
	case token.Import:
		s = p.parseImportStmt()
	case token.Export:
		s = p.parseExportStmt()
	case token.Struct:
		s = p.parseStructStmt()
	case token.Operator:
		s = p.parseOperatorStmt()
	case token.Async:
		if p.nextTok == token.Operator {
			s = p.parseOperatorStmt()
		} else {
			s = p.parseSimpleStmt()
		}
	case token.Native:
		s = p.parseNativeStmt()
	default:
		// no statement found
		pos := p.pos
		p.errorExpected(pos, "statement")
		p.next() // make progress
		s = &ast.BadStmt{From: pos, To: p.pos}
	}

	return
}

func (p *Parser) parseSimpleStmt() ast.Stmt {
	if p.trace {
		defer un(trace(p, "SimpleStmt"))
	}

	x := p.ParseExpr()

	switch p.tok {
	case token.Assign, token.AddAssign, token.SubAssign, token.MulAssign, token.QuoAssign, token.RemAssign, token.PowAssign, token.LogAndAssign, token.LogOrAssign, token.BitAndAssign, token.BitOrAssign, token.BitXorAssign, token.BitShlAssign, token.BitShrAssign, token.LogNullAssign:
		tok := p.tok
		tokPos := p.pos
		p.next()
		rhs := p.ParseExpr()
		return &ast.AssignStmt{
			Lhs:    x,
			TokPos: tokPos,
			Tok:    tok,
			Rhs:    rhs,
		}
	case token.Inc, token.Dec:
		s := &ast.IncDecStmt{X: x, TokPos: p.pos, Tok: p.tok}
		p.next()
		return s
	}

	return &ast.ExprStmt{X: x}
}

func (p *Parser) ParseFile() *ast.Module {
	if p.trace {
		defer un(trace(p, "File"))
	}

	// Don't bother parsing the rest if we had errors scanning the first token.
	// Likely not a Go source file at all.
	if p.errors.Len() != 0 {
		return nil
	}

	var stmts []ast.Stmt

	// import decls
	// for p.tok == token.IMPORT {
	// 	decls = append(decls, p.parseGenDecl(token.IMPORT, p.parseImportSpec))
	// }

	for p.tok != token.EOF {
		stmts = append(stmts, p.parseStmt())
	}

	f := &ast.Module{
		Size:      p.file.Size(),
		Stmts:     stmts,
		Specifier: p.file.Specifier(),
		Scheme:    p.file.Scheme(),
	}

	return f
}

func (p *Parser) parseIdent() *ast.Ident {
	pos := p.pos
	name := "_"
	if p.tok == token.Ident {
		name = p.lit
		p.next()
	} else {
		p.expectMsg(token.Ident, "identifier")
	}
	return &ast.Ident{NamePos: pos, Name: name}
}
