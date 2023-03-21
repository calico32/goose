package interpreter

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/calico32/goose/ast"
)

func (i *interp) runNativeStmt(scope *Scope, stmt ast.NativeStmt) StmtResult {
	var name string
	switch stmt := stmt.(type) {
	case *ast.NativeConst:
		name = "C/" + stmt.Ident.Name
	case *ast.NativeStruct:
		name = "S/" + stmt.Name.Name
	case *ast.NativeFunc:
		if stmt.Receiver != nil {
			name = "F/" + stmt.Receiver.Name + "." + name
		} else {
			name = "F/" + stmt.Name.Name
		}
	case *ast.NativeOperator:
		name = "O/" + stmt.Receiver.Name + "." + stmt.Tok.String()
	default:
		i.throw(fmt.Sprintf("invalid native stmt type %T", stmt))
	}

	scheme := "file"
	module := scope.Module()
	moduleName := strings.TrimPrefix(module.Name, i.gooseRoot+"/")
	if moduleName != module.Name {
		scheme = "pkg"
	}
	if moduleNatives, ok := Natives[scheme+":"+moduleName]; ok {
		if value, ok := moduleNatives[name]; ok {
			scope.idents[name[2:]] = &Variable{
				Value:    value,
				Constant: true,
			}
			return &Decl{
				Name:  name,
				Value: value,
			}
		}
	}

	i.throw("native symbol %s not found in module %s:%s", name, scheme, moduleName)
	return nil
}

