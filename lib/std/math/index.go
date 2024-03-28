package std_math

import (
	"math"
	"math/big"

	. "github.com/calico32/goose/interpreter/lib"
	"github.com/calico32/goose/lib/types"
)

var Doc = types.StdlibDoc{
	Name:        "math",
	Description: "Mathematical functions and constants.",
}

var Index = map[string]Value{
	"F/sin": &Func{Executor: func(ctx *FuncContext) *Return {
		if len(ctx.Args) < 1 {
			ctx.Interp.Throw("std/math:sin(x): expected 1 argument")
			return &Return{}
		}
		x := ctx.Args[0]
		if numeric, ok := x.(Numeric); ok {
			return NewReturn(math.Sin(numeric.Float64()))
		} else {
			ctx.Interp.Throw("std/math:sin(x): expected number")
			return &Return{}
		}
	}},
	"F/cos": &Func{Executor: func(ctx *FuncContext) *Return {
		if len(ctx.Args) < 1 {
			ctx.Interp.Throw("std/math:cos(x): expected 1 argument")
			return &Return{}
		}
		x := ctx.Args[0]
		if numeric, ok := x.(Numeric); ok {
			return NewReturn(math.Cos(numeric.Float64()))
		} else {
			ctx.Interp.Throw("std/math:cos(x): expected number")
			return &Return{}
		}
	}},
	"F/tan": &Func{Executor: func(ctx *FuncContext) *Return {
		if len(ctx.Args) < 1 {
			ctx.Interp.Throw("std/math:tan(x): expected 1 argument")
			return &Return{}
		}
		x := ctx.Args[0]
		if numeric, ok := x.(Numeric); ok {
			return NewReturn(math.Tan(numeric.Float64()))
		} else {
			ctx.Interp.Throw("std/math:tan(x): expected number")
			return &Return{}
		}
	}},
	"F/asin": &Func{Executor: func(ctx *FuncContext) *Return {
		if len(ctx.Args) < 1 {
			ctx.Interp.Throw("std/math:asin(x): expected 1 argument")
			return &Return{}
		}
		x := ctx.Args[0]
		if numeric, ok := x.(Numeric); ok {
			return NewReturn(math.Asin(numeric.Float64()))
		} else {
			ctx.Interp.Throw("std/math:asin(x): expected number")
			return &Return{}
		}
	}},
	"F/acos": &Func{Executor: func(ctx *FuncContext) *Return {
		if len(ctx.Args) < 1 {
			ctx.Interp.Throw("std/math:acos(x): expected 1 argument")
			return &Return{}
		}
		x := ctx.Args[0]
		if numeric, ok := x.(Numeric); ok {
			return NewReturn(math.Acos(numeric.Float64()))
		} else {
			ctx.Interp.Throw("std/math:acos(x): expected number")
			return &Return{}
		}
	}},
	"F/atan": &Func{Executor: func(ctx *FuncContext) *Return {
		if len(ctx.Args) < 1 {
			ctx.Interp.Throw("std/math:atan(x): expected 1 argument")
			return &Return{}
		}
		x := ctx.Args[0]
		if numeric, ok := x.(Numeric); ok {
			return NewReturn(math.Atan(numeric.Float64()))
		} else {
			ctx.Interp.Throw("std/math:atan(x): expected number")
			return &Return{}
		}
	}},
	"F/atan2": &Func{Executor: func(ctx *FuncContext) *Return {
		if len(ctx.Args) < 2 {
			ctx.Interp.Throw("std/math:atan2(y, x): expected 2 arguments")
			return &Return{}
		}
		y := ctx.Args[0]
		x := ctx.Args[1]
		if numericY, ok := y.(Numeric); ok {
			if numericX, ok := x.(Numeric); ok {
				return NewReturn(math.Atan2(numericY.Float64(), numericX.Float64()))
			} else {
				ctx.Interp.Throw("std/math:atan2(y, x): expected number")
				return &Return{}
			}
		} else {
			ctx.Interp.Throw("std/math:atan2(y, x): expected number")
			return &Return{}
		}
	}},
	"F/sinh": &Func{Executor: func(ctx *FuncContext) *Return {
		if len(ctx.Args) < 1 {
			ctx.Interp.Throw("std/math:sinh(x): expected 1 argument")
			return &Return{}
		}
		x := ctx.Args[0]
		if numeric, ok := x.(Numeric); ok {
			return NewReturn(math.Sinh(numeric.Float64()))
		} else {
			ctx.Interp.Throw("std/math:sinh(x): expected number")
			return &Return{}
		}
	}},
	"F/cosh": &Func{Executor: func(ctx *FuncContext) *Return {
		if len(ctx.Args) < 1 {
			ctx.Interp.Throw("std/math:cosh(x): expected 1 argument")
			return &Return{}
		}
		x := ctx.Args[0]
		if numeric, ok := x.(Numeric); ok {
			return NewReturn(math.Cosh(numeric.Float64()))
		} else {
			ctx.Interp.Throw("std/math:cosh(x): expected number")
			return &Return{}
		}
	}},
	"F/tanh": &Func{Executor: func(ctx *FuncContext) *Return {
		if len(ctx.Args) < 1 {
			ctx.Interp.Throw("std/math:tanh(x): expected 1 argument")
			return &Return{}
		}
		x := ctx.Args[0]
		if numeric, ok := x.(Numeric); ok {
			return NewReturn(math.Tanh(numeric.Float64()))
		} else {
			ctx.Interp.Throw("std/math:tanh(x): expected number")
			return &Return{}
		}
	}},
	"F/round": &Func{Executor: func(ctx *FuncContext) *Return {
		if len(ctx.Args) < 1 {
			ctx.Interp.Throw("std/math:round(x): expected 1 argument")
			return &Return{}
		}
		x := ctx.Args[0]
		if numeric, ok := x.(Numeric); ok {
			return NewReturn(math.Round(numeric.Float64()))
		} else {
			ctx.Interp.Throw("std/math:round(x): expected number")
			return &Return{}
		}
	}},
	"F/floor": &Func{Executor: func(ctx *FuncContext) *Return {
		if len(ctx.Args) < 1 {
			ctx.Interp.Throw("std/math:floor(x): expected 1 argument")
			return &Return{}
		}
		x := ctx.Args[0]
		if numeric, ok := x.(Numeric); ok {
			return NewReturn(math.Floor(numeric.Float64()))
		} else {
			ctx.Interp.Throw("std/math:floor(x): expected number")
			return &Return{}
		}
	}},
	"F/ceil": &Func{Executor: func(ctx *FuncContext) *Return {
		if len(ctx.Args) < 1 {
			ctx.Interp.Throw("std/math:ceil(x): expected 1 argument")
			return &Return{}
		}
		x := ctx.Args[0]
		if numeric, ok := x.(Numeric); ok {
			return NewReturn(math.Ceil(numeric.Float64()))
		} else {
			ctx.Interp.Throw("std/math:ceil(x): expected number")
			return &Return{}
		}
	}},
	"F/sqrt": &Func{Executor: func(ctx *FuncContext) *Return {
		if len(ctx.Args) < 1 {
			ctx.Interp.Throw("std/math:sqrt(x): expected 1 argument")
			return &Return{}
		}
		x := ctx.Args[0]
		if numeric, ok := x.(Numeric); ok {
			return NewReturn(math.Sqrt(numeric.Float64()))
		} else {
			ctx.Interp.Throw("std/math:sqrt(x): expected number")
			return &Return{}
		}
	}},
	"F/log": &Func{Executor: func(ctx *FuncContext) *Return {
		if len(ctx.Args) < 1 {
			ctx.Interp.Throw("std/math:log(x): expected 1 argument")
			return &Return{}
		}
		x := ctx.Args[0]
		if numeric, ok := x.(Numeric); ok {
			return NewReturn(math.Log(numeric.Float64()))
		} else {
			ctx.Interp.Throw("std/math:log(x): expected number")
			return &Return{}
		}
	}},
	"F/log2": &Func{Executor: func(ctx *FuncContext) *Return {
		if len(ctx.Args) < 1 {
			ctx.Interp.Throw("std/math:log2(x): expected 1 argument")
			return &Return{}
		}
		x := ctx.Args[0]
		if numeric, ok := x.(Numeric); ok {
			return NewReturn(math.Log2(numeric.Float64()))
		} else {
			ctx.Interp.Throw("std/math:log2(x): expected number")
			return &Return{}
		}
	}},
	"F/log10": &Func{Executor: func(ctx *FuncContext) *Return {
		if len(ctx.Args) < 1 {
			ctx.Interp.Throw("std/math:log10(x): expected 1 argument")
			return &Return{}
		}
		x := ctx.Args[0]
		if numeric, ok := x.(Numeric); ok {
			return NewReturn(math.Log10(numeric.Float64()))
		} else {
			ctx.Interp.Throw("std/math:log10(x): expected number")
			return &Return{}
		}
	}},
	"F/exp": &Func{Executor: func(ctx *FuncContext) *Return {
		if len(ctx.Args) < 1 {
			ctx.Interp.Throw("std/math:exp(x): expected 1 argument")
			return &Return{}
		}
		x := ctx.Args[0]
		if numeric, ok := x.(Numeric); ok {
			return NewReturn(math.Exp(numeric.Float64()))
		} else {
			ctx.Interp.Throw("std/math:exp(x): expected number")
			return &Return{}
		}
	}},
	"F/pow": &Func{Executor: func(ctx *FuncContext) *Return {
		if len(ctx.Args) < 2 {
			ctx.Interp.Throw("std/math:pow(x, y): expected 2 arguments")
			return &Return{}
		}
		x := ctx.Args[0]
		y := ctx.Args[1]
		if numericX, ok := x.(Numeric); ok {
			if numericY, ok := y.(Numeric); ok {
				return NewReturn(math.Pow(numericX.Float64(), numericY.Float64()))
			} else {
				ctx.Interp.Throw("std/math:pow(x, y): expected number")
				return &Return{}
			}
		} else {
			ctx.Interp.Throw("std/math:pow(x, y): expected number")
			return &Return{}
		}
	}},
	"F/abs": &Func{Executor: func(ctx *FuncContext) *Return {
		if len(ctx.Args) < 1 {
			ctx.Interp.Throw("std/math:abs(x): expected 1 argument")
			return &Return{}
		}
		x := ctx.Args[0]
		if numeric, ok := x.(Numeric); ok {
			return NewReturn(math.Abs(numeric.Float64()))
		} else {
			ctx.Interp.Throw("std/math:abs(x): expected number")
			return &Return{}
		}
	}},
	"F/sign": &Func{Executor: func(ctx *FuncContext) *Return {
		if len(ctx.Args) < 1 {
			ctx.Interp.Throw("std/math:sign(x): expected 1 argument")
			return &Return{}
		}
		x := ctx.Args[0]
		if numeric, ok := x.(Numeric); ok {
			if numeric.Float64() < 0 {
				return NewReturn(NewFloat(-1))
			} else if numeric.Float64() > 0 {
				return NewReturn(NewFloat(1))
			} else {
				return NewReturn(NewFloat(0))
			}
		} else {
			ctx.Interp.Throw("std/math:sign(x): expected number")
			return &Return{}
		}
	}},
	"F/min": &Func{Executor: func(ctx *FuncContext) *Return {
		if len(ctx.Args) < 1 {
			ctx.Interp.Throw("std/math:min(nums): expected at least 1 argument")
			return &Return{}
		}

		var min float64
		for i, arg := range ctx.Args {
			if numeric, ok := arg.(Numeric); ok {
				if i == 0 {
					min = numeric.Float64()
				} else {
					min = math.Min(min, numeric.Float64())
				}
			} else {
				ctx.Interp.Throw("std/math:min(nums): expected number")
				return &Return{}
			}
		}

		return NewReturn(min)
	}},
	"F/max": &Func{Executor: func(ctx *FuncContext) *Return {
		if len(ctx.Args) < 1 {
			ctx.Interp.Throw("std/math:max(nums): expected at least 1 argument")
			return &Return{}
		}

		var max float64
		for i, arg := range ctx.Args {
			if numeric, ok := arg.(Numeric); ok {
				if i == 0 {
					max = numeric.Float64()
				} else {
					max = math.Max(max, numeric.Float64())
				}
			} else {
				ctx.Interp.Throw("std/math:max(nums): expected number")
				return &Return{}
			}
		}

		return NewReturn(max)
	}},
	"F/clamp": &Func{Executor: func(ctx *FuncContext) *Return {
		if len(ctx.Args) < 3 {
			ctx.Interp.Throw("std/math:clamp(x, min, max): expected 3 arguments")
			return &Return{}
		}
		x := ctx.Args[0]
		min := ctx.Args[1]
		max := ctx.Args[2]
		if numericX, ok := x.(Numeric); ok {
			if numericMin, ok := min.(Numeric); ok {
				if numericMax, ok := max.(Numeric); ok {
					return NewReturn(math.Min(math.Max(numericX.Float64(), numericMin.Float64()), numericMax.Float64()))
				} else {
					ctx.Interp.Throw("std/math:clamp(x, min, max): expected number")
					return &Return{}
				}
			} else {
				ctx.Interp.Throw("std/math:clamp(x, min, max): expected number")
				return &Return{}
			}
		} else {
			ctx.Interp.Throw("std/math:clamp(x, min, max): expected number")
			return &Return{}
		}
	}},
	"F/parseInt": &Func{Executor: func(ctx *FuncContext) *Return {
		if len(ctx.Args) < 1 {
			ctx.Interp.Throw("std/math:parseInt(s): expected 1 argument")
			return &Return{}
		}
		s := ctx.Args[0]
		base := 10
		if len(ctx.Args) > 1 {
			if integer, ok := ctx.Args[1].(*Integer); ok {
				base = int(integer.Value.Int64())
			} else {
				ctx.Interp.Throw("std/math:parseInt(s, base): expected integer")
			}
		}
		if str, ok := s.(*String); ok {
			i := new(big.Int)
			_, ok := i.SetString(str.Value, base)
			if !ok {
				ctx.Interp.Throw("std/math:parseInt(s): failed to parse integer")
				return &Return{}
			}
			return NewReturn(NewInteger(i))
		} else {
			ctx.Interp.Throw("std/math:parseInt(s): expected string")
			return &Return{}
		}
	}},

	"C/pi":     NewFloat(math.Pi),
	"C/e":      NewFloat(math.E),
	"C/phi":    NewFloat(math.Phi),
	"C/sqrt2":  NewFloat(math.Sqrt2),
	"C/ln2":    NewFloat(math.Ln2),
	"C/log2e":  NewFloat(math.Log2E),
	"C/ln10":   NewFloat(math.Ln10),
	"C/log10e": NewFloat(math.Log10E),
}
