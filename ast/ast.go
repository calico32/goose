package ast

import (
	"github.com/wiisportsresort/goose/token"
)

func isWhitespace(ch byte) bool { return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' }

func trimRight(s string) string {
	i := len(s)
	for i > 0 && isWhitespace(s[i-1]) {
		i--
	}
	return s[0:i]
}

type Token struct {
	Type   token.Token
	Source string
}

type Node interface {
	// Pos returns the position of the first character belonging to the node
	Pos() token.Pos
	// End returns the position of the first character immediately after the node
	End() token.Pos
}

type Expr interface {
	Node
	exprNode()
}

type Stmt interface {
	Node
	stmtNode()
}

type StringLiteralPart interface {
	Node
	stringLiteralPart()
}

type Comment struct {
	Slash token.Pos
	Text  string
}

func (c *Comment) Pos() token.Pos { return c.Slash }
func (c *Comment) End() token.Pos { return token.Pos(int(c.Slash) + len(c.Text)) }

// type CommentGroup struct{ List []*Comment }

// func (g *CommentGroup) Pos() token.Pos { return g.List[0].Pos() }
// func (g *CommentGroup) End() token.Pos { return g.List[len(g.List)-1].End() }

// func (g *CommentGroup) Text() string {
// 	if g == nil {
// 		return ""
// 	}
// 	comments := make([]string, len(g.List))
// 	for i, c := range g.List {
// 		comments[i] = c.Text
// 	}

// 	lines := make([]string, 0, 10) // most comments are less than 10 lines
// 	for _, c := range comments {
// 		// Remove comment markers.
// 		// The parser has given us exactly the comment text.
// 		switch c[1] {
// 		case '/':
// 			//-style comment (no newline at the end)
// 			c = c[2:]
// 			if len(c) == 0 {
// 				// empty line
// 				break
// 			}
// 			if c[0] == ' ' {
// 				// strip first space - required for Example tests
// 				c = c[1:]
// 				break
// 			}

// 		case '*':
// 			/*-style comment */
// 			c = c[2 : len(c)-2]
// 		}

// 		// Split on newlines.
// 		cl := strings.Split(c, "\n")

// 		// Walk lines, stripping trailing white space and adding to list.
// 		for _, l := range cl {
// 			lines = append(lines, trimRight(l))
// 		}
// 	}

// 	// Remove leading blank lines; convert runs of
// 	// interior blank lines to a single blank line.
// 	n := 0
// 	for _, line := range lines {
// 		if line != "" || n > 0 && lines[n-1] != "" {
// 			lines[n] = line
// 			n++
// 		}
// 	}
// 	lines = lines[0:n]

// 	// Add final "" entry to get trailing newline from Join.
// 	if n > 0 && lines[n-1] != "" {
// 		lines = append(lines, "")
// 	}

// 	return strings.Join(lines, "\n")
// }

type CompositeField struct {
	Key   Expr
	Value Expr
}

type FuncParam struct {
	Ident *Ident
	Value Expr
}

func (f *FuncParam) Pos() token.Pos { return f.Ident.Pos() }
func (f *FuncParam) End() token.Pos { return f.Ident.End() }

type FuncParamList struct {
	Opening token.Pos
	List    []*FuncParam
	Closing token.Pos
}

func (f *FuncParamList) Pos() token.Pos {
	if f.Opening.IsValid() {
		return f.Opening
	}
	if len(f.List) > 0 {
		return f.List[0].Pos()
	}
	return token.NoPos
}

func (f *FuncParamList) End() token.Pos {
	if f.Closing.IsValid() {
		return f.Closing + 1
	}
	if n := len(f.List); n > 0 {
		return f.List[n-1].End()
	}
	return token.NoPos
}

func (f *FuncParamList) NumFields() int { return len(f.List) }

type (
	BadExpr struct {
		From, To token.Pos
	}

	Ident struct {
		NamePos token.Pos
		Name    string
	}

	Literal struct {
		ValuePos token.Pos
		Kind     token.Token
		Value    string
	}

	StringLiteral struct {
		StringStart *StringLiteralStart
		Parts       []StringLiteralPart
		StringEnd   *StringLiteralEnd
	}

	StringLiteralStart struct {
		Quote   token.Pos
		Content string
	}

	StringLiteralMiddle struct {
		StartPos token.Pos
		Content  string
	}

	StringLiteralInterpIdent struct {
		InterpPos token.Pos
		Name      string
	}

	StringLiteralInterpExpr struct {
		InterpPos token.Pos
		Expr      Expr
	}

	StringLiteralEnd struct {
		StartPos token.Pos
		Content  string
		Quote    token.Pos
	}

	ArrayLiteral struct {
		Opening token.Pos
		List    []Expr
		Closing token.Pos
	}

	ArrayInitializer struct {
		Opening token.Pos
		Value   Expr
		Semi    token.Pos
		Count   Expr
		Closing token.Pos
	}

	CompositeLiteral struct {
		Lbrace token.Pos
		Fields []*CompositeField
		Rbrace token.Pos
	}

	ParenExpr struct {
		Lparen token.Pos
		X      Expr
		Rparen token.Pos
	}

	SelectorExpr struct {
		X   Expr
		Sel *Ident
	}

	BracketSelectorExpr struct {
		X      Expr
		LBrack token.Pos
		Sel    Expr
		RBrack token.Pos
	}

	SliceExpr struct {
		X      Expr
		LBrack token.Pos
		Low    Expr
		High   Expr
		RBrack token.Pos
	}

	CallExpr struct {
		Fun    Expr
		LParen token.Pos
		Args   []Expr
		RParen token.Pos
	}

	UnaryExpr struct {
		OpPos token.Pos
		Op    token.Token
		X     Expr
	}

	BinaryExpr struct {
		X     Expr
		OpPos token.Pos
		Op    token.Token
		Y     Expr
	}

	KeyValueExpr struct {
		Key   Expr
		Colon token.Pos
		Value Expr
	}
)

func (x *BadExpr) Pos() token.Pos                  { return x.From }
func (x *Ident) Pos() token.Pos                    { return x.NamePos }
func (x *Literal) Pos() token.Pos                  { return x.ValuePos }
func (x *StringLiteral) Pos() token.Pos            { return x.StringStart.Pos() }
func (x *StringLiteralStart) Pos() token.Pos       { return x.Quote }
func (x *StringLiteralMiddle) Pos() token.Pos      { return x.StartPos }
func (x *StringLiteralInterpIdent) Pos() token.Pos { return x.InterpPos }
func (x *StringLiteralInterpExpr) Pos() token.Pos  { return x.InterpPos }
func (x *StringLiteralEnd) Pos() token.Pos         { return x.StartPos }
func (x *ArrayLiteral) Pos() token.Pos             { return x.Opening }
func (x *ArrayInitializer) Pos() token.Pos         { return x.Opening }
func (x *CompositeLiteral) Pos() token.Pos         { return x.Lbrace }
func (x *ParenExpr) Pos() token.Pos                { return x.Lparen }
func (x *SelectorExpr) Pos() token.Pos             { return x.X.Pos() }
func (x *BracketSelectorExpr) Pos() token.Pos      { return x.X.Pos() }
func (x *SliceExpr) Pos() token.Pos                { return x.X.Pos() }
func (x *CallExpr) Pos() token.Pos                 { return x.Fun.Pos() }
func (x *UnaryExpr) Pos() token.Pos                { return x.OpPos }
func (x *BinaryExpr) Pos() token.Pos               { return x.X.Pos() }
func (x *KeyValueExpr) Pos() token.Pos             { return x.Key.Pos() }

func (x *BadExpr) End() token.Pos             { return x.To }
func (x *Ident) End() token.Pos               { return token.Pos(int(x.NamePos) + len(x.Name)) }
func (x *Literal) End() token.Pos             { return token.Pos(int(x.ValuePos) + len(x.Value)) }
func (x *StringLiteral) End() token.Pos       { return x.StringEnd.Quote + 1 }
func (x *StringLiteralStart) End() token.Pos  { return token.Pos(int(x.Quote) + len(x.Content)) }
func (x *StringLiteralMiddle) End() token.Pos { return token.Pos(int(x.StartPos) + len(x.Content)) }
func (x *StringLiteralInterpIdent) End() token.Pos {
	return token.Pos(int(x.InterpPos)+len(x.Name)) + 1
}
func (x *StringLiteralInterpExpr) End() token.Pos { return x.Expr.End() }
func (x *StringLiteralEnd) End() token.Pos        { return x.Quote + 1 }
func (x *ArrayLiteral) End() token.Pos            { return x.Closing + 1 }
func (x *ArrayInitializer) End() token.Pos        { return x.Closing + 1 }
func (x *CompositeLiteral) End() token.Pos        { return x.Rbrace + 1 }
func (x *ParenExpr) End() token.Pos               { return x.Rparen + 1 }
func (x *SelectorExpr) End() token.Pos            { return x.Sel.End() }
func (x *BracketSelectorExpr) End() token.Pos     { return x.RBrack + 1 }
func (x *SliceExpr) End() token.Pos               { return x.RBrack + 1 }
func (x *CallExpr) End() token.Pos                { return x.RParen + 1 }
func (x *UnaryExpr) End() token.Pos               { return x.X.End() }
func (x *BinaryExpr) End() token.Pos              { return x.Y.End() }
func (x *KeyValueExpr) End() token.Pos            { return x.Value.End() }

func (*BadExpr) exprNode()                           {}
func (*Ident) exprNode()                             {}
func (*Literal) exprNode()                           {}
func (*StringLiteral) exprNode()                     {}
func (*StringLiteralStart) stringLiteralPart()       {}
func (*StringLiteralMiddle) stringLiteralPart()      {}
func (*StringLiteralInterpIdent) stringLiteralPart() {}
func (*StringLiteralInterpExpr) stringLiteralPart()  {}
func (*StringLiteralEnd) stringLiteralPart()         {}
func (*ArrayLiteral) exprNode()                      {}
func (*ArrayInitializer) exprNode()                  {}
func (*CompositeLiteral) exprNode()                  {}
func (*ParenExpr) exprNode()                         {}
func (*SelectorExpr) exprNode()                      {}
func (*BracketSelectorExpr) exprNode()               {}
func (*SliceExpr) exprNode()                         {}
func (*CallExpr) exprNode()                          {}
func (*UnaryExpr) exprNode()                         {}
func (*BinaryExpr) exprNode()                        {}
func (*KeyValueExpr) exprNode()                      {}

func (id *Ident) String() string {
	if id != nil {
		return id.Name
	}
	return "<nil>"
}

type (
	BadStmt struct {
		From, To token.Pos
	}

	EmptyStmt struct {
	}

	LabeledStmt struct {
		Label *Ident
		Colon token.Pos
		Stmt  Stmt
	}

	ExprStmt struct {
		X Expr
	}

	FuncStmt struct {
		Memo     token.Pos
		Func     token.Pos
		Name     *Ident
		Params   *FuncParamList
		Body     []Stmt
		BlockEnd token.Pos
	}

	IncDecStmt struct {
		X      Expr
		TokPos token.Pos
		Tok    token.Token
	}

	ConstStmt struct {
		ConstPos token.Pos
		Ident    *Ident
		TokPos   token.Pos
		Value    Expr
	}

	LetStmt struct {
		LetPos token.Pos
		Ident  *Ident
		TokPos token.Pos
		Value  Expr
	}

	AssignStmt struct {
		Lhs    Expr
		TokPos token.Pos
		Tok    token.Token
		Rhs    Expr
	}

	ReturnStmt struct {
		Return token.Pos
		Result Expr
	}

	BranchStmt struct {
		TokPos token.Pos
		Tok    token.Token
		Label  *Ident
	}

	IfStmt struct {
		If       token.Pos
		Cond     Expr
		Body     []Stmt
		Else     []Stmt
		BlockEnd token.Pos
	}

	ForStmt struct {
		For      token.Pos
		Var      *Ident
		Iterable Expr
		Body     []Stmt
		BlockEnd token.Pos
	}

	RepeatWhileStmt struct {
		Repeat   token.Pos
		While    token.Pos
		Cond     Expr
		Body     []Stmt
		BlockEnd token.Pos
	}

	RepeatForeverStmt struct {
		Repeat   token.Pos
		Forever  token.Pos
		Body     []Stmt
		BlockEnd token.Pos
	}

	RepeatCountStmt struct {
		Repeat   token.Pos
		Count    Expr
		Times    token.Pos
		Body     []Stmt
		BlockEnd token.Pos
	}
)

func (s *BadStmt) Pos() token.Pos           { return s.From }
func (s *EmptyStmt) Pos() token.Pos         { return 9999999 }
func (s *LabeledStmt) Pos() token.Pos       { return s.Label.Pos() }
func (s *ExprStmt) Pos() token.Pos          { return s.X.Pos() }
func (s *FuncStmt) Pos() token.Pos          { return s.Func }
func (s *IncDecStmt) Pos() token.Pos        { return s.X.Pos() }
func (s *ConstStmt) Pos() token.Pos         { return s.ConstPos }
func (s *LetStmt) Pos() token.Pos           { return s.LetPos }
func (s *AssignStmt) Pos() token.Pos        { return s.Lhs.Pos() }
func (s *ReturnStmt) Pos() token.Pos        { return s.Return }
func (s *BranchStmt) Pos() token.Pos        { return s.TokPos }
func (s *IfStmt) Pos() token.Pos            { return s.If }
func (s *ForStmt) Pos() token.Pos           { return s.For }
func (s *RepeatWhileStmt) Pos() token.Pos   { return s.Repeat }
func (s *RepeatForeverStmt) Pos() token.Pos { return s.Repeat }
func (s *RepeatCountStmt) Pos() token.Pos   { return s.Repeat }

func (s *BadStmt) End() token.Pos     { return s.To }
func (s *EmptyStmt) End() token.Pos   { return s.Pos() }
func (s *LabeledStmt) End() token.Pos { return s.Stmt.End() }
func (s *ExprStmt) End() token.Pos    { return s.X.End() }
func (s *FuncStmt) End() token.Pos    { return s.BlockEnd + 3 }
func (s *IncDecStmt) End() token.Pos  { return s.TokPos + 2 /* len("++") */ }
func (s *ConstStmt) End() token.Pos   { return s.Value.End() }
func (s *LetStmt) End() token.Pos {
	if s.Value != nil {
		return s.Value.End()
	}
	return s.Ident.End()
}
func (s *AssignStmt) End() token.Pos { return s.Rhs.End() }
func (s *ReturnStmt) End() token.Pos {
	if s.Result != nil {
		return s.Result.End()
	}
	return s.Return + 6 // len("return")
}
func (s *BranchStmt) End() token.Pos {
	if s.Label != nil {
		return s.Label.End()
	}
	return token.Pos(int(s.TokPos) + len(s.Tok.String()))
}
func (s *IfStmt) End() token.Pos {
	if s.Else != nil {
		return s.Else[len(s.Else)-1].End()
	}
	return s.BlockEnd
}
func (s *ForStmt) End() token.Pos           { return s.BlockEnd + 3 }
func (s *RepeatWhileStmt) End() token.Pos   { return s.BlockEnd + 3 }
func (s *RepeatForeverStmt) End() token.Pos { return s.BlockEnd + 3 }
func (s *RepeatCountStmt) End() token.Pos   { return s.BlockEnd + 3 }

func (*BadStmt) stmtNode()           {}
func (*EmptyStmt) stmtNode()         {}
func (*LabeledStmt) stmtNode()       {}
func (*ExprStmt) stmtNode()          {}
func (*FuncStmt) stmtNode()          {}
func (*IncDecStmt) stmtNode()        {}
func (*ConstStmt) stmtNode()         {}
func (*LetStmt) stmtNode()           {}
func (*AssignStmt) stmtNode()        {}
func (*ReturnStmt) stmtNode()        {}
func (*BranchStmt) stmtNode()        {}
func (*IfStmt) stmtNode()            {}
func (*ForStmt) stmtNode()           {}
func (*RepeatWhileStmt) stmtNode()   {}
func (*RepeatForeverStmt) stmtNode() {}
func (*RepeatCountStmt) stmtNode()   {}

type (
	Spec interface {
		Node
		specNode()
	}

	ValueSpec struct {
		Names  []*Ident
		Values []Expr
	}
)

func (s *ValueSpec) Pos() token.Pos { return s.Names[0].Pos() }
func (s *ValueSpec) End() token.Pos {
	if n := len(s.Values); n > 0 {
		return s.Values[n-1].End()
	}
	return s.Names[len(s.Names)-1].End()
}
func (*ValueSpec) specNode() {}

type File struct {
	Size       int
	Stmts      []Stmt
	Unresolved []*Ident
}
