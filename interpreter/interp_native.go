package interpreter

import (
	"fmt"

	"github.com/calico32/goose/ast"
	. "github.com/calico32/goose/interpreter/lib"

	std_fs "github.com/calico32/goose/lib/std/fs"
	std_json "github.com/calico32/goose/lib/std/json"
	std_language "github.com/calico32/goose/lib/std/language"
	std_math "github.com/calico32/goose/lib/std/math"
	std_platform "github.com/calico32/goose/lib/std/platform"
	std_random "github.com/calico32/goose/lib/std/random"
	std_readline "github.com/calico32/goose/lib/std/readline"
)

func (i *interp) evalNativeExpr(scope *Scope, expr *ast.NativeExpr) Value {
	module := scope.Module()
	if moduleNatives, ok := Natives[module.Specifier]; ok {
		if value, ok := moduleNatives[expr.Id]; ok {
			return value
		} else {
			i.Throw("native symbol %s not found in module %s", expr.Id, module.Specifier)
			panic("unreachable")
		}
	} else {
		i.Throw("native module %s not found", module.Specifier)
		panic("unreachable")
	}
}

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
		i.Throw(fmt.Sprintf("invalid native stmt type %T", stmt))
	}

	specifier := scope.Module().Specifier

	if moduleNatives, ok := Natives[specifier]; ok {
		if value, ok := moduleNatives[name]; ok {
			if fn, ok := stmt.(*ast.NativeFunc); ok && fn.Receiver != nil {
				// find proto
				// TODO: limit to current module
				constructor := scope.Get(fn.Receiver.Name)
				if constructor == nil {
					i.Throw("unknown type %s", fn.Receiver.Name)
				}

				if val, ok := constructor.Value.(*Func); !ok || val.NewableProto == nil {
					i.Throw("%s is not a type", fn.Receiver.Name)
				}

				proto := constructor.Value.(*Func).NewableProto

				if proto.Properties[PKString] == nil {
					proto.Properties[PKString] = make(map[string]Value)
				}

				if _, ok := proto.Properties[PKString][fn.Name.Name]; ok {
					i.Throw("duplicate receiver function %s", fn.Name.Name)
				}

				proto.Properties[PKString][fn.Name.Name] = value
			}

			scope.Set(name[2:], &Variable{
				Value:    value,
				Constant: true,
			})
			return &Decl{
				Name:  name[2:],
				Value: value,
			}
		}
	}

	i.Throw("native symbol %s not found in module %s", name, specifier)
	return nil
}

var Natives = map[string]map[string]Value{
	"std:language/builtin.goose": std_language.Builtin,

	"std:fs/index.goose":       std_fs.Index,
	"std:json/index.goose":     std_json.Index,
	"std:math/index.goose":     std_math.Index,
	"std:platform/index.goose": std_platform.Index,
	"std:random/index.goose":   std_random.Index,
	"std:readline/index.goose": std_readline.Index,
}
