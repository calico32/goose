package ast

import "github.com/calico32/goose/token"

type (
	NativeStmt interface {
		Stmt
		nativeStmt()
	}

	NativeConst struct {
		Native token.Pos
		Const  token.Pos
		Ident  *Ident
	}

	NativeStruct struct {
		Native token.Pos
		Struct token.Pos
		Name   *Ident
		Fields *StructFieldList
	}

	NativeFunc struct {
		Native   token.Pos
		Async    token.Pos
		Memo     token.Pos
		Fn       token.Pos
		Receiver *Ident
		Name     *Ident
		Params   *FuncParamList
	}

	NativeOperator struct {
		Native   token.Pos
		Async    token.Pos
		Operator token.Pos
		Receiver *Ident
		TokPos   token.Pos
		Tok      token.Token
		Params   *FuncParamList
	}

	NativeExpr struct {
		Native token.Pos
		Id     string
	}
)

func (s *NativeConst) Pos() token.Pos    { return s.Native }
func (s *NativeStruct) Pos() token.Pos   { return s.Native }
func (s *NativeFunc) Pos() token.Pos     { return s.Native }
func (s *NativeOperator) Pos() token.Pos { return s.Native }
func (s *NativeExpr) Pos() token.Pos     { return s.Native }

func (s *NativeConst) End() token.Pos    { return s.Ident.End() }
func (s *NativeStruct) End() token.Pos   { return s.Fields.End() }
func (s *NativeFunc) End() token.Pos     { return s.Params.End() }
func (s *NativeOperator) End() token.Pos { return s.Params.End() }
func (s *NativeExpr) End() token.Pos     { return s.Native + token.Pos(len(s.Id)) }

func (*NativeConst) stmtNode()    {}
func (*NativeStruct) stmtNode()   {}
func (*NativeFunc) stmtNode()     {}
func (*NativeOperator) stmtNode() {}
func (*NativeExpr) exprNode()     {}

func (*NativeConst) nativeStmt()    {}
func (*NativeStruct) nativeStmt()   {}
func (*NativeFunc) nativeStmt()     {}
func (*NativeOperator) nativeStmt() {}

func (s *NativeConst) Flatten() []Node  { return nil }
func (s *NativeStruct) Flatten() []Node { return nil }
func (s *NativeFunc) Flatten() []Node {
	if s.Receiver != nil {
		return s.Receiver.Flatten()
	}
	return nil
}
func (s *NativeOperator) Flatten() []Node {
	if s.Receiver != nil {
		return s.Receiver.Flatten()
	}
	return nil
}
func (s *NativeExpr) Flatten() []Node { return nil }
