package ast

import "github.com/calico32/goose/token"

type (
	GeneratorExpr struct {
		Async     token.Pos
		Generator token.Pos
		Receiver  *Ident
		Name      *Ident
		Params    *FuncParamList
		Body      []Stmt
		BlockEnd  token.Pos
	}

	YieldStmt struct {
		Yield  token.Pos
		Result Expr
	}

	RangeExpr struct {
		Start   Expr
		ToPos   token.Pos
		Stop    Expr
		StepPos token.Pos
		Step    Expr
	}
)

func (s *GeneratorExpr) Pos() token.Pos {
	if s.Async.IsValid() {
		return s.Async
	}
	return s.Generator
}
func (s *YieldStmt) Pos() token.Pos { return s.Yield }
func (x *RangeExpr) Pos() token.Pos { return x.Start.Pos() }

func (s *GeneratorExpr) End() token.Pos { return s.BlockEnd + 3 }
func (s *YieldStmt) End() token.Pos {
	if s.Result != nil {
		return s.Result.End()
	}
	return s.Yield + 5 // len("yield")
}
func (x *RangeExpr) End() token.Pos {
	if x.Step != nil {
		return x.Step.End()
	}
	return x.Stop.End()
}
func (*GeneratorExpr) exprNode() {}
func (*YieldStmt) stmtNode()     {}
func (*RangeExpr) exprNode()     {}

func (s *GeneratorExpr) Flatten() []Node {
	nodes := make([]Node, 0, len(s.Body))
	if s.Receiver != nil {
		nodes = append(nodes, s.Receiver.Flatten()...)
	}
	if s.Name != nil {
		nodes = append(nodes, s.Name.Flatten()...)
	}
	if s.Params != nil {
		nodes = append(nodes, s.Params.Flatten()...)
	}
	for _, stmt := range s.Body {
		nodes = append(nodes, stmt.Flatten()...)
	}
	return nodes
}

func (s *YieldStmt) Flatten() []Node { return s.Result.Flatten() }

func (x *RangeExpr) Flatten() []Node {
	nodes := make([]Node, 0, 3)
	nodes = append(nodes, x.Start.Flatten()...)
	nodes = append(nodes, x.Stop.Flatten()...)
	if x.Step != nil {
		nodes = append(nodes, x.Step.Flatten()...)
	}
	return nodes
}
