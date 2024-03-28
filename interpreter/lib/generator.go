package lib

type GeneratorMessage interface {
	generatorMessage()
}

type (
	GeneratorReturn struct{ Value Value }
	GeneratorYield  struct{ Value Value }
	GeneratorError  struct{ Error error }
	GeneratorNext   struct{ Value Value }
	GeneratorIsDone struct{ Done bool }
	GeneratorDone   struct{}
)

func (*GeneratorReturn) generatorMessage() {}
func (*GeneratorYield) generatorMessage()  {}
func (*GeneratorError) generatorMessage()  {}
func (*GeneratorNext) generatorMessage()   {}
func (*GeneratorIsDone) generatorMessage() {}
func (*GeneratorDone) generatorMessage()   {}
