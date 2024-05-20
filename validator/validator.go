package validator

import (
	"errors"
	"fmt"
	"io"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/calico32/goose/ast"
	"github.com/calico32/goose/interpreter"
	. "github.com/calico32/goose/interpreter/lib"
	"github.com/calico32/goose/lib"
	"github.com/calico32/goose/parser"
	"github.com/calico32/goose/token"
	"go.lsp.dev/protocol"
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

type Diagnostic struct {
	Module   *Module
	Node     ast.Node
	Severity protocol.DiagnosticSeverity
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

func (v *Validator) Report(severity protocol.DiagnosticSeverity, node ast.Node, message string, parts ...any) {
	v.diagnostics = append(v.diagnostics, &Diagnostic{
		Module:   v.CurrentModule(),
		Node:     node,
		Severity: severity,
		Message:  fmt.Sprintf(message, parts...),
	})
}

func (v *Validator) Throw(msg string, parts ...any) {
	// panic(fmt.Errorf("%s: Validation error: %s", v.fset.Position(v.currentNode().Pos()), fmt.Sprintf(msg, parts...)))
	err := fmt.Sprintf("Validation error: %s", fmt.Sprintf(msg, parts...))
	v.Report(protocol.DiagnosticSeverityError, v.currentNode(), "Internal error: %s", err)
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
			v.Report(protocol.DiagnosticSeverityError, stmt, "cannot return from top-level")
		case *Break:
			v.Report(protocol.DiagnosticSeverityError, stmt, "cannot break from top-level")
		case *Continue:
			v.Report(protocol.DiagnosticSeverityError, stmt, "cannot continue from top-level")
		}
	}

	v.moduleStack = v.moduleStack[:len(v.moduleStack)-1]
}

func (v *Validator) checkStmt(scope *Scope, stmt ast.Stmt) StmtResult {
	defer pop(push(v, stmt))
	switch stmt := stmt.(type) {
	case *ast.ExprStmt:
		val := v.checkExpr(scope, stmt.X)
		if fn, ok := stmt.X.(*ast.FuncExpr); ok && fn.Name != nil && fn.Receiver == nil {
			// special case: if the expression is a function expression with a name, mark it as a declaration
			return &Decl{
				Name:  fn.Name.Name,
				Value: scope.Get(fn.Name.Name).Value,
			}
		}
		return &LoneValue{Value: val}
	case *ast.ImportStmt:
		return v.checkImportStmt(scope, stmt)
	case *ast.ExportDeclStmt:
		return v.checkExportDeclStmt(scope, stmt)
	case *ast.ConstStmt:
		return v.checkConstStmt(scope, stmt)
	case *ast.LetStmt:
		return v.checkLetStmt(scope, stmt)
	case *ast.AssignStmt:
		return v.checkAssignStmt(scope, stmt)
	case ast.NativeStmt:
		return v.checkNativeStmt(scope, stmt)
	case *ast.StructStmt:
		return v.checkStructStmt(scope, stmt)
	case *ast.BranchStmt:
		return v.checkBranchStmt(scope, stmt)
	case *ast.ReturnStmt:
		return v.checkReturnStmt(scope, stmt)
	case *ast.IfStmt:
		return v.checkIfStmt(scope, stmt)
	case *ast.ForStmt:
		return v.checkForStmt(scope, stmt)
	case *ast.RepeatForeverStmt:
		return v.checkRepeatForeverStmt(scope, stmt)
	case *ast.RepeatWhileStmt:
		return v.checkRepeatWhileStmt(scope, stmt)
	case *ast.RepeatCountStmt:
		return v.checkRepeatCountStmt(scope, stmt)
	case *ast.IncDecStmt:
		return v.checkIncDecStmt(scope, stmt)
	case *ast.OperatorStmt:
		return v.checkOperatorStmt(scope, stmt)
	default:
		fmt.Fprintf(os.Stderr, "unhandled statement type: %T\n", stmt)
		return &Void{}
	}
}

