package interpreter

import "github.com/calico32/goose/ast"

func (i *interp) evalCompositeLiteral(scope *Scope, expr *ast.CompositeLiteral) Value {
	defer un(trace(i, "composite literal"))

	composite := NewComposite()

	for _, field := range expr.Fields {
		var keyValue PropertyKey
		switch key := field.Key.(type) {
		case *ast.Ident:
			keyValue = wrap(key.Name).(PropertyKey)
		case *ast.StringLiteral:
			k := i.evalString(scope, key)

			keyValue = k.(PropertyKey)
		default:
			lit := i.evalExpr(scope, key)

			if int, ok := lit.(*Integer); ok {
				keyValue = wrap(int.Value).(PropertyKey)
				break
			}

			i.throw("invalid composite literal key type %s", lit.Type())
		}

		val := i.evalExpr(scope, field.Value)

		SetProperty(composite, keyValue, val)
	}

	return composite
}
