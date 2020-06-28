package main

import (
	"html/template"
	"io"

	"github.com/labstack/echo"
)

type templateSet struct {
	templates map[string]*template.Template
}

func (t *templateSet) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return  t.templates[name].Execute(w, data)
}

func newTemplateSet() *templateSet {
	return &templateSet{
		templates: map[string]*template.Template {
			"createAccount": template.Must(template.ParseFiles("templates/createAccount.tmpl", "templates/layout.tmpl", "templates/layoutVue.tmpl")),
		},
	}
}
