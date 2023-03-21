package parser

import (
	"github.com/calico32/goose/ast"
	"github.com/calico32/goose/token"
)

func (p *Parser) parseImportStmt() *ast.ImportStmt {
	if p.trace {
		defer un(trace(p, "ImportStmt"))
	}

	pos := p.expect(token.Import)
	spec := p.parseImportExportSpec()

	return &ast.ImportStmt{
		Import: pos,
		Spec:   spec,
	}
}

func (p *Parser) parseExportStmt() ast.Stmt {
	if p.trace {
		defer un(trace(p, "ExportStmt"))
	}

	exportPos := p.expect(token.Export)

	switch p.tok {
	case token.StringStart:
		return &ast.ExportSpecStmt{
			Export: exportPos,
			Spec:   p.parseImportExportSpec(),
		}
	case token.LBrace:
		return &ast.ExportListStmt{
			Export: exportPos,
			List:   p.parseExportList(),
		}
	case token.Symbol,
		token.Func, token.Async, token.Memo,
		token.Let, token.Const,
		token.Struct, token.Generator,
		token.Native:
		return &ast.ExportDeclStmt{
			Export: exportPos,
			Stmt:   p.parseStmt(),
		}
	}

	p.errorExpected(p.pos, "module specifier, export list, or declaration")
	return &ast.BadStmt{From: exportPos, To: p.pos}
}

func (p *Parser) parseImportExportSpec() ast.ModuleSpec {
	specifierString := p.parseString()
	if len(specifierString.Parts) != 0 {
		// string has interpolation, so we can't parse it
		p.errorExpected(specifierString.Pos(), "import/export specifier")
	}

	specifier := specifierString.StringStart.Content

	switch p.tok {
	case token.As:
		spec := &ast.ModuleSpecAs{
			SpecifierPos: specifierString.Pos(),
			Specifier:    specifier,
			As:           p.pos,
		}
		p.next()
		alias := p.parseIdent()
		spec.Alias = alias
		return spec
	case token.Show:
		show := &ast.Show{Show: p.pos}
		p.next()
		switch p.tok {
		case token.Ellipsis:
			show.Ellipsis = p.pos
			p.next()
		case token.LBrace:
			show.LBrace = p.pos
			p.next()
			show.Fields = p.parseImportList()
			show.RBrace = p.expect(token.RBrace)
		default:
			p.errorExpected(p.pos, "ellipsis or import list")
		}

		return &ast.ModuleSpecShow{
			SpecifierPos: specifierString.Pos(),
			Specifier:    specifier,
			Show:         show,
		}
	}

	return &ast.ModuleSpecPlain{
		SpecifierPos: specifierString.Pos(),
		Specifier:    specifier,
	}
}

func (p *Parser) parseImportList() []ast.ShowField {
	if p.trace {
		defer un(trace(p, "ImportList"))
	}

	list := []ast.ShowField{}
	noMoreIdents := false

	for p.tok != token.RBrace && p.tok != token.EOF {
		switch p.tok {
		case token.Ellipsis:
			p.next()
			ident := p.parseIdent()
			list = append(list, &ast.ShowFieldEllipsis{Ellipsis: ident.Pos(), Ident: ident})
			noMoreIdents = true
		case token.Ident:
			if noMoreIdents {
				p.errorExpected(p.pos, "submodule import or end of import list")
			}
			ident := p.parseIdent()
			if p.tok == token.As {
				field := &ast.ShowFieldAs{Ident: ident}
				field.As = p.pos
				p.next()
				field.Alias = p.parseIdent()
				list = append(list, field)
			} else {
				list = append(list, &ast.ShowFieldIdent{Ident: ident})
			}
		case token.StringStart:
			spec := p.parseImportExportSpec()
			field := &ast.ShowFieldSpec{Spec: spec}
			list = append(list, field)
		case token.Comma:
			p.next()
			// two commas in a row, error
			fallthrough
		default:
			p.errorExpected(p.pos, "identifier or submodule import")
		}

		if p.tok != token.RBrace {
			p.expect(token.Comma)
		}
	}

	return list
}

func (p *Parser) parseExportList() *ast.ExportList {
	if p.trace {
		defer un(trace(p, "ExportList"))
	}

	list := &ast.ExportList{
		LBrace: p.expect(token.LBrace),
		Fields: []ast.ExportField{},
	}

	for p.tok != token.RBrace && p.tok != token.EOF {
		switch p.tok {
		case token.Ident:
			ident := p.parseIdent()
			if p.tok == token.As {
				field := &ast.ExportFieldAs{Ident: ident}
				field.As = p.pos
				p.next()
				field.Alias = p.parseIdent()
				list.Fields = append(list.Fields, field)
			} else {
				list.Fields = append(list.Fields, &ast.ExportFieldIdent{Ident: ident})
			}
		case token.Comma:
			p.next()
			// two commas in a row, error
			fallthrough
		default:
			p.errorExpected(p.pos, "identifier")
		}

		if p.tok != token.RBrace {
			p.expect(token.Comma)
		}
	}

	list.RBrace = p.expect(token.RBrace)

	return list
}
