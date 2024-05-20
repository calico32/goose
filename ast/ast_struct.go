package ast

import "github.com/calico32/goose/token"

type (
	StructStmt struct {
		Struct token.Pos
		Name   *Ident
		Fields *StructFieldList
		Init   *StructInit
	}

	PropertyExpr struct {
		Hash  token.Pos
		Ident *Ident
	}

	BracketPropertyExpr struct {
		HashLBracket token.Pos
		X            Expr
		RBracket     token.Pos
	}

	OperatorStmt struct {
		Async     token.Pos
		Memo      token.Pos
		Operator  token.Pos
		Receiver  *Ident
		TokPos    token.Pos
		Tok       token.Token
		Params    *FuncParamList
		Arrow     token.Pos
		ArrowExpr Expr
		Body      []Stmt
		BlockEnd  token.Pos
	}
)

type StructFieldList struct {
	Opening token.Pos
	List    []*StructField
	Closing token.Pos
}

type StructField struct {
	Ident *Ident
	Value Expr
}

type StructInit struct {
	Init     token.Pos
	Body     []Stmt
	BlockEnd token.Pos
}

func (x *StructStmt) Pos() token.Pos { return x.Struct }
func (x *OperatorStmt) Pos() token.Pos {
	if x.Async.IsValid() {
		return x.Async
	}
	return x.Operator
}
func (x *PropertyExpr) Pos() token.Pos        { return x.Hash }
func (x *BracketPropertyExpr) Pos() token.Pos { return x.HashLBracket }

func (x *StructStmt) End() token.Pos {
	if x.Init != nil {
		return x.Init.End()
	}
	return x.Fields.End()
}
func (x *OperatorStmt) End() token.Pos { return x.BlockEnd + 3 }
func (x *PropertyExpr) End() token.Pos { return x.Ident.End() }
func (x *BracketPropertyExpr) End() token.Pos {
	return x.RBracket + 1
}

func (x *StructStmt) stmtNode()          {}
func (x *OperatorStmt) stmtNode()        {}
func (x *PropertyExpr) exprNode()        {}
func (x *BracketPropertyExpr) exprNode() {}

func (x *StructInit) Pos() token.Pos { return x.Init }
func (x *StructInit) End() token.Pos { return x.BlockEnd + 3 }

func (f *StructFieldList) Pos() token.Pos {
	if f.Opening.IsValid() {
		return f.Opening
	}
	if len(f.List) > 0 {
		return f.List[0].Pos()
	}
	return token.NoPos
}

func (f *StructFieldList) End() token.Pos {
	if f.Closing.IsValid() {
		return f.Closing + 1
	}
	if n := len(f.List); n > 0 {
		return f.List[n-1].End()
	}
	return token.NoPos
}

func (f *StructFieldList) NumFields() int { return len(f.List) }

func (f *StructField) Pos() token.Pos { return f.Ident.Pos() }
func (f *StructField) End() token.Pos { return f.Ident.End() }

func (x *StructStmt) Flatten() []Node {
	nodes := make([]Node, 0, len(x.Fields.List))
	for _, field := range x.Fields.List {
		nodes = append(nodes, field.Value.Flatten()...)
	}
	if x.Init != nil {
		nodes = append(nodes, x.Init.Flatten()...)
	}
	return nodes
}

func (x *OperatorStmt) Flatten() []Node {
	nodes := make([]Node, 0, len(x.Params.List)+len(x.Body)+1)

	if x.Receiver != nil {
		nodes = append(nodes, x.Receiver)
	}
	if x.Params != nil {
		nodes = append(nodes, x.Params.Flatten()...)
	}
	for _, stmt := range x.Body {
		nodes = append(nodes, stmt.Flatten()...)
	}
	return nodes
}

func (x *PropertyExpr) Flatten() []Node        { return nil }
func (x *BracketPropertyExpr) Flatten() []Node { return x.X.Flatten() }
func (x *StructInit) Flatten() []Node {
	nodes := make([]Node, 0, len(x.Body))
	for _, stmt := range x.Body {
		nodes = append(nodes, stmt.Flatten()...)
	}
	return nodes
}
func (f *StructField) Flatten() []Node {
	if f.Value != nil {
		return f.Value.Flatten()
	}
	return nil
}
