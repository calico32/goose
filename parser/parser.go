package parser

import (
	"fmt"
	"io"
	"strconv"

	"github.com/wiisportsresort/goose/ast"
	"github.com/wiisportsresort/goose/scanner"
	"github.com/wiisportsresort/goose/token"
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

	syncPos   token.Pos
	syncCount int

	exprLev int // < 0: in control clause, >= 0: in expression
	inRhs   bool
}

func (p *Parser) init(fset *token.FileSet, filename string, src []byte, trace io.Writer) {
	p.file = fset.AddFile(filename, -1, len(src))
	p.trace = trace != nil
	if p.trace {
		p.traceWriter = trace
	}

	errorHandler := func(pos token.Position, msg string) { p.errors.Add(pos, msg) }
	p.scanner.Init(p.file, src, errorHandler)

	p.next()
}

func (p *Parser) printTrace(a ...any) {
	const dots = ". . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . "
	const n = len(dots)
	pos := p.file.Position(p.pos)
	fmt.Fprintf(p.traceWriter, "%5d:%3d: ", pos.Line, pos.Column)
	i := 2 * p.indent
	for i > n {
		fmt.Fprintf(p.traceWriter, dots)
		i -= n
	}
	// i <= n
	fmt.Fprint(p.traceWriter, dots[0:i])
	fmt.Fprintln(p.traceWriter, a...)
}

func trace(p *Parser, msg string) *Parser {
	p.printTrace(msg, "(")
	p.indent++
	return p
}

// Usage pattern: defer un(trace(p, "..."))
func un(p *Parser) {
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

	p.pos, p.tok, p.lit = p.scanner.Scan()
}

func (p *Parser) next() {
	p._next()

	// skip comments
	for p.tok == token.Comment {
		p._next()
	}
}

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

func (p *Parser) parseIdent() *ast.Ident {
	pos := p.pos
	name := "_"
	if p.tok == token.Ident {
		name = p.lit
		p.next()
	} else {
		p.expect(token.Ident)
	}
	return &ast.Ident{NamePos: pos, Name: name}
}

func (p *Parser) parseParameters() (params *ast.FieldList) {
	if p.trace {
		defer un(trace(p, "Parameters"))
	}

	opening := p.expect(token.LParen)
	var fields []*ast.Field
	for p.tok != token.RParen {
		ident := p.parseIdent()
		fields = append(fields, &ast.Field{Ident: ident})
		if p.tok != token.RParen {
			p.expect(token.Comma)
		}
	}

	rparen := p.expect(token.RParen)
	params = &ast.FieldList{Opening: opening, List: fields, Closing: rparen}

	return
}

