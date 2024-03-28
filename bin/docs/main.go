package main

import (
	"embed"
	"net/http"
	"strings"

	"github.com/calico32/goose/interpreter"
	"github.com/calico32/goose/lib"
	std_language "github.com/calico32/goose/lib/std/language"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

//go:embed static
var static embed.FS

//go:generate bun run build

func main() {
	var _ = interpreter.Natives

	e := echo.New()
	e.Renderer = createRenderer()
	e.HTTPErrorHandler = customHTTPErrorHandler
	e.Logger.SetLevel(log.DEBUG)
	e.Logger.SetHeader("${time_rfc3339} ${level} ${short_file}:${line} ${message}")

	templates := map[string]string{
		"/":                 "index.html",
		"/docs":             "docs.html",
		"/docs/api":         "docs_api.html",
		"/docs/api/std":     "docs_api_std.html",
		"/docs/api/builtin": "docs_api_builtin.html",
	}

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: func(c echo.Context) bool {
			return strings.HasPrefix(c.Request().URL.Path, "/static")
		},
		Format: "${time_rfc3339} ${method} ${uri} ${status} ${latency_human}\n",
	}))

	context := &TemplateContext{
		Path:            "",
		StandardLibrary: lib.StdlibDocs,
		Builtins:        append(interpreter.GlobalDocs, std_language.Builtins...),
	}

	for path, tmpl := range templates {
		path := path
		tmpl := tmpl
		e.GET(path, func(c echo.Context) error {
			return c.Render(http.StatusOK, tmpl, context.CloneFor(c))
		})
	}

	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Filesystem: http.FS(static),
	}))

	e.Logger.Fatal(e.Start(":3000"))
}
