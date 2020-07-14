package main

import (
	"encoding/json"
	"html/template"
	"io"
	"simple-auth/pkg/lib/multitemplate"

	"github.com/labstack/echo"
)

var helpers = template.FuncMap{
	"json": func(obj interface{}) template.HTML {
		jsonBytes, _ := json.Marshal(obj)
		return template.HTML(jsonBytes)
	},
}

var templateDefinitions = map[string][]string{
	"createAccount": {"templates/web/createAccount.tmpl", "templates/web/layout.tmpl", "templates/web/layoutVue.tmpl"},
	"manageAccount": {"templates/web/manageAccount.tmpl", "templates/web/layout.tmpl", "templates/web/layoutVue.tmpl"},
	"home":          {"templates/web/home.tmpl", "templates/web/layout.tmpl"},
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
