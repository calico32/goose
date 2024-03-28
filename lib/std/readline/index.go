package std_readline

import (
	"bufio"

	. "github.com/calico32/goose/interpreter/lib"
	"github.com/calico32/goose/lib/types"
)

var Doc = types.StdlibDoc{
	Name:        "readline",
	Description: "Read text from standard input.",
}

var Index = map[string]Value{
	"F/read": &Func{Executor: func(ctx *FuncContext) *Return {
		r := bufio.NewReader(ctx.Interp.Stdin())
		line, err := r.ReadString('\n')

		if err != nil {
			ctx.Interp.Throw(err.Error())
		}

		return NewReturn(NewString(line[:len(line)-1]))
	}},
}
