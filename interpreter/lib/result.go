package lib

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
	LoneValue struct {
		Value Value
	}
)

func (*Return) stmtResult()    {}
func (*Yield) stmtResult()     {}
func (*Break) stmtResult()     {}
func (*Continue) stmtResult()  {}
func (*Void) stmtResult()      {}
func (*Export) stmtResult()    {}
func (*Decl) stmtResult()      {} // TODO: module variable update leads to update in export?
func (*LoneValue) stmtResult() {}

func NewReturn[T ValueType](value T) *Return     { return &Return{Wrap(value)} }
func NewYield[T ValueType](value T) *Yield       { return &Yield{Wrap(value)} }
func NewBreak() *Break                           { return &Break{} }
func NewContinue() *Continue                     { return &Continue{} }
func NewVoid() *Void                             { return &Void{} }
func NewExport(value Value, name string) *Export { return &Export{value, name} }
func NewDecl(value Value, name string) *Decl     { return &Decl{value, name} }
func NewLoneValue(value Value) *LoneValue        { return &LoneValue{value} }

var ReturnNull = NewReturn(NullValue)
var YieldNull = NewYield(NullValue)
var ReturnTrue = NewReturn(TrueValue)
var ReturnFalse = NewReturn(FalseValue)
