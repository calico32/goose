package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"os"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type TemplateRenderer struct {
	templates *template.Template
}

func (r *TemplateRenderer) ReadFile(path string) ([]byte, error) {
	var fs fs.FS
	if mode == "release" {
		fs = tmplFS
	} else {
		fs = os.DirFS(".")
	}

	file, err := fs.Open("views/" + path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)
}

func (r *TemplateRenderer) WriteCachable(content []byte, c echo.Context) error {
	hash := sha256.Sum256(content)
	etag := `"` + base64.StdEncoding.EncodeToString(hash[:]) + `"`
	c.Response().Header().Set("ETag", etag)
	if c.Request().Header.Get("If-None-Match") == etag {
		c.Response().WriteHeader(http.StatusNotModified)
		return nil
	} else {
		_, err := c.Response().Write(content)
		if err != nil {
			c.Logger().Error(err)
			return err
		}
	}
	return nil
}

func createRenderer() *TemplateRenderer {
	funcMap := template.FuncMap{
		"codeBlock": codeBlock,
		"getCurrentYear": func() int {
			return time.Now().Year()
		},
	}

	templates := template.New("").Funcs(funcMap)
	return &TemplateRenderer{
		templates: templates,
	}
}

func (r *TemplateRenderer) Render(w io.Writer, tmplName string, ctx any, c echo.Context) error {
	tmpl := template.Must(r.templates.Clone())
	tmpl, _, err := loadTemplateChainFromFile(tmpl, tmplName, make([]string, 0, 10))
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusNotFound, "Page not found")
	}
	return tmpl.ExecuteTemplate(w, tmplName, ctx)
}

func ExecuteToString(t *template.Template, templateName string, ctx any) (string, error) {
	var b strings.Builder
	err := t.ExecuteTemplate(&b, templateName, ctx)
	return b.String(), err
}

func DoRender[Ctx Context](r *TemplateRenderer, path string, templateName string, ctx Ctx, c echo.Context) error {
	isHTMX := strings.Contains(c.Request().Header.Get("HX-Request"), "true")
	if isHTMX {
		tmpl := template.Must(r.templates.Clone())
		tmpl, _, err := loadTemplateChainFromFile(tmpl, templateName, make([]string, 0, 10))
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusNotFound, "Page not found")
		}

		title, err := ExecuteToString(tmpl, "docs-title", ctx.CloneFor(c))
		if err != nil {
			c.Logger().Error(err)
			return err
		}

		page, err := ExecuteToString(tmpl, "docs-content", ctx.CloneFor(c))
		if err != nil {
			c.Logger().Error(err)
			return err
		}

		c.Response().Header().Set("Content-Type", "text/html")
		htmxContent := fmt.Sprintf(`
			<span id="title" class="hidden" hx-swap-oob="true">%s | goose</span>
			%s
		`, title, page)
		err = r.WriteCachable([]byte(htmxContent), c)
		if err != nil {
			c.Logger().Error(err)
			return err
		}

		return nil
	}

	isDocs := strings.HasPrefix(path, "/docs")
	if isDocs {
		// not HTMX inside of /docs, which means we need to render the full page
		f, err := r.ReadFile(templateName)
		if err != nil {
			c.Logger().Error(err)
			return err
		}

		tmpl := fmt.Sprintf(`{{template "+docs_layout.html" .}} %s`, f)
		c.Response().Header().Set("Content-Type", "text/html")

		buf := new(bytes.Buffer)
		err = r.RenderFromString(buf, tmpl, ctx.CloneFor(c), c)
		if err != nil {
			c.Logger().Error(err)
			return err
		}
		err = r.WriteCachable(buf.Bytes(), c)
		if err != nil {
			c.Logger().Error(err)
			return err
		}

		return nil
	}

	err := c.Render(http.StatusOK, templateName, ctx.CloneFor(c))
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	return nil
}

func (t *TemplateRenderer) RenderToString(name string, data any, c echo.Context) (string, error) {
	var b strings.Builder
	err := t.Render(&b, name, data, c)
	return b.String(), err
}

func (t *TemplateRenderer) RenderFromString(w io.Writer, content string, data any, c echo.Context) error {
	tmpl := template.Must(t.templates.Clone())
	tmpl, _, err := loadTemplateChain(tmpl, "<content>", content, make([]string, 0, 10))
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusNotFound, "Page not found")
	}
	// fmt.Printf("Loaded templates: %v\n", chain)
	return tmpl.Execute(w, data)
}

func loadTemplateChainFromFile(tmpl *template.Template, name string, chain []string) (*template.Template, []string, error) {
	var fs fs.FS
	if mode == "release" {
		fs = tmplFS
	} else {
		fs = os.DirFS(".")
	}

	file, err := fs.Open("views/" + name)
	if err != nil {
		return nil, chain, err
	}
	defer file.Close()

	f, err := io.ReadAll(file)
	if err != nil {
		return nil, chain, err
	}

	return loadTemplateChain(tmpl, name, string(f), chain)
}

func loadTemplateChain(tmpl *template.Template, name string, content string, chain []string) (*template.Template, []string, error) {
	// check if the file has template directives
	pattern := regexp.MustCompile(`{{template "([^\"]+)"`)
	matches := pattern.FindAllStringSubmatch(content, -1)
	if len(matches) == 0 {
		tmpl, err := tmpl.New(name).Parse(content)
		chain = append(chain, name)
		if err != nil {
			return nil, chain, err
		}
		return tmpl, chain, nil
	}

	// load chains
	for _, m := range matches {
		for _, match := range m[1:] {
			if !strings.HasSuffix(match, ".html") {
				// not a template file
				continue
			}
			if slices.Contains(chain, match) {
				return nil, chain, echo.NewHTTPError(http.StatusInternalServerError, "Circular template dependency")
			}
			var err error
			tmpl, chain, err = loadTemplateChainFromFile(tmpl, match, chain)
			if err != nil {
				return nil, chain, err
			}
		}
	}

	// load the current template
	tmpl, err := tmpl.New(name).Parse(content)
	chain = append(chain, name)
	if err != nil {
		return nil, chain, err
	}

	return tmpl, chain, err
}

func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	message := "Internal Server Error"
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		switch msg := he.Message.(type) {
		case string:
			message = msg
		case nil:
			message = http.StatusText(code)
		}
	}
	if code == 404 {
		c.Render(code, "+error_404.html", nil)
	} else {
		c.Render(code, "+error.html", map[string]any{
			"Code":    code,
			"Message": message,
		})
	}
}
