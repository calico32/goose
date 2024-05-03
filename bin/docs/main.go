package main

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/calico32/goose/interpreter"
	"github.com/calico32/goose/lib"
	std_language "github.com/calico32/goose/lib/std/language"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

var mode = "release"

//go:generate bun run build

func main() {
	var _ = interpreter.Natives

	e := echo.New()
	e.Renderer = createRenderer()
	e.HTTPErrorHandler = customHTTPErrorHandler

	var refresh *PageRefresh
	if mode != "release" {
		refresh = NewPageRefresh(e.Logger)
	}

	if mode == "release" {
		e.Logger.SetLevel(log.INFO)
		e.Logger.SetHeader("${time_rfc3339} ${level} ${message}")
	} else {
		e.Logger.SetLevel(log.DEBUG)
		e.Logger.SetHeader("${time_rfc3339} ${level} ${short_file}:${line} ${message}")
	}

	e.Use(middleware.Gzip())
	e.Use(middleware.NonWWWRedirect())

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
			if mode == "release" {
				c.Response().Header().Set("Cache-Control", "public, max-age=3600")
			}
			return c.Render(http.StatusOK, tmpl, context.CloneFor(c))
		})
	}

	if mode == "release" {
		e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
			Filesystem: http.FS(staticFS),
		}))
	} else {
		e.GET("/static/*", func(c echo.Context) error {
			if mode == "release" {
				c.Response().Header().Set("Cache-Control", "public, max-age=3600")
			}
			return c.File("." + c.Request().URL.Path)
		})
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	go func() {
		time.Sleep(100 * time.Millisecond)
		if mode != "release" {
			refresh.Refresh()
		}
	}()

	e.Logger.Fatal(e.Start(":" + port))
}
