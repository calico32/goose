package ast

import (
	"errors"
	"strings"

	"github.com/calico32/goose/token"
)

type (
	ImportStmt struct {
		Import token.Pos
		Spec   ModuleSpec
	}

	ExportDeclStmt struct {
		Export token.Pos
		Stmt   Stmt
	}

	ExportListStmt struct {
		Export token.Pos
		List   *ExportList
	}

	ExportSpecStmt struct {
		Export token.Pos
		Spec   ModuleSpec
	}
)

func (s *ImportStmt) Pos() token.Pos     { return s.Import }
func (s *ExportDeclStmt) Pos() token.Pos { return s.Export }
func (s *ExportListStmt) Pos() token.Pos { return s.Export }
func (s *ExportSpecStmt) Pos() token.Pos { return s.Export }

func (s *ImportStmt) End() token.Pos     { return s.Spec.End() }
func (s *ExportDeclStmt) End() token.Pos { return s.Stmt.End() }
func (s *ExportListStmt) End() token.Pos { return s.List.End() }
func (s *ExportSpecStmt) End() token.Pos { return s.Spec.End() }

func (*ImportStmt) stmtNode()     {}
func (*ExportDeclStmt) stmtNode() {}
func (*ExportListStmt) stmtNode() {}
func (*ExportSpecStmt) stmtNode() {}

type (
	ModuleSpec interface {
		Node
		importExportSpec()
		ModuleSpecifier() string
	}

	ModuleSpecPlain struct {
		SpecifierPos token.Pos
		Specifier    string
	}

	ModuleSpecAs struct {
		SpecifierPos token.Pos
		Specifier    string
		As           token.Pos
		Alias        *Ident
	}

	ModuleSpecShow struct {
		SpecifierPos token.Pos
		Specifier    string
		Show         *Show
	}

	Show struct {
		Show token.Pos

		// import "foo" show ...
		Ellipsis token.Pos

		LBrace token.Pos
		Fields []ShowField
		RBrace token.Pos
	}

	ShowField interface {
		Node
		showField()
	}

	// import "foo" show { foo }
	ShowFieldIdent struct {
		Ident *Ident
	}

	// import "foo" show { ...foo }
	ShowFieldEllipsis struct {
		Ellipsis token.Pos
		Ident    *Ident
	}

	// import "foo" show { foo as foo2 }
	ShowFieldAs struct {
		Ident *Ident
		As    token.Pos
		Alias *Ident
	}

	// import "foo" show { "utils", "utils" as fooUtils, "utils" show { foo } }
	ShowFieldSpec struct {
		Spec ModuleSpec
	}

	ExportList struct {
		LBrace token.Pos
		Fields []ExportField
		RBrace token.Pos
	}

	ExportField interface {
		Node
		exportField()
	}

	// import "foo" show { foo }
	ExportFieldIdent struct {
		Ident *Ident
	}

	// import "foo" show { foo as foo2 }
	ExportFieldAs struct {
		Ident *Ident
		As    token.Pos
		Alias *Ident
	}
)

func (s *ModuleSpecPlain) Pos() token.Pos { return s.SpecifierPos }
func (s *ModuleSpecAs) Pos() token.Pos    { return s.SpecifierPos }
func (s *ModuleSpecShow) Pos() token.Pos  { return s.SpecifierPos }

func (s *Show) Pos() token.Pos              { return s.Show }
func (s *ShowFieldIdent) Pos() token.Pos    { return s.Ident.Pos() }
func (s *ShowFieldEllipsis) Pos() token.Pos { return s.Ellipsis }
func (s *ShowFieldAs) Pos() token.Pos       { return s.Ident.Pos() }
func (s *ShowFieldSpec) Pos() token.Pos     { return s.Spec.Pos() }

func (s *ExportList) Pos() token.Pos       { return s.LBrace }
func (s *ExportFieldIdent) Pos() token.Pos { return s.Ident.Pos() }
func (s *ExportFieldAs) Pos() token.Pos    { return s.Ident.Pos() }

func (s *ModuleSpecPlain) End() token.Pos { return s.SpecifierPos + token.Pos(len(s.Specifier)) + 2 }
func (s *ModuleSpecAs) End() token.Pos    { return s.Alias.End() }
func (s *ModuleSpecShow) End() token.Pos  { return s.Show.RBrace + 1 }

func (s *Show) End() token.Pos {
	if s.Ellipsis.IsValid() {
		return s.Ellipsis + 3
	} else {
		return s.RBrace + 1
	}
}
func (s *ShowFieldIdent) End() token.Pos    { return s.Ident.End() }
func (s *ShowFieldEllipsis) End() token.Pos { return s.Ident.End() }
func (s *ShowFieldAs) End() token.Pos       { return s.Alias.End() }
func (s *ShowFieldSpec) End() token.Pos     { return s.Spec.End() }

func (s *ExportList) End() token.Pos       { return s.RBrace + 1 }
func (s *ExportFieldIdent) End() token.Pos { return s.Ident.End() }
func (s *ExportFieldAs) End() token.Pos    { return s.Alias.End() }

func (*ModuleSpecPlain) importExportSpec() {}
func (*ModuleSpecAs) importExportSpec()    {}
func (*ModuleSpecShow) importExportSpec()  {}

func (s *ModuleSpecPlain) ModuleSpecifier() string { return s.Specifier }
func (s *ModuleSpecAs) ModuleSpecifier() string    { return s.Specifier }
func (s *ModuleSpecShow) ModuleSpecifier() string  { return s.Specifier }

func (*ShowFieldIdent) showField()    {}
func (*ShowFieldEllipsis) showField() {}
func (*ShowFieldAs) showField()       {}
func (*ShowFieldSpec) showField()     {}

func (*ExportFieldIdent) exportField() {}
func (*ExportFieldAs) exportField()    {}

func ModuleName(specifier string) (name string, err error) {
	if specifier == "" {
		return "", errors.New("empty module name")
	}

	if strings.HasSuffix(specifier, "/") {
		return "", errors.New("trailing slash in module name")
	}

	isRelative := strings.HasPrefix(specifier, "./") || strings.HasPrefix(specifier, "../") || strings.HasPrefix(specifier, "/")

	parts := strings.Split(specifier, "/")
	if !isRelative {
		name = parts[len(parts)-1]
		if name == "_module.goose" {
			if len(parts) == 1 {
				return "", errors.New("invalid module name")
			}
			name = parts[len(parts)-2]
		}
	} else {
		if len(parts) == 1 { // ., .., /
			return "", errors.New("invalid module name")
		}

		name = parts[len(parts)-1]

		if name == "_module.goose" {
			if len(parts) == 2 {
				return "", errors.New("invalid module name")
			}
			name = parts[len(parts)-2]
		}
	}

	if token.Lookup(name).IsKeyword() {
		return "", errors.New("invalid module name")
	}

	name = strings.TrimSuffix(name, ".goose")
	name = strings.ReplaceAll(name, ".", "")

	cleanName := ""

	validChars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	for _, c := range name {
		if strings.ContainsRune(validChars, c) {
			cleanName += string(c)
		}
	}

	if cleanName == "" {
		return "", errors.New("invalid module name")
	}

	invalidStartChars := "0123456789"
	if strings.ContainsRune(invalidStartChars, rune(cleanName[0])) {
		return "", errors.New("invalid module name")
	}

	if cleanName == "" || cleanName == "_" || cleanName == "_module" {
		return "", errors.New("invalid module name")
	}

	name = cleanName

	return
}
