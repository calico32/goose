package interpreter

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/calico32/goose/ast"
	. "github.com/calico32/goose/interpreter/lib"
	"github.com/calico32/goose/lib"
	"github.com/calico32/goose/parser"
)

func (i *interp) runExportDeclStmt(scope *Scope, stmt *ast.ExportDeclStmt) StmtResult {
	result := i.runStmt(scope, stmt.Stmt)

	if scope != scope.ModuleScope() {
		i.Throw("export declarations must be at the top level")
	}

	if decl, ok := result.(*Decl); !ok {
		i.Throw("declaration expected")
	} else if _, ok := scope.Module().Exports[decl.Name]; ok {
		i.Throw("duplicate export %s", decl.Name)
	} else {
		scope.Module().Exports[decl.Name] = &Variable{
			Source:   VariableSourceImport,
			Value:    decl.Value,
			Constant: true,
		}
	}

	return &Void{}
}

func (i *interp) runExportListStmt(scope *Scope, stmt *ast.ExportListStmt) StmtResult {
	if scope != scope.ModuleScope() {
		i.Throw("export declarations must be at the top level")
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
			i.Throw("unknown export field type %T", field)
		}

		if _, ok := scope.Module().Exports[exportedName]; ok {
			i.Throw("duplicate export %s", exportedName)
		}

		if !scope.IsDefinedInCurrentScope(localName) {
			i.Throw("undefined name %s", localName)
		}

		scope.Module().Exports[exportedName] = scope.Get(localName)
	}

	return &Void{}
}

func (i *interp) runExportSpecStmt(scope *Scope, stmt *ast.ExportSpecStmt) StmtResult {
	if scope != scope.ModuleScope() {
		i.Throw("export declarations must be at the top level")
	}

	importScope := scope.Fork(ScopeOwnerImport)

	i.runImportSpec(importScope, stmt.Spec)

	for name := range importScope.Idents() {
		if _, ok := scope.Module().Exports[name]; ok {
			i.Throw("duplicate export %s", name)
		}

		scope.Module().Exports[name] = importScope.Get(name)
	}

	return &Void{}
}

func (i *interp) runImportStmt(scope *Scope, stmt *ast.ImportStmt) StmtResult {
	if scope != scope.ModuleScope() {
		i.Throw("import declarations must be at the top level")
	}

	return i.runImportSpec(scope, stmt.Spec)
}

func (i *interp) runImportSpec(scope *Scope, spec ast.ModuleSpec) StmtResult {
	// TODO: custom import schemes
	name, module := i.loadModule(spec.ModuleSpecifier(), scope)

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
						i.Throw("name %s is already defined", field.Ident.Name)
					}
					// put all remaining exports into an object
					object := NewComposite()
					for name, value := range module.Exports {
						if imported[name] {
							continue
						}
						SetProperty(object, NewString(name), value.Value)
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
						i.Throw("unknown import field type %T", field)
					}

					if _, ok := module.Exports[exportedName]; !ok {
						if module.Scope.IsDefinedInCurrentScope(exportedName) {
							i.Throw("value %s is defined locally in module %s but is not exported", exportedName, name)
						}
						i.Throw("undefined export %s", exportedName)
					}

					imported[exportedName] = true

					value := module.Exports[exportedName]
					scope.Set(localName, value)
				}
			}
		}
	default:
		object := NewComposite()
		for name, value := range module.Exports {
			SetProperty(object, NewString(name), value.Value) // TODO: reassignment can change the value
		}
		object.Frozen = true

		var moduleName string

		if aliased, ok := spec.(*ast.ModuleSpecAs); ok {
			moduleName = aliased.Alias.Name
		} else {
			var err error
			moduleName, err = ast.ModuleName(strings.TrimPrefix(name, module.Scheme+":"))
			if err != nil {
				i.Throw(err.Error())
			}
		}

		if scope.IsDefinedInCurrentScope(moduleName) {
			i.Throw("name %s is already defined", moduleName)
		}

		scope.Set(moduleName, &Variable{
			Constant: true,
			Value:    object,
		})
	}

	return &Void{}
}

