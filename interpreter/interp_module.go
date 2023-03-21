package interpreter

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/calico32/goose/ast"
	"github.com/calico32/goose/parser"
)

func (i *interp) runExportDeclStmt(scope *Scope, stmt *ast.ExportDeclStmt) StmtResult {
	result := i.runStmt(scope, stmt.Stmt)

	if scope != scope.ModuleScope() {
		i.throw("export declarations must be at the top level")
	}

	if decl, ok := result.(*Decl); !ok {
		i.throw("declaration expected")
	} else if _, ok := scope.Module().Exports[decl.Name]; ok {
		i.throw("duplicate export %s", decl.Name)
	} else {
		scope.Module().Exports[decl.Name] = scope.Get(decl.Name)
	}

	return &Void{}
}

func (i *interp) runExportListStmt(scope *Scope, stmt *ast.ExportListStmt) StmtResult {
	if scope != scope.ModuleScope() {
		i.throw("export declarations must be at the top level")
	}

	for _, field := range stmt.List.Fields {
		var exportedName string
		var localName string
		switch field := field.(type) {
		case *ast.ExportFieldIdent:
			localName = field.Ident.Name
			exportedName = field.Ident.Name
		case *ast.ExportFieldAs:
			localName = field.Ident.Name
			exportedName = field.Alias.Name
		default:
			i.throw("unknown export field type %T", field)
		}

		if _, ok := scope.Module().Exports[exportedName]; ok {
			i.throw("duplicate export %s", exportedName)
		}

		if !scope.IsDefinedInCurrentScope(localName) {
			i.throw("undefined name %s", localName)
		}

		scope.Module().Exports[exportedName] = scope.Get(localName)
	}

	return &Void{}
}

func (i *interp) runExportSpecStmt(scope *Scope, stmt *ast.ExportSpecStmt) StmtResult {
	if scope != scope.ModuleScope() {
		i.throw("export declarations must be at the top level")
	}

	importScope := &Scope{
		owner:  ScopeOwnerImport,
		idents: make(map[string]*Variable),
	}

	i.runImportSpec(importScope, stmt.Spec)

	for name := range importScope.idents {
		if _, ok := scope.Module().Exports[name]; ok {
			i.throw("duplicate export %s", name)
		}

		scope.Module().Exports[name] = importScope.Get(name)
	}

	return &Void{}
}

func (i *interp) runImportStmt(scope *Scope, stmt *ast.ImportStmt) StmtResult {
	if scope != scope.ModuleScope() {
		i.throw("import declarations must be at the top level")
	}

	return i.runImportSpec(scope, stmt.Spec)
}

