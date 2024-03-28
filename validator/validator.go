package validator

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/calico32/goose/ast"
	"github.com/calico32/goose/interpreter"
	. "github.com/calico32/goose/interpreter/lib"
	"github.com/calico32/goose/lib"
	"github.com/calico32/goose/parser"
	"github.com/calico32/goose/token"
)

type Validator struct {
	fset        *token.FileSet
	moduleStack []*Module
	modules     map[string]*Module
	global      *Scope
	stdin       io.Reader
	stdout      io.Writer
	stderr      io.Writer
	gooseRoot   string
	diagnostics []*Diagnostic

	// internal state
	trace    bool
	indent   int
	lastNode ast.Node
	stack    []ast.Node
}

type DiagnosticSeverity int

func (s DiagnosticSeverity) String() string {
	return Severity[s]
}

const (
	Error DiagnosticSeverity = iota
	Warning
	Info
	Hint
)

var Severity = [...]string{
	"error",
	"warning",
	"info",
	"hint",
}

type Diagnostic struct {
	Module   *Module
	Node     ast.Node
	Severity DiagnosticSeverity
	Message  string
}

func (v *Validator) Fset() *token.FileSet        { return v.fset }
func (v *Validator) ImportStack() []*Module      { return v.moduleStack }
func (v *Validator) Diagnostics() []*Diagnostic  { return v.diagnostics }
func (v *Validator) Modules() map[string]*Module { return v.modules }
func (v *Validator) Global() *Scope              { return v.global }
func (v *Validator) Stdin() io.Reader            { return v.stdin }
func (v *Validator) Stdout() io.Writer           { return v.stdout }
func (v *Validator) Stderr() io.Writer           { return v.stderr }
func (v *Validator) GooseRoot() string           { return v.gooseRoot }

func (v *Validator) CurrentModule() *Module {
	if len(v.moduleStack) == 0 {
		if len(v.modules) == 1 {
			for _, m := range v.modules {
				v.moduleStack = append(v.moduleStack, m)
				return m
			}
		}
		v.Throw("no current module")
	}

	return v.moduleStack[len(v.moduleStack)-1]
}

func (v *Validator) Report(severity DiagnosticSeverity, node ast.Node, message string, parts ...any) {
	v.diagnostics = append(v.diagnostics, &Diagnostic{
		Module:   v.CurrentModule(),
		Node:     node,
		Severity: severity,
		Message:  fmt.Sprintf(message, parts...),
	})
}

func (v *Validator) Throw(msg string, parts ...any) {
	panic(fmt.Errorf("%s: Validation error: %s", v.fset.Position(v.currentNode().Pos()), fmt.Sprintf(msg, parts...)))
}

func (v *Validator) currentNode() ast.Node {
	if len(v.stack) == 0 {
		return v.lastNode
	}
	return v.stack[len(v.stack)-1]
}

func push(v *Validator, n ast.Node) *Validator {
	v.lastNode = n
	v.stack = append(v.stack, n)
	return v
}

func pop(v *Validator) {
	v.stack = v.stack[:len(v.stack)-1]
}

func New(file *ast.Module, fset *token.FileSet, trace bool, stdin io.Reader, stdout io.Writer, stderr io.Writer) (i *Validator, err error) {
	i = &Validator{
		modules:     make(map[string]*Module),
		global:      NewGlobalScope(interpreter.GlobalConstants),
		trace:       trace,
		fset:        fset,
		stdin:       stdin,
		stdout:      stdout,
		stderr:      stderr,
		gooseRoot:   os.Getenv("GOOSEROOT"),
		moduleStack: make([]*Module, 0, 10),
	}

	if i.gooseRoot == "" {
		if xdgDataHome := os.Getenv("XDG_DATA_HOME"); xdgDataHome != "" {

			i.gooseRoot = filepath.Join(xdgDataHome, "goose")
		} else {
			home := os.Getenv("HOME")
			if home == "" {
				home = os.Getenv("USERPROFILE")
			}

			i.gooseRoot = filepath.Join(home, ".goose")
		}
	}

	err = interpreter.CreateGooseRoot(i.gooseRoot)
	if err != nil {
		return
	}

	module := &Module{
		Module:  file,
		Scope:   i.global.Fork(ScopeOwnerModule),
		Exports: make(map[string]*Variable),
	}

	module.Scope.SetModule(module)
	i.modules[file.Specifier] = module
	i.moduleStack = append(i.moduleStack, module)

	return
}

