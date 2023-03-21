package ast

import "github.com/calico32/goose/token"

type (
	CompositeLiteral struct {
		Lbrace token.Pos
		Fields []*CompositeField
		Rbrace token.Pos
	}
)

type CompositeField struct {
	Key   Expr
	Value Expr
}

func (x *CompositeLiteral) Pos() token.Pos { return x.Lbrace }
func (x *CompositeLiteral) End() token.Pos { return x.Rbrace + 1 }

func (*CompositeLiteral) exprNode() {}
