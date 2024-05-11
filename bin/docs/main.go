package main

import (
	"fmt"
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
	etag "github.com/pablor21/echo-etag/v4"
)

var mode = "release"

//go:generate bun run build

func main() {
	var _ = interpreter.Natives

	e := echo.New()
	renderer := createRenderer()
	e.Renderer = renderer
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

	e.Use(middleware.BodyLimit("256K"))
	e.Use(middleware.Gzip())
	e.Use(middleware.NonWWWRedirect())
	e.Use(middleware.Recover())
	e.Use(etag.Etag())

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: func(c echo.Context) bool {
			return strings.HasPrefix(c.Request().URL.Path, "/static")
		},
		Format: "${time_rfc3339} ${method} ${uri} ${status} ${latency_human}\n",
	}))

	context := &TemplateContext{
		path:            "",
		standardLibrary: lib.StdlibDocs,
		builtins:        append(interpreter.GlobalDocs, std_language.Builtins...),
	}

	prerenderQueue := make(chan *string, 200)

	for _, doc := range context.builtins {
		doc := doc
		path := fmt.Sprintf("/docs/api/builtin/%s", doc.Name)
		prerenderQueue <- &path
		e.GET(path, func(c echo.Context) error {
			builtinCtx := &BuiltinContext{
				TemplateContext: context.CloneFor(c).(*TemplateContext),
				Doc:             doc,
			}
			return DoRender(renderer, path, "+docs_api_builtin_x.html", builtinCtx, c)
		})
	}

	for _, doc := range context.standardLibrary {
		pkg := doc
		path := fmt.Sprintf("/docs/api/std/%s", doc.Name)
		prerenderQueue <- &path
		e.GET(path, func(c echo.Context) error {
			stdlibCtx := &StdlibContext{
				TemplateContext: context.CloneFor(c).(*TemplateContext),
				Pkg:             pkg,
			}
			return DoRender(renderer, path, "+docs_api_std_x.html", stdlibCtx, c)
		})
	}

	for path, tmpl := range templates {
		path := path
		tmpl := tmpl
		prerenderQueue <- &path
		e.GET(path, func(c echo.Context) error {
			return DoRender(renderer, path, tmpl, context.CloneFor(c), c)
		})
	}

	if mode == "release" {
		e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
			Filesystem: http.FS(staticFS),
		}))
	} else {
		e.GET("/static/*", func(c echo.Context) error {
			if !strings.HasPrefix(c.Request().URL.Path, "/static/css") {
				c.Response().Header().Set("Cache-Control", "public, max-age=31536000")
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

	go func() {
		e.Logger.Debugf("Pre-rendering code blocks...")
		for i := 0; i < 4; i++ {
			go prerenderWorker(e, port, prerenderQueue)
		}
	}()

	for _, route := range e.Routes() {
		e.Logger.Debugf("Registered route: %s %s", route.Method, route.Path)
	}

	e.Logger.Fatal(e.Start(":" + port))
}

func prerenderWorker(e *echo.Echo, port string, queue chan *string) {
	for {
		path := <-queue
		if path == nil {
			break
		}
		resp, err := http.Get("http://localhost:" + port + *path)
		if err != nil {
			e.Logger.Errorf("Failed to get response for %s: %v", *path, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			e.Logger.Errorf("Failed to get 200 OK for %s: %v", *path, resp.Status)
			continue
		}
	}
}