func (i *interp) loadModule(specifier string, scope *Scope) (string, *Module) {
	var scheme string
	var name string
	if strings.Contains(specifier, ":") {
		colon := strings.Index(specifier, ":")
		scheme = specifier[:colon]
		name = specifier[colon+1:]
	} else if isFilePath(specifier) {
		// inherit the scheme from the parent module
		scheme = scope.Module().Scheme
		name = filepath.Join(filepath.Dir(strings.TrimPrefix(scope.Module().Specifier, scope.Module().Scheme+":")), specifier)
		if scheme == "file" && !isFilePath(specifier) {
			i.Throw("invalid file import path %s", specifier)
		}
	} else {
		scheme = "pkg"
		name = strings.TrimPrefix(specifier, "pkg:")
		if isFilePath(specifier) {
			i.Throw("invalid package import path %s", specifier)
		}
	}

	if module, ok := i.modules[specifier]; ok {
		return specifier, module
	}

	var module *Module
	switch scheme {
	case "file":
		dir := strings.TrimPrefix(filepath.Dir(scope.Module().Specifier), "file:")
		module = i.loadFileModule(name, dir)
	case "pkg":
		module = i.loadPackageModule(scheme + ":" + name)
	case "std":
		module = i.loadStdModule(scheme + ":" + name)
	default:
		i.Throw("unknown import scheme %s", scheme)
	}
	return scheme + ":" + name, module
}

func (i *interp) loadFileModule(specifier string, dir string) *Module {
	if module, ok := i.modules[specifier]; ok {
		return module
	}

	path := strings.TrimPrefix(specifier, "file:")

	if !filepath.IsAbs(path) {
		path = filepath.Join(dir, path)
	}

	info, err := os.Stat(path)
	if err != nil {
		i.Throw(err.Error())
	}

	if info.IsDir() {
		newPath := filepath.Join(path, "index.goose")
		_, err = os.Stat(newPath)
		if err != nil {
			if os.IsNotExist(err) {
				i.Throw("index.goose not found in directory %s", specifier)
			}

			i.Throw(err.Error())
		}

		path = newPath
	}

	content, err := os.ReadFile(path)
	if err != nil {
		i.Throw(err.Error())
	}

	astFile, err := parser.ParseFile(i.fset, specifier, content, nil)
	if err != nil {
		i.Throw(err.Error())
	}

	module := &Module{
		Module:  astFile,
		Exports: make(map[string]*Variable),
		Scope:   i.global.Fork(ScopeOwnerModule),
	}

	module.Scope.SetModule(module)
	i.modules[specifier] = module
	i.runModule(module)

	return module
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

	dir := filepath.Join(i.gooseRoot, "pkg", name)

	return i.loadFileModule(path, dir)
}

func isFilePath(path string) bool {
	return path == "." ||
		path == ".." ||
		strings.HasPrefix(path, "./") ||
		strings.HasPrefix(path, "../") ||
		strings.HasPrefix(path, "/")
}

func (i *interp) loadStdModule(specifier string) *Module {
	if module, ok := i.modules[specifier]; ok {
		return module
	}

	bindataPath := filepath.Join("std", strings.TrimPrefix(specifier, "std:"))
	content, err := lib.Stdlib.ReadFile(bindataPath)
	if err != nil {
		// try to find index.goose in the directory
		bindataPath = filepath.Join(bindataPath, "index.goose")
		content, err = lib.Stdlib.ReadFile(bindataPath)
		if err != nil {
			i.Throw("module %s or %s/index.goose not found", specifier, specifier)
		} else {
			specifier += "/index.goose"
		}
	}

	file, err := parser.ParseFile(i.fset, specifier, content, nil)
	if err != nil {
		i.Throw(err.Error())
	}

	module := &Module{
		Module:  file,
		Exports: make(map[string]*Variable),
		Scope:   i.global.Fork(ScopeOwnerModule),
	}

	module.Scope.SetModule(module)
	i.modules[specifier] = module
	i.runModule(module)

	return module
}

func (i *interp) copyModuleExportsToGlobal(module *Module) {
	for name, value := range module.Exports {
		i.global.Set(name, &Variable{
			Constant: true,
			Value:    value.Value,
		})
	}
}
