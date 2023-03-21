package ast

import "github.com/calico32/goose/token"

type (
	ForStmt struct {
		For      token.Pos
		Await    token.Pos
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

	BranchStmt struct {
		TokPos token.Pos
		Tok    token.Token
		Label  *Ident
	}
)

func (s *ForStmt) Pos() token.Pos           { return s.For }
func (s *RepeatWhileStmt) Pos() token.Pos   { return s.Repeat }
func (s *RepeatForeverStmt) Pos() token.Pos { return s.Repeat }
func (s *RepeatCountStmt) Pos() token.Pos   { return s.Repeat }
func (s *BranchStmt) Pos() token.Pos        { return s.TokPos }

func (s *ForStmt) End() token.Pos           { return s.BlockEnd + 3 }
func (s *RepeatWhileStmt) End() token.Pos   { return s.BlockEnd + 3 }
func (s *RepeatForeverStmt) End() token.Pos { return s.BlockEnd + 3 }
func (s *RepeatCountStmt) End() token.Pos   { return s.BlockEnd + 3 }
func (s *BranchStmt) End() token.Pos {
	if s.Label != nil {
		return s.Label.End()
	}
	return token.Pos(int(s.TokPos) + len(s.Tok.String()))
}

func (*ForStmt) stmtNode()           {}
func (*RepeatWhileStmt) stmtNode()   {}
func (*RepeatForeverStmt) stmtNode() {}
func (*RepeatCountStmt) stmtNode()   {}
func (*BranchStmt) stmtNode()        {}