func (v *Validator) checkOperatorStmt(scope *Scope, stmt *ast.OperatorStmt) StmtResult {
	defer pop(push(v, stmt))

	constructor := scope.Get(stmt.Receiver.Name)
	if constructor == nil {
		v.Report(protocol.DiagnosticSeverityError, stmt.Receiver, "undefined type %s", stmt.Receiver.Name)
	} else if c, ok := constructor.Value.(*Func); !ok {
		v.Report(protocol.DiagnosticSeverityError, stmt.Receiver, "value %s is not a type", stmt.Receiver.Name)
	} else if c.NewableProto == nil {
		v.Report(protocol.DiagnosticSeverityError, stmt.Receiver, "value %s is not a type", stmt.Receiver.Name)
	} else {
		proto := c.NewableProto

		if proto.Operators[stmt.Tok] != nil {
			r := &ast.PosRange{From: stmt.TokPos, To: stmt.TokPos + token.Pos(len(stmt.Tok.String()))}
			v.Report(protocol.DiagnosticSeverityError, r, "duplicate operator %s", stmt.Tok)
		}

		proto.Operators[stmt.Tok] = &OperatorFunc{
			Async:   stmt.Async.IsValid(),
			Builtin: false,
		}
	}

	// validate parameters
	paramNames := map[string]bool{}
	for _, param := range stmt.Params.List {
		if paramNames[param.Ident.Name] {
			v.Report(protocol.DiagnosticSeverityError, param.Ident, "duplicate parameter %s", param.Ident.Name)
		}
		paramNames[param.Ident.Name] = true
		if param.Value != nil {
			v.checkExpr(scope, param.Value)
		}
	}

	closure := scope.Fork(ScopeOwnerClosure)
	funcScope := closure.Fork(ScopeOwnerFunc)

	// set parameters in scope
	for _, param := range stmt.Params.List {
		funcScope.Set(param.Ident.Name, &Variable{
			Constant: false,
		})
	}

	// TODO: better this
	funcScope.Set("this", &Variable{
		Constant: true,
	})

	if stmt.Arrow.IsValid() {
		v.checkExpr(funcScope, stmt.ArrowExpr)
	}

	v.checkStmts(funcScope, stmt.Body)

	return &Void{}
}

func (v *Validator) checkIncDecStmt(scope *Scope, stmt *ast.IncDecStmt) StmtResult {
	defer pop(push(v, stmt))

	v.checkExpr(scope, stmt.X)

	return &Void{}
}

func (v *Validator) checkRepeatForeverStmt(scope *Scope, stmt *ast.RepeatForeverStmt) StmtResult {
	defer pop(push(v, stmt))

	loopScope := scope.Fork(ScopeOwnerRepeat)
	v.checkStmts(loopScope, stmt.Body)
	return &Void{}
}

func (v *Validator) checkRepeatWhileStmt(scope *Scope, stmt *ast.RepeatWhileStmt) StmtResult {
	defer pop(push(v, stmt))

	loopScope := scope.Fork(ScopeOwnerRepeat)
	v.checkExpr(loopScope, stmt.Cond)
	v.checkStmts(loopScope, stmt.Body)
	return &Void{}
}

func (v *Validator) checkRepeatCountStmt(scope *Scope, stmt *ast.RepeatCountStmt) StmtResult {
	defer pop(push(v, stmt))

	loopScope := scope.Fork(ScopeOwnerRepeat)
	v.checkExpr(loopScope, stmt.Count)
	v.checkStmts(loopScope, stmt.Body)
	return &Void{}
}

func (v *Validator) checkReturnStmt(scope *Scope, stmt *ast.ReturnStmt) StmtResult {
	defer pop(push(v, stmt))

	// find the nearest function scope
	var funcScope *Scope
	for funcScope = scope; funcScope != nil; funcScope = funcScope.Parent() {
		if funcScope.Owner() == ScopeOwnerFunc {
			break
		}
	}

	if funcScope == nil {
		v.Report(protocol.DiagnosticSeverityError, stmt, "return outside of function")
	}

	return &Return{}
}

func (v *Validator) checkIfStmt(scope *Scope, stmt *ast.IfStmt) StmtResult {
	defer pop(push(v, stmt))

	v.checkExpr(scope, stmt.Cond)

	bodyScope := scope.Fork(ScopeOwnerIf)
	v.checkStmts(bodyScope, stmt.Body)

	if stmt.Else != nil {
		elseScope := scope.Fork(ScopeOwnerIf)
		v.checkStmts(elseScope, stmt.Else)
	}

	return &Void{}
}

