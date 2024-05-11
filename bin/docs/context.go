package main

import (
	"github.com/calico32/goose/lib/types"
	"github.com/labstack/echo/v4"
)

type Context interface {
	Clone() Context
	CloneFor(c echo.Context) Context
	Path() string
	StandardLibrary() []types.StdlibDoc
	Builtins() []types.BuiltinDoc
}

type TemplateContext struct {
	path            string
	standardLibrary []types.StdlibDoc
	builtins        []types.BuiltinDoc
}

func (ctx *TemplateContext) Path() string                       { return ctx.path }
func (ctx *TemplateContext) StandardLibrary() []types.StdlibDoc { return ctx.standardLibrary }
func (ctx *TemplateContext) Builtins() []types.BuiltinDoc       { return ctx.builtins }
func (ctx *TemplateContext) Clone() Context {
	return &TemplateContext{
		path:            ctx.path,
		standardLibrary: ctx.standardLibrary,
		builtins:        ctx.builtins,
	}
}
func (ctx *TemplateContext) CloneFor(c echo.Context) Context {
	new := ctx.Clone().(*TemplateContext)
	new.path = c.Request().URL.Path
	return new
}

type BuiltinContext struct {
	*TemplateContext
	Doc types.BuiltinDoc
}

func (ctx *BuiltinContext) Path() string                       { return ctx.path }
func (ctx *BuiltinContext) StandardLibrary() []types.StdlibDoc { return ctx.standardLibrary }
func (ctx *BuiltinContext) Builtins() []types.BuiltinDoc       { return ctx.builtins }
func (ctx *BuiltinContext) Clone() Context {
	return &BuiltinContext{
		TemplateContext: ctx.TemplateContext.Clone().(*TemplateContext),
		Doc:             ctx.Doc,
	}
}
func (ctx *BuiltinContext) CloneFor(c echo.Context) Context {
	new := ctx.Clone().(*BuiltinContext)
	new.path = c.Request().URL.Path
	return new
}

type StdlibContext struct {
	*TemplateContext
	Pkg types.StdlibDoc
}

func (ctx *StdlibContext) Path() string                       { return ctx.path }
func (ctx *StdlibContext) StandardLibrary() []types.StdlibDoc { return ctx.standardLibrary }
func (ctx *StdlibContext) Builtins() []types.BuiltinDoc       { return ctx.builtins }
func (ctx *StdlibContext) Clone() Context {
	return &StdlibContext{
		TemplateContext: ctx.TemplateContext.Clone().(*TemplateContext),
		Pkg:             ctx.Pkg,
	}
}
func (ctx *StdlibContext) CloneFor(c echo.Context) Context {
	new := ctx.Clone().(*StdlibContext)
	new.path = c.Request().URL.Path
	return new
}
