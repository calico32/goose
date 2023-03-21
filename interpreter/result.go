package interpreter

type (
	StmtResult interface {
		stmtResult()
	}

	Return   struct{ Value Value }
	Yield    struct{ Value Value }
	Break    struct{}
	Continue struct{}
	Void     struct{}
	Export   struct {
		Value Value
		Name  string
	}
	Decl struct {
		Value Value
		Name  string
	}
)

func (*Return) stmtResult()   {}
func (*Yield) stmtResult()    {}
func (*Break) stmtResult()    {}
func (*Continue) stmtResult() {}
func (*Void) stmtResult()     {}
func (*Export) stmtResult()   {}
func (*Decl) stmtResult()     {} // TODO: module variable update leads to update in export?
