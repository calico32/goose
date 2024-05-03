package ast

import "github.com/calico32/goose/token"

type (
	FuncExpr struct {
		Async    token.Pos
		Memo     token.Pos
		Func     token.Pos
		Receiver *Ident
		Name     *Ident
		Params   *FuncParamList

		Arrow     token.Pos
		ArrowExpr Expr

		Body     []Stmt
		BlockEnd token.Pos
	}

	ReturnStmt struct {
		Return token.Pos
		Result Expr
	}
)

type FuncParamList struct {
	Opening token.Pos
	List    []*FuncParam
	Closing token.Pos
}

type FuncParam struct {
	Ellipsis token.Pos
	Ident    *Ident
	Value    Expr
}

func (s *FuncExpr) Pos() token.Pos {
	if s.Async.IsValid() {
		return s.Async
	}
	if s.Memo.IsValid() {
		return s.Memo
	}
	return s.Func
}
func (s *ReturnStmt) Pos() token.Pos { return s.Return }

func (s *FuncExpr) End() token.Pos {
	if s.BlockEnd.IsValid() {
		return s.BlockEnd + 1
	}
	return s.ArrowExpr.End()
}
func (s *ReturnStmt) End() token.Pos {
	if s.Result != nil {
		return s.Result.End()
	}
	return s.Return + 6 // len("return")
}

func (*FuncExpr) exprNode() {}

func (*ReturnStmt) stmtNode() {}
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

func (f *FuncParam) Pos() token.Pos {
	if f.Ellipsis.IsValid() {
		return f.Ellipsis
	}
	return f.Ident.Pos()
}
func (f *FuncParam) End() token.Pos { return f.Ident.End() }

func (f *FuncExpr) Flatten() []Node {
	nodes := make([]Node, 0, len(f.Body))
	if f.Receiver != nil {
		nodes = append(nodes, f.Receiver.Flatten()...)
	}

	if f.Name != nil {
		nodes = append(nodes, f.Name.Flatten()...)
	}

	if f.Params != nil {
		nodes = append(nodes, f.Params.Flatten()...)
	}

	if f.ArrowExpr != nil {
		nodes = append(nodes, f.ArrowExpr.Flatten()...)
	}

	if len(f.Body) > 0 {
		for _, stmt := range f.Body {
			nodes = append(nodes, stmt.Flatten()...)
		}
	}
	return nodes
}

func (f *ReturnStmt) Flatten() []Node { return f.Result.Flatten() }

func (f *FuncParamList) Flatten() []Node {
	nodes := make([]Node, 0, len(f.List))
	for _, param := range f.List {
		nodes = append(nodes, param.Flatten()...)
	}
	return nodes
}

func (f *FuncParam) Flatten() []Node {
	nodes := make([]Node, 0, 2)
	// nodes = append(nodes, f.Ident.Flatten()...)
	if f.Value != nil {
		nodes = append(nodes, f.Value.Flatten()...)
	}
	return nodes
}