func (v *Validator) checkBranchStmt(scope *Scope, stmt *ast.BranchStmt) StmtResult {
	defer pop(push(v, stmt))

	// find the nearest loop scope
	var loopScope *Scope
	for loopScope = scope; loopScope != nil; loopScope = loopScope.Parent() {
		if loopScope.Owner() == ScopeOwnerFor || loopScope.Owner() == ScopeOwnerRepeat {
			break
		}
	}

	if loopScope == nil {
		v.Report(protocol.DiagnosticSeverityError, stmt, "break/continue must be inside a loop")
	}

	return &Break{}
}

func (v *Validator) checkStructStmt(scope *Scope, stmt *ast.StructStmt) StmtResult {
	defer pop(push(v, stmt))

	if scope.IsDefinedInCurrentScope(stmt.Name.Name) {
		v.Report(protocol.DiagnosticSeverityError, stmt.Name, "cannot redefine struct %s", stmt.Name.Name)
	}

	// validate parameters
	fieldNames := map[string]bool{}
	for _, param := range stmt.Fields.List {
		if fieldNames[param.Ident.Name] {
			v.Report(protocol.DiagnosticSeverityError, param, "duplicate field %s", param.Ident.Name)

		}
		fieldNames[param.Ident.Name] = true
	}

	for _, field := range stmt.Fields.List {
		if field.Value != nil {
			v.checkExpr(scope, field.Value)
		}
	}

	proto := NewComposite()
	proto.Name = stmt.Name.Name

	closure := scope.Fork(ScopeOwnerClosure)

	// create new composite for testing
	obj := &Composite{
		Proto:      proto,
		Properties: make(Properties),
		Operators:  make(Operators),
	}

	if stmt.Init != nil {
		newScope := closure.Fork(ScopeOwnerStruct)

		// set parameters in scope
		for _, param := range stmt.Fields.List {
			if stmt.Init != nil {
				newScope.Set(param.Ident.Name, &Variable{
					Constant: false,
				})
			}

			if set, ok := obj.Properties[PKString]; ok {
				set[param.Ident.Name] = NullValue
			} else {
				obj.Properties[PKString] = map[string]Value{
					param.Ident.Name: NullValue,
				}
			}
		}

		// set this
		// TODO: better this
		newScope.Set("this", &Variable{
			Constant: false,
			Value:    obj,
		})

		v.checkStmts(newScope, stmt.Init.Body)
	}

	value := &Func{
		NewableProto: proto,
	}

	scope.Set(stmt.Name.Name, &Variable{
		Constant: false,
		Value:    value,
	})

	return &Decl{
		Name:  stmt.Name.Name,
		Value: value,
	}
}

func (v *Validator) checkForStmt(scope *Scope, stmt *ast.ForStmt) StmtResult {
	defer pop(push(v, stmt))

	v.checkExpr(scope, stmt.Iterable)

	loopScope := scope.Fork(ScopeOwnerFor)

	loopScope.Set(stmt.Var.Name, &Variable{
		Constant: false,
	})

	v.checkStmts(loopScope, stmt.Body)

	return &Void{}
}

func (v *Validator) checkStmts(scope *Scope, body []ast.Stmt) StmtResult {
	for _, stmt := range body {
		v.checkStmt(scope, stmt)
	}

	return &Void{}
}

func (v *Validator) checkExportDeclStmt(scope *Scope, stmt *ast.ExportDeclStmt) StmtResult {
	defer pop(push(v, stmt))
	if scope != scope.ModuleScope() {
		v.Report(protocol.DiagnosticSeverityError, stmt, "export declarations must be at the top level")
	}

	result := v.checkStmt(scope, stmt.Stmt)

	if scope != scope.ModuleScope() {
		v.Report(protocol.DiagnosticSeverityError, stmt, "export declarations must be at the top level")
	}

	if decl, ok := result.(*Decl); !ok {
		v.Report(protocol.DiagnosticSeverityError, stmt, "declaration expected")
	} else if _, ok := scope.Module().Exports[decl.Name]; ok {
		v.Report(protocol.DiagnosticSeverityError, stmt, "duplicate export %s", decl.Name)
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
		v.Report(protocol.DiagnosticSeverityError, stmt, "import declarations must be at the top level")
	}

	return v.checkImportSpec(scope, stmt.Spec)
}

