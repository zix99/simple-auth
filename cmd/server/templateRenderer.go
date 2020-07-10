package main

import (
	"encoding/json"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"simple-auth/pkg/config"
	"time"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

var helpers = template.FuncMap{
	"json": func(obj interface{}) template.HTML {
		jsonBytes, _ := json.Marshal(obj)
		return template.HTML(jsonBytes)
	},
}

var templateDefinitions = map[string][]string{
	"createAccount": {"templates/web/createAccount.tmpl", "templates/web/layout.tmpl", "templates/web/layoutVue.tmpl"},
	"home":          {"templates/web/home.tmpl", "templates/web/layout.tmpl"},
}

type templateSet struct {
	templates  map[string]*template.Template
	lastUpdate map[string]time.Time
}

func getFileLastUpdate(files ...string) time.Time {
	var maxTime time.Time
	for _, fname := range files {
		if info, err := os.Stat(fname); err == nil {
			if info.ModTime().After(maxTime) {
				maxTime = info.ModTime()
			}
		}
	}
	return maxTime
}

func compileTemplate(files ...string) (*template.Template, error) {
	basename := filepath.Base(files[0])
	return template.New(basename).Funcs(helpers).ParseFiles(files...)
}

func (t *templateSet) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if !config.Global.Production {
		defn := templateDefinitions[name]
		lastUpdate := t.lastUpdate[name]
		updated := getFileLastUpdate(defn...)

		if updated.After(lastUpdate) {
			logrus.Infof("Reloading template %s...", name)
			t.lastUpdate[name] = updated
			runtimeTemplate, err := compileTemplate(defn...)
			if err != nil {
				logrus.Error(err)
			} else {
				t.templates[name] = runtimeTemplate
			}
		}
	}
	return t.templates[name].Execute(w, data)
}

func newTemplateSet() *templateSet {
	ret := &templateSet{
		templates:  map[string]*template.Template{},
		lastUpdate: map[string]time.Time{},
	}
	for k, v := range templateDefinitions {
		logrus.Infof("Loading template %s...", k)
		ret.templates[k] = template.Must(compileTemplate(v...))
		ret.lastUpdate[k] = time.Now()
	}
	return ret
}
