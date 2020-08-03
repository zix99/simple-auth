package main

import (
	"encoding/json"
	"html/template"
	"io"
	"simple-auth/pkg/lib/multitemplate"

	"github.com/labstack/echo/v4"
)

var helpers = template.FuncMap{
	"json": func(obj interface{}) template.JS {
		jsonBytes, _ := json.Marshal(obj)
		return template.JS(jsonBytes)
	},
}

var templateDefinitions = map[string][]string{
	"home": {"templates/web/home.tmpl", "templates/web/layout.tmpl"},
}

type templateRenderer struct {
	templates *multitemplate.TemplateSet
}

func (t *templateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.Render(w, name, data)
}

func newTemplateRenderer(autoreload bool) *templateRenderer {
	engine := multitemplate.New().
		Helpers(helpers).
		AutoReload(autoreload).
		LoadTemplates(templateDefinitions)

	return &templateRenderer{
		templates: engine,
	}
}