func (v *Validator) checkImportSpec(scope *Scope, spec ast.ModuleSpec) StmtResult {
	defer pop(push(v, spec))
	name, module := v.loadModule(spec.ModuleSpecifier(), scope)
	if module == nil {
		return &Void{}
	}

	switch spec := spec.(type) {
	case *ast.ModuleSpecShow:
		defer pop(push(v, spec.Show))
		if spec.Show.Ellipsis.IsValid() {
			for name, value := range module.Exports {
				scope.Set(name, value)
				if scope.IsDefinedInCurrentScope(name) {
					v.Report(protocol.DiagnosticSeverityError, spec.Show, "name %s is already defined", name)
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
						v.Report(protocol.DiagnosticSeverityError, field.Ident, "name %s is already defined", field.Ident.Name)
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
							v.Report(protocol.DiagnosticSeverityError, field, "value %s is defined locally in module %s but is not exported", exportedName, name)
						} else {
							v.Report(protocol.DiagnosticSeverityError, field, "undefined export %s", exportedName)
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
			v.Report(protocol.DiagnosticSeverityError, spec, "name %s is already defined", moduleName)
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
		module = v.loadFileModule(name, dir, false)
	case "pkg":
		module = v.loadPackageModule(scheme + ":" + name)
	case "std":
		module = v.loadStdModule(scheme + ":" + name)
	default:
		v.Throw("unknown import scheme %s", scheme)
	}
	return scheme + ":" + name, module
}

func (v *Validator) loadFileModule(specifier string, dir string, isPackage bool) *Module {
	if module, ok := v.modules[specifier]; ok {
		return module
	}

	path := strings.TrimPrefix(specifier, "file:")

	if !filepath.IsAbs(path) {
		path = filepath.Join(dir, path)
	}

	info, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if isPackage {
				s := strings.TrimPrefix(dir, v.gooseRoot+"/pkg/")
				if specifier == "" {
					v.Report(protocol.DiagnosticSeverityError, v.currentNode(), "package %s not found", s)
				} else {
					v.Report(protocol.DiagnosticSeverityError, v.currentNode(), "module %s not found in package %s", s+"/"+specifier, s)
				}
			} else {
				v.Report(protocol.DiagnosticSeverityError, v.currentNode(), "module %s not found", specifier)
			}
			return nil
		}

		v.Throw(err.Error())
		return nil
	}

	if info == nil {
		v.Report(protocol.DiagnosticSeverityError, v.currentNode(), "module %s not found", specifier)
		return nil
	}

	if info.IsDir() {
		newPath := filepath.Join(path, "index.goose")
		_, err = os.Stat(newPath)
		if err != nil {
			if os.IsNotExist(err) {
				v.Report(protocol.DiagnosticSeverityError, v.currentNode(), "index.goose not found in directory %s", specifier)
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
	specifier = strings.TrimPrefix(specifier, "pkg:")
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

	return v.loadFileModule(path, dir, true)
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
			v.Report(protocol.DiagnosticSeverityError, v.currentNode(), "module %s or %s/index.goose not found", specifier, specifier)
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
		return v.checkIdent(scope, expr)
	case *ast.FuncExpr:
		return v.checkFuncExpr(scope, expr)
	case *ast.FrozenExpr:
		return v.checkFrozenExpr(scope, expr)
	case *ast.Literal:
		return v.checkLiteral(scope, expr)
	case *ast.StringLiteral:
		return v.checkStringLiteral(scope, expr)
	case *ast.BinaryExpr:
		return v.checkBinaryExpr(scope, expr)
	case *ast.UnaryExpr:
		return v.checkUnaryExpr(scope, expr)
	case *ast.CallExpr:
		return v.checkCallExpr(scope, expr)
	case *ast.ParenExpr:
		return v.checkExpr(scope, expr.X)
	case *ast.SelectorExpr:
		return v.checkSelectorExpr(scope, expr)
	case *ast.ArrayLiteral:
		return v.checkArrayLiteral(scope, expr)
	case *ast.BracketSelectorExpr:
		return v.checkBracketSelectorExpr(scope, expr)
	case *ast.CompositeLiteral:
		return v.checkCompositeLiteral(scope, expr)
	case *ast.SliceExpr:
		return v.checkSliceExpr(scope, expr)
	case *ast.IfExpr:
		return v.checkIfExpr(scope, expr)
	case *ast.MatchExpr:
		return v.checkMatchExpr(scope, expr)
	case *ast.DoExpr:
		return v.checkDoExpr(scope, expr)
	case *ast.NativeExpr:
		return v.checkNativeExpr(scope, expr)
	case *ast.RangeExpr:
		return v.checkRangeExpr(scope, expr)
	default:
		fmt.Fprintf(os.Stderr, "unhandled expression type: %T\n", expr)
		return nil
	}
}

func (v *Validator) checkRangeExpr(scope *Scope, expr *ast.RangeExpr) Value {
	defer pop(push(v, expr))

	v.checkExpr(scope, expr.Start)
	v.checkExpr(scope, expr.Stop)
	v.checkExpr(scope, expr.Step)

	return nil
}

func (v *Validator) checkNativeExpr(scope *Scope, expr *ast.NativeExpr) Value {
	defer pop(push(v, expr))

	module := scope.Module()
	natives, ok := interpreter.Natives[module.Specifier]
	if !ok {
		v.Report(protocol.DiagnosticSeverityError, expr, "module %s has no native components", module.Specifier)
		return nil
	}

	_, ok = natives[expr.Id]
	if !ok {
		v.Report(protocol.DiagnosticSeverityError, expr, "native %s not found", expr.Id)
		return nil
	}

	return nil
}

func (v *Validator) checkDoExpr(scope *Scope, expr *ast.DoExpr) Value {
	defer pop(push(v, expr))

	doScope := scope.Fork(ScopeOwnerDo)
	v.checkStmts(doScope, expr.Body)

	return nil
}

func (v *Validator) checkMatchExpr(scope *Scope, expr *ast.MatchExpr) Value {
	defer pop(push(v, expr))

	v.checkExpr(scope, expr.Expr)

	for _, branch := range expr.Clauses {
		switch branch := branch.(type) {
		case *ast.MatchPattern:
			// TODO: check pattern
			// v.checkPatternExpr(scope, branch.Pattern)
			v.checkExpr(scope, branch.Expr)
		case *ast.MatchElse:
			v.checkExpr(scope, branch.Expr)
		default:
			v.Throw("unhandled match clause type: %T", branch)
		}
	}

	return nil

}

func (v *Validator) checkBracketSelectorExpr(scope *Scope, expr *ast.BracketSelectorExpr) Value {
	defer pop(push(v, expr))

	v.checkExpr(scope, expr.X)
	v.checkExpr(scope, expr.Sel)
	v.checkPropertyKey(expr.Sel)

	return nil
}

func (v *Validator) checkCompositeLiteral(scope *Scope, lit *ast.CompositeLiteral) Value {
	defer pop(push(v, lit))

	composite := NewComposite()

	for _, field := range lit.Fields {
		// var keyValue PropertyKey
		switch key := field.Key.(type) {
		// case *ast.Ident:
		// keyValue = Wrap(key.Name).(PropertyKey)
		case *ast.StringLiteral:
			v.checkStringLiteral(scope, key)
		default:
			v.checkExpr(scope, key)

			// if int, ok := lit.(*Integer); ok {
			// 	keyValue = Wrap(int.Value).(PropertyKey)
			// 	break
			// }

			// i.Throw("invalid composite literal key type %s", lit.Type())
		}

		v.checkExpr(scope, field.Value)

		// SetProperty(composite, keyValue, val)
	}

	return composite
}

func (v *Validator) checkArrayLiteral(scope *Scope, lit *ast.ArrayLiteral) Value {
	defer pop(push(v, lit))

	for _, elem := range lit.List {
		v.checkExpr(scope, elem)
	}

	return nil
}

func (v *Validator) checkSliceExpr(scope *Scope, expr *ast.SliceExpr) Value {
	defer pop(push(v, expr))

	v.checkExpr(scope, expr.X)

	if expr.Low != nil {
		v.checkExpr(scope, expr.Low)
	}

	if expr.High != nil {
		v.checkExpr(scope, expr.High)
	}

	if expr.Low == nil && expr.High == nil {
		v.Report(protocol.DiagnosticSeverityError, expr, "slice expression must have at least one bound")
	}

	return nil
}

func (v *Validator) checkIfExpr(scope *Scope, expr *ast.IfExpr) Value {
	defer pop(push(v, expr))

	v.checkExpr(scope, expr.Cond)
	v.checkExpr(scope, expr.Then)
	v.checkExpr(scope, expr.Else)

	return nil
}

func (v *Validator) checkSelectorExpr(scope *Scope, expr *ast.SelectorExpr) Value {
	defer pop(push(v, expr))

	v.checkExpr(scope, expr.X)
	return nil
}

func (v *Validator) checkCallExpr(scope *Scope, expr *ast.CallExpr) Value {
	defer pop(push(v, expr))

	v.checkExpr(scope, expr.Func)
	for _, arg := range expr.Args {
		v.checkExpr(scope, arg)
	}

	return nil
}

func (v *Validator) checkBinaryExpr(scope *Scope, expr *ast.BinaryExpr) Value {
	defer pop(push(v, expr))

	v.checkExpr(scope, expr.X)

	if expr.Op == token.Arrow {
		rightScope := scope.Fork(ScopeOwnerPipeline)
		rightScope.Set("_", &Variable{
			Constant: true,
		})

		v.checkExpr(rightScope, expr.Y)
	} else {
		v.checkExpr(scope, expr.Y)
	}

	return nil
}

func (v *Validator) checkUnaryExpr(scope *Scope, expr *ast.UnaryExpr) Value {
	defer pop(push(v, expr))

	v.checkExpr(scope, expr.X)
	return nil
}

func (v *Validator) checkLiteral(_ *Scope, expr *ast.Literal) Value {
	switch expr.Kind {
	case token.Int:
		strVal := strings.ReplaceAll(expr.Value, "_", "")
		base := 10
		switch {
		case strings.HasPrefix(expr.Value, "0x"):
			strVal = strVal[2:]
			base = 16
		case strings.HasPrefix(expr.Value, "0o"):
			strVal = strVal[2:]
			base = 8
		case strings.HasPrefix(expr.Value, "0b"):
			strVal = strVal[2:]
			base = 2
		}
		val := new(big.Int)
		val, ok := val.SetString(strVal, base)
		if !ok {
			v.Report(protocol.DiagnosticSeverityError, expr, "failed to parse integer")
		}

		return Wrap(val)

	case token.Float:
		val, err := strconv.ParseFloat(expr.Value, 64)
		if err != nil {
			v.Report(protocol.DiagnosticSeverityError, expr, err.Error())
			return nil
		}

		return Wrap(val)

	case token.Null:
		return NullValue

	case token.True:
		return TrueValue

	case token.False:
		return FalseValue

	default:
		v.Report(protocol.DiagnosticSeverityError, expr, "unexpected literal kind %s", expr.Kind)
		return nil
	}
}

func (v *Validator) checkStringLiteral(scope *Scope, lit *ast.StringLiteral) Value {
	defer pop(push(v, lit))

	// TODO: check start and end parts

	for _, part := range lit.Parts {
		switch part := part.(type) {
		case *ast.StringLiteralInterpExpr:
			v.checkExpr(scope, part.Expr)
			continue
		case *ast.StringLiteralInterpIdent:
			v.checkIdent(scope, part.Ident)
			continue
		case *ast.StringLiteralMiddle:
			// TODO: check for escape sequences
		default:
			v.Throw("unhandled string literal part type: %T", part)
		}
	}

	return &String{}
}

func (v *Validator) checkIdent(scope *Scope, ident *ast.Ident) Value {
	defer pop(push(v, ident))
	if ident.Name[0] == '#' {
		// ignore
		return nil
	}

	if !scope.IsDefined(ident.Name) {
		v.Report(protocol.DiagnosticSeverityError, ident, "%s is not defined", ident.Name)
	}
	return nil
}

func (v *Validator) checkConstStmt(scope *Scope, stmt *ast.ConstStmt) StmtResult {
	defer pop(push(v, stmt))

	if stmt.Ident.Name == "_" {
		v.Report(protocol.DiagnosticSeverityError, stmt.Ident, "cannot declare _")
	}

	if scope.IsDefinedInCurrentScope(stmt.Ident.Name) {
		v.Report(protocol.DiagnosticSeverityError, stmt.Ident, "cannot redefine variable %s", stmt.Ident.Name)
	}

	value := v.checkExpr(scope, stmt.Value)

	scope.Set(stmt.Ident.Name, &Variable{
		Constant: true,
		Value:    value,
	})

	return &Decl{
		Name:  stmt.Ident.Name,
		Value: value,
	}
}

func (v *Validator) checkLetStmt(scope *Scope, stmt *ast.LetStmt) StmtResult {
	defer pop(push(v, stmt))

	if stmt.Ident.Name == "_" {
		v.Report(protocol.DiagnosticSeverityError, stmt.Ident, "cannot declare _")
	}

	if scope.IsDefinedInCurrentScope(stmt.Ident.Name) {
		v.Report(protocol.DiagnosticSeverityError, stmt.Ident, "cannot redefine variable %s", stmt.Ident.Name)
	}

	var value Value

	if stmt.Value != nil {
		value = v.checkExpr(scope, stmt.Value)
	}

	if value == nil {
		value = NullValue
	}

	scope.Set(stmt.Ident.Name, &Variable{
		Constant: false,
		Value:    value,
	})

	return &Decl{
		Name:  stmt.Ident.Name,
		Value: value,
	}
}

func (v *Validator) checkAssignStmt(scope *Scope, stmt *ast.AssignStmt) StmtResult {
	defer pop(push(v, stmt))

	switch lhs := stmt.Lhs.(type) {
	case *ast.Ident:
		ident := lhs.Name

		if ident[0] == '#' {
			x := scope.Get("this")
			if x == nil {
				v.Report(protocol.DiagnosticSeverityError, lhs, "invalid property assignment: 'this' is not defined")
			}
			v.checkExpr(scope, stmt.Rhs)

			return &Void{}
		}

		existing := scope.Get(ident)
		if existing == nil {
			v.Report(protocol.DiagnosticSeverityError, lhs, "%s is not defined", ident)
		} else if existing.Constant {
			v.Report(protocol.DiagnosticSeverityError, lhs, "cannot assign to constant %s", ident)
		}

		// op := GetOperator(existing.Value, stmt.Tok)
		// if op == nil {
		// 	node := &ast.PosRange{From: stmt.TokPos, To: stmt.TokPos + token.Pos(len(stmt.Tok.String()))}
		// 	v.Report(protocol.DiagnosticSeverityError, node, "operator %s not defined for type %s", stmt.Tok, existing.Value.Type())
		// }
		v.checkExpr(scope, stmt.Rhs)
		return &Void{}
	case *ast.SelectorExpr:
		v.checkExpr(scope, lhs.X)
		v.checkExpr(scope, stmt.Rhs)
	case *ast.BracketSelectorExpr:
		v.checkExpr(scope, lhs.X)
		v.checkExpr(scope, lhs.Sel)
		v.checkPropertyKey(lhs.Sel)
		v.checkExpr(scope, stmt.Rhs)
	}

	return &Void{}
}

func (v *Validator) checkPropertyKey(expr ast.Expr) {
	switch expr := expr.(type) {
	case *ast.StringLiteral, *ast.Literal:
		// valid key
	case *ast.Ident:
		if expr.Name == "true" || expr.Name == "false" || expr.Name == "null" {
			v.Report(protocol.DiagnosticSeverityError, expr, "invalid key %s", expr.Name)
		}
		// otherwise, we have no way of knowing if this is a valid key
	case *ast.CompositeLiteral, *ast.FrozenExpr, *ast.BindExpr, *ast.ArrayLiteral, *ast.DoExpr, *ast.FuncExpr, *ast.NativeExpr, *ast.ArrayInitializer, *ast.EllipsisExpr, *ast.GeneratorExpr:
		// TODO: exhaustive list of invalid keys
		v.Report(protocol.DiagnosticSeverityError, expr, "invalid key %s", expr)
	default:
		// we have no way of knowing if this is a valid key
	}
}

func (v *Validator) checkFrozenExpr(scope *Scope, expr *ast.FrozenExpr) Value {
	defer pop(push(v, expr))

	if _, ok := expr.X.(*ast.CompositeLiteral); !ok {
		v.Report(protocol.DiagnosticSeverityError, expr, "frozen keyword can only be used with composite literals")
	}

	return v.checkExpr(scope, expr.X)
}
func (v *Validator) checkNativeStmt(scope *Scope, stmt ast.NativeStmt) StmtResult {
	defer pop(push(v, stmt))

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
		v.Report(protocol.DiagnosticSeverityError, stmt, "invalid native stmt type %T", stmt)
	}

	specifier := scope.Module().Specifier

	if moduleNatives, ok := interpreter.Natives[specifier]; ok {
		if value, ok := moduleNatives[name]; ok {
			if fn, ok := stmt.(*ast.NativeFunc); ok && fn.Receiver != nil {
				// find proto
				// TODO: limit to current module
				constructor := scope.Get(fn.Receiver.Name)
				if constructor == nil {
					v.Report(protocol.DiagnosticSeverityError, stmt, "unknown type %s", fn.Receiver.Name)
					return nil
				} else if val, ok := constructor.Value.(*Func); !ok || val.NewableProto == nil {
					v.Report(protocol.DiagnosticSeverityError, stmt, "%s cannot have receiver functions", fn.Receiver.Name)
					return nil
				}

				proto := constructor.Value.(*Func).NewableProto

				if proto.Properties[PKString] == nil {
					proto.Properties[PKString] = make(map[string]Value)
				}

				if _, ok := proto.Properties[PKString][fn.Name.Name]; ok {
					v.Report(protocol.DiagnosticSeverityError, stmt, "duplicate receiver function %s", fn.Name.Name)
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

	v.Report(protocol.DiagnosticSeverityError, stmt, "native %s not found", name)
	return nil
}

func (v *Validator) checkFuncExpr(scope *Scope, expr *ast.FuncExpr) Value {
	defer pop(push(v, expr))

	if expr.Name != nil && expr.Receiver == nil {
		if scope.IsDefinedInCurrentScope(expr.Name.Name) {
			v.Report(protocol.DiagnosticSeverityError, expr.Name, "cannot redefine function %s", expr.Name.Name)
		}
	}

	// validate parameters
	paramNames := map[string]bool{}
	for _, param := range expr.Params.List {
		if paramNames[param.Ident.Name] {
			v.Report(protocol.DiagnosticSeverityError, param.Ident, "duplicate parameter %s", param.Ident.Name)
		}
		paramNames[param.Ident.Name] = true
	}

	for _, param := range expr.Params.List {
		if param.Value != nil {
			v.checkExpr(scope, param.Value)
		}
	}

	closure := scope.Fork(ScopeOwnerClosure)
	funcScope := closure.Fork(ScopeOwnerFunc)

	// set parameters in scope
	for _, param := range expr.Params.List {
		funcScope.Set(param.Ident.Name, &Variable{
			Constant: false,
		})
	}

	// TODO: better this
	funcScope.Set("this", &Variable{
		Constant: true,
	})

	if expr.Arrow.IsValid() {
		v.checkExpr(funcScope, expr.ArrowExpr)
	}

	v.checkStmts(funcScope, expr.Body)

	value := &Func{
		Async:    expr.Async.IsValid(),
		Memoized: expr.Memo.IsValid(),
	}

	if expr.Name != nil {
		if expr.Receiver != nil {
			// find proto
			// TODO: limit to current module
			constructor := scope.Get(expr.Receiver.Name)
			exit := false

			if constructor == nil {
				v.Report(protocol.DiagnosticSeverityError, expr, "unknown type %s", expr.Receiver.Name)
				exit = true
			} else if val, ok := constructor.Value.(*Func); !ok || val == nil || val.NewableProto == nil {
				v.Report(protocol.DiagnosticSeverityError, expr, "%s cannot have receiver functions", expr.Receiver.Name)
				exit = true
			}

			if exit {
				return nil
			}

			proto := constructor.Value.(*Func).NewableProto

			if proto.Properties[PKString] == nil {
				proto.Properties[PKString] = make(map[string]Value)
			}

			if _, ok := proto.Properties[PKString][expr.Name.Name]; ok {
				v.Report(protocol.DiagnosticSeverityError, expr.Name, "duplicate receiver function %s", expr.Name.Name)
			}

			proto.Properties[PKString][expr.Name.Name] = value
		} else {
			scope.Set(expr.Name.Name, &Variable{
				Constant: true, // functions are constants
				Value:    value,
			})
		}
	}

	return value
}