func (i *interp) runImportSpec(scope *Scope, spec ast.ModuleSpec) StmtResult {
	// TODO: custom import schemes
	var scheme string
	var name string
	specifier := spec.ModuleSpecifier()
	if strings.Contains(specifier, ":") {
		colon := strings.Index(specifier, ":")
		scheme = specifier[:colon]
		name = specifier[colon+1:]
	} else if isFilePath(specifier) {
		scheme = "file"
		name = specifier
		if !isFilePath(specifier) {
			i.throw("invalid file import path %s", specifier)
		}
	} else {
		scheme = "pkg"
		name = specifier
		if isFilePath(specifier) {
			i.throw("invalid package import path %s", specifier)
		}
	}

	var module *Module
	switch scheme {
	case "file":
		dir := filepath.Dir(scope.Module().Name)
		module = i.loadFileModule(name, dir)
	case "pkg":
		module = i.loadPackageModule(name)
	}

	switch spec := spec.(type) {
	case *ast.ModuleSpecShow:
		if spec.Show.Ellipsis.IsValid() {
			for name, value := range module.Exports {
				scope.Set(name, value)
			}
		} else {
			imported := make(map[string]bool)
			for _, field := range spec.Show.Fields {

				switch field := field.(type) {
				case *ast.ShowFieldSpec:
					var spec ast.ModuleSpec
					switch original := field.Spec.(type) {
					case *ast.ModuleSpecAs:
						spec = &ast.ModuleSpecAs{
							SpecifierPos: original.SpecifierPos,
							Specifier:    name + "/" + original.Specifier,
							As:           original.As,
							Alias:        original.Alias,
						}
					case *ast.ModuleSpecShow:
						spec = &ast.ModuleSpecShow{
							SpecifierPos: original.SpecifierPos,
							Specifier:    name + "/" + original.Specifier,
							Show:         original.Show,
						}
					case *ast.ModuleSpecPlain:
						spec = &ast.ModuleSpecPlain{
							SpecifierPos: original.SpecifierPos,
							Specifier:    name + "/" + original.Specifier,
						}
					}
					i.runImportSpec(scope, spec)
				case *ast.ShowFieldEllipsis:
					if scope.IsDefinedInCurrentScope(field.Ident.Name) {
						i.throw("name %s is already defined", field.Ident.Name)
					}
					// put all remaining exports into an object
					object := NewComposite()
					for name, value := range module.Exports {
						if imported[name] {
							continue
						}
						SetProperty(object, &String{name}, value.Value)
					}
					object.Frozen = true

					scope.Set(field.Ident.Name, &Variable{
						Value:    object,
						Constant: true,
					})
				default:
					var localName string
					var exportedName string

					if identField, ok := field.(*ast.ShowFieldIdent); ok {
						localName = identField.Ident.Name
						exportedName = identField.Ident.Name
					} else if asField, ok := field.(*ast.ShowFieldAs); ok {
						localName = asField.Alias.Name
						exportedName = asField.Ident.Name
					} else {
						i.throw("unknown import field type %T", field)
					}

					if _, ok := module.Exports[exportedName]; !ok {
						if module.Scope.IsDefinedInCurrentScope(exportedName) {
							i.throw("value %s is defined locally in module %s but is not exported", exportedName, name)
						}
						i.throw("undefined export %s", exportedName)
					}

					imported[exportedName] = true

					value := module.Scope.Get(exportedName)
					scope.Set(localName, value)
				}
			}
		}
	default:
		object := NewComposite()
		for name, value := range module.Exports {
			SetProperty(object, &String{name}, value.Value) // TODO: reassignment can change the value
		}
		object.Frozen = true

		var moduleName string

		if aliased, ok := spec.(*ast.ModuleSpecAs); ok {
			moduleName = aliased.Alias.Name
		} else {
			var err error
			moduleName, err = ast.ModuleName(name)
			if err != nil {
				i.throw(err.Error())
			}
		}

		if scope.IsDefinedInCurrentScope(moduleName) {
			i.throw("name %s is already defined", moduleName)
		}

		scope.Set(moduleName, &Variable{
			Constant: true,
			Value:    object,
		})
	}

	return &Void{}
}

func (i *interp) loadFileModule(name string, dir string) *Module {
	if module, ok := i.modules[name]; ok {
		return module
	}

	if !filepath.IsAbs(name) {
		name = filepath.Join(dir, name)
	}

	info, err := os.Stat(name)
	if err != nil {
		i.throw(err.Error())
	}

	if info.IsDir() {
		newName := filepath.Join(name, "_module.goose")
		_, err = os.Stat(newName)
		if err != nil {
			if os.IsNotExist(err) {
				i.throw("_module.goose not found in directory %s", name)
			}

			i.throw(err.Error())
		}

		name = newName
	}

	content, err := os.ReadFile(name)
	if err != nil {
		i.throw(err.Error())
	}

	astFile, err := parser.ParseFile(i.fset, name, content, nil)
	if err != nil {
		i.throw(err.Error())
	}

	module := &Module{
		File:    astFile,
		Exports: make(map[string]*Variable),
		Scope:   i.global.Fork(ScopeOwnerModule),
	}

	module.Scope.module = module
	i.modules[name] = module
	i.runModule(i.modules[name])

	return i.modules[name]
}

func (i *interp) loadPackageModule(specifier string) *Module {
	slash := strings.Index(specifier, "/")
	var name string
	var path string

	if slash == -1 {
		name = specifier
		path = ""
	} else {
		name = specifier[:slash]
		path = specifier[slash+1:]
	}

	var dir string
	if name == "std" {
		dir = filepath.Join(i.gooseRoot, "std")
	} else {
		dir = filepath.Join(i.gooseRoot, "pkg", name)
	}

	return i.loadFileModule(path, dir)
}

func isFilePath(path string) bool {
	return path == "." ||
		path == ".." ||
		strings.HasPrefix(path, "./") ||
		strings.HasPrefix(path, "../") ||
		strings.HasPrefix(path, "/")
}