func (v *Validator) Check() (exitCode int, err error) {
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		if exit, ok := r.(gooseExit); ok {
	// 			exitCode = int(exit.code)
	// 		} else {
	// 			exitCode = PANIC
	// 			err = fmt.Errorf("%s", r)
	// 		}
	// 	}
	// }()
	v.diagnostics = make([]*Diagnostic, 0, 10)
	for name, val := range interpreter.Globals {
		v.global.Set(name, &Variable{
			Constant: true,
			Value:    &Func{Executor: val},
		})
	}

	v.checkBuiltins()
	// if len(v.diagnostics) > 0 {
	// v.Throw("internal error: builtins failed validation")
	// }
	v.checkModule(v.CurrentModule())
	return 0, nil
}

func (v *Validator) checkBuiltins() {
	specs := []string{
		"std:language/builtin.goose",
	}

	for _, spec := range specs {
		mod := v.loadStdModule(spec)
		v.copyModuleExportsToGlobal(mod)
	}
}
func (v *Validator) checkModule(module *Module) {
	defer pop(push(v, module.Module))
	v.moduleStack = append(v.moduleStack, module)

	for _, stmt := range module.Stmts {
		result := v.checkStmt(module.Scope, stmt)
		switch result.(type) {
		case *Return:
			v.Report(Error, stmt, "cannot return from top-level")
		case *Break:
			v.Report(Error, stmt, "cannot break from top-level")
		case *Continue:
			v.Report(Error, stmt, "cannot continue from top-level")
		}
	}

	v.moduleStack = v.moduleStack[:len(v.moduleStack)-1]
}

func (v *Validator) checkStmt(scope *Scope, stmt ast.Stmt) StmtResult {
	defer pop(push(v, stmt))
	switch stmt := stmt.(type) {
	case *ast.ExprStmt:
		v.checkExpr(scope, stmt.X)
		return &Void{}
	case *ast.ImportStmt:
		return v.checkImportStmt(scope, stmt)
	case *ast.ExportDeclStmt:
		return v.checkExportDeclStmt(scope, stmt)
	default:
		fmt.Printf("unhandled statement type: %T\n", stmt)
		return &Void{}
	}
}

func (v *Validator) checkStmts(scope *Scope, body []ast.Stmt) StmtResult {
	var last StmtResult
	for _, stmt := range body {
		result := v.checkStmt(scope, stmt)
		switch result.(type) {
		case *Return, *Break, *Continue:
			return result
		}
		last = result
	}

	if val, ok := last.(*LoneValue); ok {
		return val
	}

	return &Void{}
}

func (v *Validator) checkExportDeclStmt(scope *Scope, stmt *ast.ExportDeclStmt) StmtResult {
	defer pop(push(v, stmt))
	if scope != scope.ModuleScope() {
		v.Report(Error, stmt, "export declarations must be at the top level")
	}

	result := v.checkStmt(scope, stmt.Stmt)

	if scope != scope.ModuleScope() {
		v.Report(Error, stmt, "export declarations must be at the top level")
	}

	if decl, ok := result.(*Decl); !ok {
		v.Report(Error, stmt, "declaration expected")
	} else if _, ok := scope.Module().Exports[decl.Name]; ok {
		v.Report(Error, stmt, "duplicate export %s", decl.Name)
	} else {
		scope.Module().Exports[decl.Name] = &Variable{
			Source:   VariableSourceImport,
			Value:    decl.Value,
			Constant: true,
		}
	}

	return &Void{}
}

func (v *Validator) checkImportStmt(scope *Scope, stmt *ast.ImportStmt) StmtResult {
	defer pop(push(v, stmt))
	if scope != scope.ModuleScope() {
		v.Report(Error, stmt, "import declarations must be at the top level")
	}

	return v.checkImportSpec(scope, stmt.Spec)
}

func (v *Validator) checkImportSpec(scope *Scope, spec ast.ModuleSpec) StmtResult {
	defer pop(push(v, spec))
	name, module := v.loadModule(spec.ModuleSpecifier(), scope)

	switch spec := spec.(type) {
	case *ast.ModuleSpecShow:
		defer pop(push(v, spec.Show))
		if spec.Show.Ellipsis.IsValid() {
			for name, value := range module.Exports {
				scope.Set(name, value)
				if scope.IsDefinedInCurrentScope(name) {
					v.Report(Error, spec.Show, "name %s is already defined", name)
				}
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
					v.checkImportSpec(scope, spec)
				case *ast.ShowFieldEllipsis:
					if scope.IsDefinedInCurrentScope(field.Ident.Name) {
						v.Report(Error, field.Ident, "name %s is already defined", field.Ident.Name)
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

					switch field := field.(type) {
					case *ast.ShowFieldIdent:
						localName = field.Ident.Name
						exportedName = field.Ident.Name
					case *ast.ShowFieldAs:
						localName = field.Alias.Name
						exportedName = field.Ident.Name
					default:
						v.Throw("unhandled show field type: %T", field)
					}

					if _, ok := module.Exports[exportedName]; !ok {
						if module.Scope.IsDefinedInCurrentScope(exportedName) {
							v.Report(Error, field, "value %s is defined locally in module %s but is not exported", exportedName, name)
						} else {
							v.Report(Error, field, "undefined export %s", exportedName)
						}
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
				v.Throw(err.Error())
			}
		}

		if scope.IsDefinedInCurrentScope(moduleName) {
			v.Report(Error, spec, "name %s is already defined", moduleName)
		}

		scope.Set(moduleName, &Variable{
			Constant: true,
			Value:    object,
		})
	}

	return &Void{}
}

func (v *Validator) loadModule(specifier string, scope *Scope) (string, *Module) {
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
			v.Throw("invalid file import path %s", specifier)
		}
	} else {
		scheme = "pkg"
		name = strings.TrimPrefix(specifier, "pkg:")
		if isFilePath(specifier) {
			v.Throw("invalid package import path %s", specifier)
		}
	}

	if module, ok := v.modules[specifier]; ok {
		return specifier, module
	}

	var module *Module
	switch scheme {
	case "file":
		dir := strings.TrimPrefix(filepath.Dir(scope.Module().Specifier), "file:")
		module = v.loadFileModule(name, dir)
	case "pkg":
		module = v.loadPackageModule(scheme + ":" + name)
	case "std":
		module = v.loadStdModule(scheme + ":" + name)
	default:
		v.Throw("unknown import scheme %s", scheme)
	}
	return scheme + ":" + name, module
}

func (v *Validator) loadFileModule(specifier string, dir string) *Module {
	if module, ok := v.modules[specifier]; ok {
		return module
	}

	path := strings.TrimPrefix(specifier, "file:")

	if !filepath.IsAbs(path) {
		path = filepath.Join(dir, path)
	}

	info, err := os.Stat(path)
	if err != nil {
		v.Throw(err.Error())
	}

	if info.IsDir() {
		newPath := filepath.Join(path, "index.goose")
		_, err = os.Stat(newPath)
		if err != nil {
			if os.IsNotExist(err) {
				v.Report(Error, v.currentNode(), "index.goose not found in directory %s", specifier)
				return nil
			}

			v.Throw(err.Error())
		}

		path = newPath
	}

	content, err := os.ReadFile(path)
	if err != nil {
		v.Throw(err.Error())
	}

	astFile, err := parser.ParseFile(v.fset, specifier, content, nil)
	if err != nil {
		v.Throw(err.Error())
	}

	module := &Module{
		Module:  astFile,
		Exports: make(map[string]*Variable),
		Scope:   v.global.Fork(ScopeOwnerModule),
	}

	module.Scope.SetModule(module)
	v.modules[specifier] = module
	v.checkModule(module)

	return module
}

func (v *Validator) loadPackageModule(specifier string) *Module {
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

	dir := filepath.Join(v.gooseRoot, "pkg", name)

	return v.loadFileModule(path, dir)
}

func isFilePath(path string) bool {
	return path == "." ||
		path == ".." ||
		strings.HasPrefix(path, "./") ||
		strings.HasPrefix(path, "../") ||
		strings.HasPrefix(path, "/")
}

func (v *Validator) loadStdModule(specifier string) *Module {
	if module, ok := v.modules[specifier]; ok {
		return module
	}

	bindataPath := filepath.Join("std", strings.TrimPrefix(specifier, "std:"))
	content, err := lib.Stdlib.ReadFile(bindataPath)
	if err != nil {
		// try to find index.goose in the directory
		bindataPath = filepath.Join(bindataPath, "index.goose")
		content, err = lib.Stdlib.ReadFile(bindataPath)
		if err != nil {
			v.Report(Error, v.currentNode(), "module %s or %s/index.goose not found", specifier, specifier)
			return nil
		} else {
			specifier += "/index.goose"
		}
	}

	file, err := parser.ParseFile(v.fset, specifier, content, nil)
	if err != nil {
		v.Throw(err.Error())
	}

	module := &Module{
		Module:  file,
		Exports: make(map[string]*Variable),
		Scope:   v.global.Fork(ScopeOwnerModule),
	}

	module.Scope.SetModule(module)
	v.modules[specifier] = module
	v.checkModule(module)

	return module
}

func (i *Validator) copyModuleExportsToGlobal(module *Module) {
	for name, value := range module.Exports {
		i.global.Set(name, &Variable{
			Constant: true,
			Value:    value.Value,
		})
	}
}

func (v *Validator) checkExpr(scope *Scope, expr ast.Expr) Value {
	defer pop(push(v, expr))
	switch expr := expr.(type) {
	case *ast.Ident:
		v.checkIdent(scope, expr)
		return nil
	default:
		fmt.Printf("unhandled expression type: %T\n", expr)
		return nil
	}
}

func (v *Validator) checkIdent(scope *Scope, ident *ast.Ident) Value {
	defer pop(push(v, ident))
	if !scope.IsDefined(ident.Name) {
		v.Report(Error, ident, "%s is not defined", ident.Name)
	}
	return nil
}