func (p *Parser) parseFuncStmt() *ast.FuncStmt {
	if p.trace {
		defer un(trace(p, "FuncStmt"))
	}

	var memoPos token.Pos
	if p.tok == token.Memo {
		memoPos = p.pos
		p.next()
	}

	pos := p.expect(token.Func)

	ident := p.parseIdent()
	params := p.parseParameters()

	var body []ast.Stmt

	for p.tok != token.EOF && p.tok != token.End {
		body = append(body, p.parseStmt())
	}

	end := p.expect(token.End)

	return &ast.FuncStmt{
		Memo:     memoPos,
		Func:     pos,
		Name:     ident,
		Params:   params,
		Body:     body,
		BlockEnd: end,
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

func (p *Parser) parseCall(fun ast.Expr) *ast.CallExpr {
	if p.trace {
		defer un(trace(p, "CallExpr"))
	}

	lparen := p.expect(token.LParen)
	var list []ast.Expr
	for p.tok != token.RParen && p.tok != token.EOF {
		list = append(list, p.parseExpr())
		if p.tok != token.RParen {
			p.expect(token.Comma)
		}
	}
	rparen := p.expect(token.RParen)

	return &ast.CallExpr{Fun: fun, LParen: lparen, Args: list, RParen: rparen}
}

func (p *Parser) parseIndexOrSlice(x ast.Expr) ast.Expr {
	if p.trace {
		defer un(trace(p, "IndexOrSlice"))
	}

	lbrack := p.expect(token.LBracket)
	if p.tok == token.RBracket {
		// empty index
		p.errorExpected(p.pos, "expression")
		rbrack := p.pos
		p.next()
		return &ast.IndexExpr{
			X:      x,
			LBrack: lbrack,
			Index:  &ast.BadExpr{From: lbrack, To: rbrack},
			RBrack: rbrack,
		}
	}

	var left ast.Expr
	var right ast.Expr
	isSlicing := false

	if p.tok != token.Colon {
		left = p.parseExpr()
		if p.tok == token.Colon {
			isSlicing = true
			p.next()
		}
	} else {
		isSlicing = true
	}

	if p.tok != token.RBracket {
		right = p.parseExpr()
	}

	if left == nil && right == nil {
		p.errorExpected(p.pos, "slice expression")
		return &ast.SliceExpr{
			X:      x,
			LBrack: lbrack,
			Low:    &ast.BadExpr{From: lbrack, To: p.pos},
			High:   &ast.BadExpr{From: p.pos, To: p.pos},
			RBrack: p.safePos(p.pos),
		}
	}

	rbrack := p.expect(token.RBracket)

	if isSlicing {
		return &ast.SliceExpr{
			X:      x,
			LBrack: lbrack,
			Low:    left,
			High:   right,
			RBrack: rbrack,
		}
	}

	return &ast.IndexExpr{
		X:      x,
		LBrack: lbrack,
		Index:  left,
		RBrack: rbrack,
	}
}

func (p *Parser) parseExpr() ast.Expr {
	if p.trace {
		defer un(trace(p, "Expression"))
	}

	return p.parseBinaryExpr(nil, token.PrecedenceLowest+1)
}

func (p *Parser) parseBinaryExpr(x ast.Expr, precedence int) ast.Expr {
	if p.trace {
		defer un(trace(p, "BinaryExpr"))
	}

	if x == nil {
		x = p.parseUnaryExpr()
	}

	for {
		operator := p.tok
		operatorPrec := p.tok.Precedence()
		if operatorPrec < precedence {
			return x
		}
		pos := p.expect(operator)
		y := p.parseBinaryExpr(nil, operatorPrec+1)

		x = &ast.BinaryExpr{X: x, OpPos: pos, Op: operator, Y: y}
	}

}

func (p *Parser) parseUnaryExpr() ast.Expr {
	switch p.tok {
	case token.Add, token.Sub, token.LogNot:
		pos := p.pos
		op := p.tok
		p.next()
		x := p.parseUnaryExpr()
		return &ast.UnaryExpr{OpPos: pos, Op: op, X: x}
	}

	return p.parsePrimaryExpr(nil)
}

func (p *Parser) parsePrimaryExpr(x ast.Expr) ast.Expr {
	if p.trace {
		defer un(trace(p, "PrimaryExpr"))
	}

	if x == nil {
		x = p.parseOperand()
	}

	for {
		switch p.tok {
		case token.Int:
			if x, ok := x.(*ast.ArrayLiteral); ok && len(x.List) == 0 {
				nulls := []ast.Expr{}
				count, err := strconv.Atoi(p.lit)
				if err != nil {
					p.error(p.pos, "invalid array length")
				}
				for i := 0; i < count; i++ {
					nulls = append(nulls, &ast.Literal{Kind: token.Null})
				}
				x.List = nulls
				p.next()
			}
		case token.LParen:
			x = p.parseCall(x)
		case token.LBracket:
			x = p.parseIndexOrSlice(x)
		}
		return x
	}
}

func (p *Parser) parseOperand() ast.Expr {
	if p.trace {
		defer un(trace(p, "Operand"))
	}

	switch p.tok {
	case token.Ident:
		x := p.parseIdent()
		return x
	case token.Int, token.Float, token.Null:
		x := &ast.Literal{Value: p.lit, ValuePos: p.pos, Kind: p.tok}
		p.next()
		return x
	case token.StringStart:
		x := p.parseString()
		return x
	case token.LBracket:
		x := p.parseArrayLitOrInitializer()
		return x
	case token.LParen:
		lparen := p.pos
		p.next()
		x := p.parseExpr()
		rparen := p.expect(token.RParen)
		return &ast.ParenExpr{Lparen: lparen, X: x, Rparen: rparen}
	}

	pos := p.pos
	p.errorExpected(pos, "operand")
	p.next() // make progress
	return &ast.BadExpr{From: pos, To: p.safePos(p.pos)}
}

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

func (p *Parser) parseStmt() (s ast.Stmt) {
	if p.trace {
		defer un(trace(p, "Statement"))
	}

	switch p.tok {
	case token.Func, token.Memo:
		s = p.parseFuncStmt()
	case
		// tokens that may start an expression
		token.Ident, token.Int, token.Float, token.LParen, token.LBracket, token.LBrace, // operands
		token.Add, token.Sub, token.Mul, token.LogNot: // unary operators
		s = p.parseSimpleStmt()
	case token.Const:
		s = p.parseConstStmt()
	case token.Let:
		s = p.parseLetStmt()
	case token.Return:
		s = p.parseReturnStmt()
	case token.Break, token.Continue:
		s = p.parseBranchStmt(p.tok)
	case token.If:
		s = p.parseIfStmt()
	case token.For:
		s = p.parseForStmt()
	case token.Repeat:
		s = p.parseRepeatStmt()
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

	x := p.parseExpr()

	switch p.tok {
	case token.Assign, token.AddAssign, token.SubAssign, token.MulAssign, token.QuoAssign, token.RemAssign, token.PowAssign:
		tok := p.tok
		tokPos := p.pos
		p.next()
		rhs := p.parseExpr()
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

func (p *Parser) parseConstStmt() *ast.ConstStmt {
	if p.trace {
		defer un(trace(p, "ConstStmt"))
	}

	pos := p.pos
	decl := p.tok
	p.next()
	if decl != token.Const {
		p.errorExpected(pos, "const")
		decl = token.Let
	}

	lhs := p.parseIdent()
	tokPos := p.expect(token.Assign)
	rhs := p.parseExpr()

	return &ast.ConstStmt{
		ConstPos: pos,
		Ident:    lhs,
		TokPos:   tokPos,
		Value:    rhs,
	}
}

func (p *Parser) parseLetStmt() *ast.LetStmt {
	if p.trace {
		defer un(trace(p, "LetStmt"))
	}

	pos := p.pos
	decl := p.tok
	p.next()
	if decl != token.Let {
		p.errorExpected(pos, "let")
		decl = token.Let
	}

	lhs := p.parseIdent()

	var tokPos token.Pos
	var rhs ast.Expr

	if p.tok == token.Assign {
		tokPos = p.expect(token.Assign)
		rhs = p.parseExpr()
	}

	return &ast.LetStmt{
		LetPos: pos,
		Ident:  lhs,
		TokPos: tokPos,
		Value:  rhs,
	}
}

func (p *Parser) parseReturnStmt() *ast.ReturnStmt {
	if p.trace {
		defer un(trace(p, "ReturnStmt"))
	}

	pos := p.pos
	p.expect(token.Return)
	var x ast.Expr
	if p.tok != token.End && p.tok != token.EOF {
		x = p.parseExpr()
	}

	return &ast.ReturnStmt{Return: pos, Result: x}
}

func (p *Parser) parseBranchStmt(tok token.Token) *ast.BranchStmt {
	if p.trace {
		defer un(trace(p, "BranchStmt"))
	}

	pos := p.expect(tok)
	var label *ast.Ident
	if p.tok == token.Ident {
		label = p.parseIdent()
	}

	return &ast.BranchStmt{TokPos: pos, Tok: tok, Label: label}
}

func (p *Parser) parseIfStmt() *ast.IfStmt {
	if p.trace {
		defer un(trace(p, "IfStmt"))
	}

	pos := p.expect(token.If)
	cond := p.parseExpr()
	var body []ast.Stmt
	for p.tok != token.End && p.tok != token.EOF && p.tok != token.Else {
		body = append(body, p.parseStmt())
	}

	var else_ []ast.Stmt
	if p.tok == token.Else {
		p.next()
		if p.tok == token.If {
			else_ = []ast.Stmt{p.parseIfStmt()}
		} else {
			for p.tok != token.End && p.tok != token.EOF {
				else_ = append(else_, p.parseStmt())
			}
			p.expect(token.End)
		}
	} else {
		p.expect(token.End)
	}

	return &ast.IfStmt{If: pos, Cond: cond, Body: body, Else: else_, BlockEnd: p.pos}
}

func (p *Parser) parseForStmt() ast.Stmt {
	if p.trace {
		defer un(trace(p, "ForStmt"))
	}

	pos := p.expect(token.For)
	ident := p.parseIdent()
	p.expect(token.In)
	expr := p.parseExpr()

	var body []ast.Stmt
	for p.tok != token.End && p.tok != token.EOF {
		body = append(body, p.parseStmt())
	}
	p.expect(token.End)

	return &ast.ForStmt{For: pos, Var: ident, Iterable: expr, Body: body, BlockEnd: p.pos}
}

func (p *Parser) parseRepeatStmt() ast.Stmt {
	if p.trace {
		defer un(trace(p, "RepeatStmt"))
	}

	pos := p.expect(token.Repeat)
	switch p.tok {
	case token.While:
		p.next()
		cond := p.parseExpr()
		var body []ast.Stmt
		for p.tok != token.End && p.tok != token.EOF {
			body = append(body, p.parseStmt())
		}
		p.expect(token.End)
		return &ast.RepeatWhileStmt{Repeat: pos, Cond: cond, Body: body, BlockEnd: p.pos}
	case token.Forever:
		p.next()
		var body []ast.Stmt
		for p.tok != token.End && p.tok != token.EOF {
			body = append(body, p.parseStmt())
		}
		p.expect(token.End)
		return &ast.RepeatForeverStmt{Repeat: pos, Body: body, BlockEnd: p.pos}
	default:
		count := p.parseExpr()
		p.expect(token.Times)
		var body []ast.Stmt
		for p.tok != token.End && p.tok != token.EOF {
			body = append(body, p.parseStmt())
		}
		p.expect(token.End)
		return &ast.RepeatCountStmt{Repeat: pos, Count: count, Body: body, BlockEnd: p.pos}
	}
}

func (p *Parser) parseFile() *ast.File {
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

	f := &ast.File{
		Size:  p.file.Size(),
		Stmts: stmts,
	}

	return f
}
