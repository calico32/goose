package main

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"os"
	"regexp"
	"slices"

	"github.com/calico32/goose/lib/types"
	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

type TemplateContext struct {
	Path            string
	StandardLibrary []types.StdlibDoc
	Builtins        []types.BuiltinDoc
}

func (t *TemplateContext) Clone() *TemplateContext {
	return &TemplateContext{
		Path:            t.Path,
		StandardLibrary: t.StandardLibrary,
		Builtins:        t.Builtins,
	}
}

func (t *TemplateContext) CloneFor(c echo.Context) *TemplateContext {
	x := t.Clone()
	x.Path = c.Request().URL.Path
	return x
}

func createRenderer() *Template {
	funcMap := template.FuncMap{
		"codeBlock": codeBlock,
	}

	templates := template.New("").Funcs(funcMap)
	return &Template{
		templates: templates,
	}
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl := template.Must(t.templates.Clone())
	tmpl, chain, err := loadTemplateChain(tmpl, name, make([]string, 0, 10))
	c.Logger().Debugf("Loaded templates: %v", chain)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusNotFound, "Page not found")
	}
	return tmpl.ExecuteTemplate(w, name, data)
}

func loadTemplateChain(tmpl *template.Template, name string, chain []string) (*template.Template, []string, error) {
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

	// check if the file has template directives
	pattern := regexp.MustCompile(`{{template "([^\"]+)"`)
	matches := pattern.FindAllStringSubmatch(string(f), -1)
	if len(matches) == 0 {
		fmt.Printf("Loading %s\n", name)
		tmpl, err := tmpl.ParseFS(fs, "views/"+name)
		chain = append(chain, name)
		if err != nil {
			return nil, chain, err
		}
		return tmpl, chain, nil
	}

	// load chains
	for _, m := range matches {
		for _, match := range m[1:] {
			if slices.Contains(chain, match) {
				return nil, chain, echo.NewHTTPError(http.StatusInternalServerError, "Circular template dependency")
			}
			tmpl, chain, err = loadTemplateChain(tmpl, match, chain)
			if err != nil {
				return nil, chain, err
			}
		}
	}

	// load the current template
	tmpl, err = tmpl.ParseFS(fs, "views/"+name)
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
		c.Logger().Error(err)
		c.Render(code, "+error.html", map[string]any{
			"Code":    code,
			"Message": message,
		})
	}
}
