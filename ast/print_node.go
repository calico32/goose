package ast

import (
	"fmt"
	"strings"
)

type NodePrinter struct {
	indent       int
	output       strings.Builder
	shouldIndent bool
}

func (p *NodePrinter) writeIndent() {
	for i := 0; i < p.indent; i++ {
		p.output.WriteString("    ")
	}
}

func (p *NodePrinter) writeBlock(block []Stmt) {
	p.write("{\n")
	p.indent++
	for _, stmt := range block {
		p.Print(stmt)
	}
	p.indent--
	p.write("}")
}

func (p *NodePrinter) write(s string) {
	if p.shouldIndent {
		p.shouldIndent = false
		p.writeIndent()
	}
	for i, c := range s {
		if c == '\n' {
			p.output.WriteRune(c)
			if i+1 < len(s) {
				p.writeIndent()
			} else {
				p.shouldIndent = true
			}
		} else {
			p.output.WriteRune(c)
		}
	}
}

func (p *NodePrinter) String() string {
	return p.output.String()
}

func (p *NodePrinter) Print(node Node) {
	switch n := node.(type) {
	case *BinaryExpr:
		p.write("(")
		p.Print(n.X)
		p.write(" ")
		p.write(n.Op.String())
		p.write(" ")
		p.Print(n.Y)
		p.write(")")
		// p.write("(\n")
		// p.indent++
		// p.Print(n.X)
		// p.write("\n")
		// p.write(n.Op.String())
		// p.write("\n")
		// p.Print(n.Y)
		// p.indent--
		// p.write("\n)")
	case *ParenExpr:
		p.Print(n.X)
	case *UnaryExpr:
		p.write("(")
		p.write(n.Op.String())
		p.write(" ")
		p.Print(n.X)
		p.write(")")
	case *Ident:
		p.write(n.Name)
	case *FuncExpr:
		if n.Async.IsValid() {
			p.write("async ")
		}
		if n.Memo.IsValid() {
			p.write("memo ")
		}
		p.write("fn")

		if n.Receiver != nil {
			p.write(" ")
			p.Print(n.Receiver)
			p.write(".")
			if n.Name != nil {
				p.Print(n.Name)
			}
		} else if n.Name != nil {
			p.write(" ")
			p.Print(n.Name)
		}

		p.write("(")
		for i, arg := range n.Params.List {
			if i > 0 {
				p.write(", ")
			}
			p.Print(arg)
		}
		p.write(")")

		if n.Arrow.IsValid() {
			p.write(" -> ")
			p.Print(n.ArrowExpr)
		} else {
			p.write(" ")
			p.writeBlock(n.Body)
		}
	case *FuncParam:
		if n.Ellipsis.IsValid() {
			p.write("...")
		}
		p.Print(n.Ident)
		if n.Value != nil {
			p.write(" = ")
			p.Print(n.Value)
		}

	case *CallExpr:
		p.Print(n.Func)

		p.write("(")
		for i, arg := range n.Args {
			if i > 0 {
				p.write(", ")
			}
			p.Print(arg)
		}
		p.write(")")

	case *SelectorExpr:
		p.Print(n.X)
		p.write(".")
		p.Print(n.Sel)
	case *ExprStmt:
		p.Print(n.X)
	case *StringLiteral:
		p.write("\"")
		p.write(n.StringStart.Content)
		for _, part := range n.Parts {
			switch part := part.(type) {
			case *StringLiteralMiddle:
				p.write(part.Content)
			case *StringLiteralInterpExpr:
				p.write("${")
				p.Print(part.Expr)
				p.write("}")
			case *StringLiteralInterpIdent:
				p.write("$")
				p.write(part.Name)
			}
		}
		p.write(n.StringEnd.Content)
		p.write("\"")
	case *LetStmt:
		p.write("let ")
		p.Print(n.Ident)
		if n.Value != nil {
			p.write(" = ")
			p.Print(n.Value)
		}
	case *ConstStmt:
		p.write("const ")
		p.Print(n.Ident)
		p.write(" = ")
		p.Print(n.Value)

	case *IfStmt:
		p.write("if ")
		p.Print(n.Cond)
		p.write(" ")
		p.writeBlock(n.Body)
		if len(n.Else) != 0 {
			p.write(" else ")
			p.writeBlock(n.Else)
		}
	case *ForStmt:
		p.write("for ")
		if n.Await.IsValid() {
			p.write("await ")
		}
		p.Print(n.Var)
		p.write(" in ")
		p.Print(n.Iterable)
		p.write(" ")
		p.writeBlock(n.Body)
	case *ReturnStmt:
		p.write("return")
		if n.Result != nil {
			p.write(" ")
			p.Print(n.Result)
		}
	case *BranchStmt:
		p.write(n.Tok.String())
		if n.Label != nil {
			p.write(" ")
			p.Print(n.Label)
		}
	case *LabeledStmt:
		if p.indent != 0 {
			p.indent--
			p.write("\n")
			p.Print(n.Label)
			p.write(":\n")
			p.indent++
		} else {
			p.Print(n.Label)
			p.write(":\n")
		}
		p.Print(n.Stmt)
	case *SymbolStmt:
		p.write("symbol ")
		p.Print(n.Ident)
	case *IfExpr:
		p.write("if ")
		p.Print(n.Cond)
		p.write(" ")
		p.Print(n.Then)
		p.write(" else ")
		p.Print(n.Else)
	case *Literal:
		p.write(n.Value)
	case *ArrayLiteral:
		p.write("[")
		for i, elem := range n.List {
			if i > 0 {
				p.write(", ")
			}
			p.Print(elem)
		}
		p.write("]")
	case *ArrayInitializer:
		p.write("[")
		p.Print(n.Count)
		p.write("; ")
		p.Print(n.Value)
		p.write("]")
	case *AssignStmt:
		p.Print(n.Lhs)
		p.write(" ")
		p.write(n.Tok.String())
		p.write(" ")
		p.Print(n.Rhs)
	case *IncDecStmt:
		p.Print(n.X)
		p.write(n.Tok.String())
	case *CompositeLiteral:
		p.write("{\n")
		p.indent++
		for i, elem := range n.Fields {
			if i > 0 {
				p.write(",\n")
			}
			p.Print(elem.Key)
			p.write(": ")
			p.Print(elem.Value)
		}
		p.indent--
		p.write("\n}")
	case *BindExpr:
		p.Print(n.X)
		p.write("::")
		p.Print(n.Sel)
	case *BracketPropertyExpr:
		p.write("#[")
		p.Print(n.X)
		p.write("]")
	case *BracketSelectorExpr:
		p.Print(n.X)
		p.write("[")
		p.Print(n.Sel)
		p.write("]")
	case *TryStmt:
		p.write("try ")
		p.writeBlock(n.Body)
		if n.Catch != nil {
			p.write(" catch ")
			if n.Catch.Ident != nil {
				p.write("as ")
				p.Print(n.Catch.Ident)
			}
			p.write(" ")
			p.writeBlock(n.Catch.Body)
		}
		if n.Finally != nil {
			p.write(" finally ")
			p.writeBlock(n.Finally.Body)
		}
	case *ThrowExpr:
		p.write("throw ")
		p.Print(n.X)
	case *RangeExpr:
		p.write("(")
		p.Print(n.Low)
		p.write(" to ")
		p.Print(n.High)
		if n.Step != nil {
			p.write(" step ")
			p.Print(n.Step)
		}
		p.write(")")
	case *RepeatCountStmt:
		p.write("repeat ")
		p.Print(n.Count)
		p.write(" times ")
		p.writeBlock(n.Body)
	case *RepeatWhileStmt:
		p.write("repeat while ")
		p.Print(n.Cond)
		p.write(" ")
		p.writeBlock(n.Body)
	case *RepeatForeverStmt:
		p.write("repeat forever ")
		p.writeBlock(n.Body)
	case *ImportStmt:
		p.write("import ")
		p.Print(n.Spec)
	case *ModuleSpecPlain:
		p.write("\"")
		p.write(n.Specifier)
		p.write("\"")
	case *ModuleSpecAs:
		p.write("\"")
		p.write(n.Specifier)
		p.write("\" as ")
		p.Print(n.Alias)
	case *ModuleSpecShow:
		p.write("\"")
		p.write(n.Specifier)
		p.write("\" show ")
		p.Print(n.Show)
	case *Show:
		if n.Ellipsis.IsValid() {
			p.write("...")
			break
		}
		p.write("{\n")
		p.indent++
		for i, elem := range n.Fields {
			if i > 0 {
				p.write(",\n")
			}
			p.Print(elem)
		}
		p.indent--
		p.write("\n}")
	case *ShowFieldAs:
		p.Print(n.Ident)
		p.write(" as ")
		p.Print(n.Alias)
	case *ShowFieldEllipsis:
		p.write("...")
		p.Print(n.Ident)
	case *ShowFieldIdent:
		p.Print(n.Ident)
	case *ShowFieldSpec:
		p.Print(n.Spec)
	case *ExportDeclStmt:
		p.write("export ")
		p.Print(n.Stmt)
	case *ExportListStmt:
		p.write("export ")
		p.write("{\n")
		p.indent++
		for i, elem := range n.List.Fields {
			if i > 0 {
				p.write(",\n")
			}
			p.Print(elem)
		}
		p.indent--
		p.write("\n}")
	case *ExportSpecStmt:
		p.write("export ")
		p.Print(n.Spec)
	case *ExportFieldAs:
		p.Print(n.Ident)
		p.write(" as ")
		p.Print(n.Alias)
	case *ExportFieldIdent:
		p.Print(n.Ident)
	case *GeneratorExpr:
		if n.Async.IsValid() {
			p.write("async ")
		}
		p.write("generator ")
		if n.Receiver != nil {
			p.Print(n.Receiver)
			p.write(".")
		}
		p.Print(n.Name)

		p.write("(")
		for i, param := range n.Params.List {
			if i > 0 {
				p.write(", ")
			}
			p.Print(param)
		}
		p.write(")")

		p.write(" ")
		p.writeBlock(n.Body)
	case *StructStmt:
		p.write("struct ")
		p.Print(n.Name)
		p.write("(")
		for i, field := range n.Fields.List {
			if i > 0 {
				p.write(", ")
			}
			p.Print(field.Ident)
			if field.Value != nil {
				p.write(" = ")
				p.Print(field.Value)
			}
		}
		p.write(")")
		if n.Init != nil {
			p.write(" init ")
			p.writeBlock(n.Init.Body)
		}
	case *OperatorStmt:
		if n.Async.IsValid() {
			p.write("async ")
		}
		p.write("operator ")
		p.Print(n.Receiver)
		p.write(" ")
		p.write(n.Tok.String())
		p.write("(")
		for i, param := range n.Params.List {
			if i > 0 {
				p.write(", ")
			}
			if param.Ellipsis.IsValid() {
				p.write("...")
			}
			p.Print(param.Ident)
			if param.Value != nil {
				p.write(" = ")
				p.Print(param.Value)
			}
		}
		p.write(") ")
		p.writeBlock(n.Body)

	default:
		p.write(fmt.Sprintf("<unhandled %T>", node))
	}

	if _, ok := node.(Stmt); ok {
		p.write(";\n")
	}
}