var Natives = map[string]map[string]Value{
	"pkg:std/fs/_module.goose": {
		"F/readFile": &Func{Executor: func(ctx *FuncContext) *Return {
			if len(ctx.Args) < 1 {
				ctx.interp.throw("std/fs:readFile(file): expected 1 argument")
			}
			file := toString(ctx.interp, ctx.Scope, ctx.Args[0])
			if file == "" {
				ctx.interp.throw("std/fs:readFile(file): expected string")
			}

			f, err := os.ReadFile(file)
			if err != nil {
				ctx.interp.throw(err.Error())
			}

			return &Return{&String{string(f)}}
		}},
		"F/writeFile": &Func{Executor: func(ctx *FuncContext) *Return {
			if len(ctx.Args) < 2 {
				ctx.interp.throw("std/fs:writeFile(file): expected 2 arguments")
			}
			file := toString(ctx.interp, ctx.Scope, ctx.Args[0])
			content := toString(ctx.interp, ctx.Scope, ctx.Args[1])
			if file == "" {
				ctx.interp.throw("std/fs:writeFile(file): expected string")
			}

			err := os.WriteFile(file, []byte(content), 0644)
			if err != nil {
				ctx.interp.throw(err.Error())
			}

			return &Return{}
		}},
		"F/appendFile": &Func{Executor: func(ctx *FuncContext) *Return {
			ctx.interp.throw("TODO: implement std/fs/_module.goose#appendFile")
			return &Return{Value: &String{Value: "TODO: implement std/fs/_module.goose#appendFile"}}
		}},
		"F/deleteFile": &Func{Executor: func(ctx *FuncContext) *Return {
			ctx.interp.throw("TODO: implement std/fs/_module.goose#deleteFile")
			return &Return{Value: &String{Value: "TODO: implement std/fs/_module.goose#deleteFile"}}
		}},
	},
	"pkg:std/math/_module.goose": {
		"F/sin": &Func{Executor: func(ctx *FuncContext) *Return {
			if len(ctx.Args) < 1 {
				ctx.interp.throw("std/math:sin(x): expected 1 argument")
				return &Return{}
			}
			x := ctx.Args[0]
			if numeric, ok := x.(Numeric); ok {
				return &Return{&Float{math.Sin(numeric.Float64())}}
			} else {
				ctx.interp.throw("std/math:sin(x): expected number")
				return &Return{}
			}
		}},
		"F/cos": &Func{Executor: func(ctx *FuncContext) *Return {
			if len(ctx.Args) < 1 {
				ctx.interp.throw("std/math:cos(x): expected 1 argument")
				return &Return{}
			}
			x := ctx.Args[0]
			if numeric, ok := x.(Numeric); ok {
				return &Return{&Float{math.Cos(numeric.Float64())}}
			} else {
				ctx.interp.throw("std/math:cos(x): expected number")
				return &Return{}
			}
		}},
		"F/tan": &Func{Executor: func(ctx *FuncContext) *Return {
			if len(ctx.Args) < 1 {
				ctx.interp.throw("std/math:tan(x): expected 1 argument")
				return &Return{}
			}
			x := ctx.Args[0]
			if numeric, ok := x.(Numeric); ok {
				return &Return{&Float{math.Tan(numeric.Float64())}}
			} else {
				ctx.interp.throw("std/math:tan(x): expected number")
				return &Return{}
			}
		}},
		"F/asin": &Func{Executor: func(ctx *FuncContext) *Return {
			if len(ctx.Args) < 1 {
				ctx.interp.throw("std/math:asin(x): expected 1 argument")
				return &Return{}
			}
			x := ctx.Args[0]
			if numeric, ok := x.(Numeric); ok {
				return &Return{&Float{math.Asin(numeric.Float64())}}
			} else {
				ctx.interp.throw("std/math:asin(x): expected number")
				return &Return{}
			}
		}},
		"F/acos": &Func{Executor: func(ctx *FuncContext) *Return {
			if len(ctx.Args) < 1 {
				ctx.interp.throw("std/math:acos(x): expected 1 argument")
				return &Return{}
			}
			x := ctx.Args[0]
			if numeric, ok := x.(Numeric); ok {
				return &Return{&Float{math.Acos(numeric.Float64())}}
			} else {
				ctx.interp.throw("std/math:acos(x): expected number")
				return &Return{}
			}
		}},
		"F/atan": &Func{Executor: func(ctx *FuncContext) *Return {
			if len(ctx.Args) < 1 {
				ctx.interp.throw("std/math:atan(x): expected 1 argument")
				return &Return{}
			}
			x := ctx.Args[0]
			if numeric, ok := x.(Numeric); ok {
				return &Return{&Float{math.Atan(numeric.Float64())}}
			} else {
				ctx.interp.throw("std/math:atan(x): expected number")
				return &Return{}
			}
		}},
		"F/atan2": &Func{Executor: func(ctx *FuncContext) *Return {
			if len(ctx.Args) < 2 {
				ctx.interp.throw("std/math:atan2(y, x): expected 2 arguments")
				return &Return{}
			}
			y := ctx.Args[0]
			x := ctx.Args[1]
			if numericY, ok := y.(Numeric); ok {
				if numericX, ok := x.(Numeric); ok {
					return &Return{&Float{math.Atan2(numericY.Float64(), numericX.Float64())}}
				} else {
					ctx.interp.throw("std/math:atan2(y, x): expected number")
					return &Return{}
				}
			} else {
				ctx.interp.throw("std/math:atan2(y, x): expected number")
				return &Return{}
			}
		}},
		"F/sinh": &Func{Executor: func(ctx *FuncContext) *Return {
			if len(ctx.Args) < 1 {
				ctx.interp.throw("std/math:sinh(x): expected 1 argument")
				return &Return{}
			}
			x := ctx.Args[0]
			if numeric, ok := x.(Numeric); ok {
				return &Return{&Float{math.Sinh(numeric.Float64())}}
			} else {
				ctx.interp.throw("std/math:sinh(x): expected number")
				return &Return{}
			}
		}},
		"F/cosh": &Func{Executor: func(ctx *FuncContext) *Return {
			if len(ctx.Args) < 1 {
				ctx.interp.throw("std/math:cosh(x): expected 1 argument")
				return &Return{}
			}
			x := ctx.Args[0]
			if numeric, ok := x.(Numeric); ok {
				return &Return{&Float{math.Cosh(numeric.Float64())}}
			} else {
				ctx.interp.throw("std/math:cosh(x): expected number")
				return &Return{}
			}
		}},
		"F/tanh": &Func{Executor: func(ctx *FuncContext) *Return {
			if len(ctx.Args) < 1 {
				ctx.interp.throw("std/math:tanh(x): expected 1 argument")
				return &Return{}
			}
			x := ctx.Args[0]
			if numeric, ok := x.(Numeric); ok {
				return &Return{&Float{math.Tanh(numeric.Float64())}}
			} else {
				ctx.interp.throw("std/math:tanh(x): expected number")
				return &Return{}
			}
		}},
		"F/round": &Func{Executor: func(ctx *FuncContext) *Return {
			if len(ctx.Args) < 1 {
				ctx.interp.throw("std/math:round(x): expected 1 argument")
				return &Return{}
			}
			x := ctx.Args[0]
			if numeric, ok := x.(Numeric); ok {
				return &Return{&Float{math.Round(numeric.Float64())}}
			} else {
				ctx.interp.throw("std/math:round(x): expected number")
				return &Return{}
			}
		}},
		"F/floor": &Func{Executor: func(ctx *FuncContext) *Return {
			if len(ctx.Args) < 1 {
				ctx.interp.throw("std/math:floor(x): expected 1 argument")
				return &Return{}
			}
			x := ctx.Args[0]
			if numeric, ok := x.(Numeric); ok {
				return &Return{&Float{math.Floor(numeric.Float64())}}
			} else {
				ctx.interp.throw("std/math:floor(x): expected number")
				return &Return{}
			}
		}},
		"F/ceil": &Func{Executor: func(ctx *FuncContext) *Return {
			if len(ctx.Args) < 1 {
				ctx.interp.throw("std/math:ceil(x): expected 1 argument")
				return &Return{}
			}
			x := ctx.Args[0]
			if numeric, ok := x.(Numeric); ok {
				return &Return{&Float{math.Ceil(numeric.Float64())}}
			} else {
				ctx.interp.throw("std/math:ceil(x): expected number")
				return &Return{}
			}
		}},
		"F/sqrt": &Func{Executor: func(ctx *FuncContext) *Return {
			if len(ctx.Args) < 1 {
				ctx.interp.throw("std/math:sqrt(x): expected 1 argument")
				return &Return{}
			}
			x := ctx.Args[0]
			if numeric, ok := x.(Numeric); ok {
				return &Return{&Float{math.Sqrt(numeric.Float64())}}
			} else {
				ctx.interp.throw("std/math:sqrt(x): expected number")
				return &Return{}
			}
		}},
		"F/log": &Func{Executor: func(ctx *FuncContext) *Return {
			if len(ctx.Args) < 1 {
				ctx.interp.throw("std/math:log(x): expected 1 argument")
				return &Return{}
			}
			x := ctx.Args[0]
			if numeric, ok := x.(Numeric); ok {
				return &Return{&Float{math.Log(numeric.Float64())}}
			} else {
				ctx.interp.throw("std/math:log(x): expected number")
				return &Return{}
			}
		}},
		"F/log2": &Func{Executor: func(ctx *FuncContext) *Return {
			if len(ctx.Args) < 1 {
				ctx.interp.throw("std/math:log2(x): expected 1 argument")
				return &Return{}
			}
			x := ctx.Args[0]
			if numeric, ok := x.(Numeric); ok {
				return &Return{&Float{math.Log2(numeric.Float64())}}
			} else {
				ctx.interp.throw("std/math:log2(x): expected number")
				return &Return{}
			}
		}},
		"F/log10": &Func{Executor: func(ctx *FuncContext) *Return {
			if len(ctx.Args) < 1 {
				ctx.interp.throw("std/math:log10(x): expected 1 argument")
				return &Return{}
			}
			x := ctx.Args[0]
			if numeric, ok := x.(Numeric); ok {
				return &Return{&Float{math.Log10(numeric.Float64())}}
			} else {
				ctx.interp.throw("std/math:log10(x): expected number")
				return &Return{}
			}
		}},
		"F/exp": &Func{Executor: func(ctx *FuncContext) *Return {
			if len(ctx.Args) < 1 {
				ctx.interp.throw("std/math:exp(x): expected 1 argument")
				return &Return{}
			}
			x := ctx.Args[0]
			if numeric, ok := x.(Numeric); ok {
				return &Return{&Float{math.Exp(numeric.Float64())}}
			} else {
				ctx.interp.throw("std/math:exp(x): expected number")
				return &Return{}
			}
		}},
		"F/pow": &Func{Executor: func(ctx *FuncContext) *Return {
			if len(ctx.Args) < 2 {
				ctx.interp.throw("std/math:pow(x, y): expected 2 arguments")
				return &Return{}
			}
			x := ctx.Args[0]
			y := ctx.Args[1]
			if numericX, ok := x.(Numeric); ok {
				if numericY, ok := y.(Numeric); ok {
					return &Return{&Float{math.Pow(numericX.Float64(), numericY.Float64())}}
				} else {
					ctx.interp.throw("std/math:pow(x, y): expected number")
					return &Return{}
				}
			} else {
				ctx.interp.throw("std/math:pow(x, y): expected number")
				return &Return{}
			}
		}},
		"F/abs": &Func{Executor: func(ctx *FuncContext) *Return {
			if len(ctx.Args) < 1 {
				ctx.interp.throw("std/math:abs(x): expected 1 argument")
				return &Return{}
			}
			x := ctx.Args[0]
			if numeric, ok := x.(Numeric); ok {
				return &Return{&Float{math.Abs(numeric.Float64())}}
			} else {
				ctx.interp.throw("std/math:abs(x): expected number")
				return &Return{}
			}
		}},
		"F/sign": &Func{Executor: func(ctx *FuncContext) *Return {
			if len(ctx.Args) < 1 {
				ctx.interp.throw("std/math:sign(x): expected 1 argument")
				return &Return{}
			}
			x := ctx.Args[0]
			if numeric, ok := x.(Numeric); ok {
				if numeric.Float64() < 0 {
					return &Return{&Float{-1}}
				} else if numeric.Float64() > 0 {
					return &Return{&Float{1}}
				} else {
					return &Return{&Float{0}}
				}
			} else {
				ctx.interp.throw("std/math:sign(x): expected number")
				return &Return{}
			}
		}},
		"F/min": &Func{Executor: func(ctx *FuncContext) *Return {
			if len(ctx.Args) < 1 {
				ctx.interp.throw("std/math:min(nums): expected at least 1 argument")
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
					ctx.interp.throw("std/math:min(nums): expected number")
					return &Return{}
				}
			}

			return &Return{wrap(min)}
		}},
		"F/max": &Func{Executor: func(ctx *FuncContext) *Return {
			if len(ctx.Args) < 1 {
				ctx.interp.throw("std/math:max(nums): expected at least 1 argument")
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
					ctx.interp.throw("std/math:max(nums): expected number")
					return &Return{}
				}
			}

			return &Return{wrap(max)}
		}},
		"F/clamp": &Func{Executor: func(ctx *FuncContext) *Return {
			if len(ctx.Args) < 3 {
				ctx.interp.throw("std/math:clamp(x, min, max): expected 3 arguments")
				return &Return{}
			}
			x := ctx.Args[0]
			min := ctx.Args[1]
			max := ctx.Args[2]
			if numericX, ok := x.(Numeric); ok {
				if numericMin, ok := min.(Numeric); ok {
					if numericMax, ok := max.(Numeric); ok {
						return &Return{&Float{math.Min(math.Max(numericX.Float64(), numericMin.Float64()), numericMax.Float64())}}
					} else {
						ctx.interp.throw("std/math:clamp(x, min, max): expected number")
						return &Return{}
					}
				} else {
					ctx.interp.throw("std/math:clamp(x, min, max): expected number")
					return &Return{}
				}
			} else {
				ctx.interp.throw("std/math:clamp(x, min, max): expected number")
				return &Return{}
			}
		}},
		"F/parseInt": &Func{Executor: func(ctx *FuncContext) *Return {
			if len(ctx.Args) < 1 {
				ctx.interp.throw("std/math:parseInt(s): expected 1 argument")
				return &Return{}
			}
			s := ctx.Args[0]
			base := 10
			if len(ctx.Args) > 1 {
				if integer, ok := ctx.Args[1].(*Integer); ok {
					base = int(integer.Value)
				} else {
					ctx.interp.throw("std/math:parseInt(s, base): expected integer")
				}
			}
			if str, ok := s.(*String); ok {
				i, err := strconv.ParseInt(string(str.Value), base, 64)
				if err != nil {
					ctx.interp.throw("std/math:parseInt(s): " + err.Error())
					return &Return{}
				}
				return &Return{wrap(i)}
			} else {
				ctx.interp.throw("std/math:parseInt(s): expected string")
				return &Return{}
			}
		}},

		"C/pi":     &Float{math.Pi},
		"C/e":      &Float{math.E},
		"C/phi":    &Float{math.Phi},
		"C/sqrt2":  &Float{math.Sqrt2},
		"C/ln2":    &Float{math.Ln2},
		"C/log2e":  &Float{math.Log2E},
		"C/ln10":   &Float{math.Ln10},
		"C/log10e": &Float{math.Log10E},
	},
}
