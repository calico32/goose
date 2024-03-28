package std_json

import (
	"encoding/json"

	. "github.com/calico32/goose/interpreter/lib"
	"github.com/calico32/goose/lib/types"
)

var Doc = types.StdlibDoc{
	Name:        "json",
	Description: "Encode and decode JSON data as strings.",
}

var Index = map[string]Value{
	"F/decode": &Func{Executor: func(ctx *FuncContext) *Return {
		if len(ctx.Args) < 1 {
			ctx.Interp.Throw("json.decode(s): expected at least 1 argument")
			return &Return{}
		}

		s := ctx.Args[0]
		if str, ok := s.(*String); ok {
			var v any
			err := json.Unmarshal([]byte(str.Value), &v)
			if err != nil {
				ctx.Interp.Throw("json.decode(s): " + err.Error())
				return &Return{}
			}
			wrapped := Wrap(v)
			return NewReturn(&wrapped)
		} else {
			ctx.Interp.Throw("json.decode(s): expected string")
			return &Return{}
		}
	}},
	"F/encode": &Func{Executor: func(ctx *FuncContext) *Return {
		if len(ctx.Args) < 1 {
			ctx.Interp.Throw("json.encode(v): expected at least 1 argument")
			return &Return{}
		}

		v := ctx.Args[0]
		data, err := json.Marshal(v.Unwrap())
		if err != nil {
			ctx.Interp.Throw("json.encode(v): " + err.Error())
			return &Return{}
		}
		return NewReturn(NewString(string(data)))
	}},
}
