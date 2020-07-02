package main

import (
	"html/template"
	"io"
	"os"
	"simple-auth/pkg/config"
	"time"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

var templateDefinitions = map[string][]string{
	"createAccount": {"templates/createAccount.tmpl", "templates/layout.tmpl", "templates/layoutVue.tmpl"},
	"home":          {"templates/home.tmpl", "templates/layout.tmpl"},
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

func (t *templateSet) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if !config.Global.Production {
		defn := templateDefinitions[name]
		lastUpdate := t.lastUpdate[name]
		updated := getFileLastUpdate(defn...)

		if updated.After(lastUpdate) {
			logrus.Infof("Reloading template %s...", name)
			t.lastUpdate[name] = updated
			runtimeTemplate, err := template.ParseFiles(defn...)
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
		ret.templates[k] = template.Must(template.ParseFiles(v...))
		ret.lastUpdate[k] = time.Now()
	}
	return ret
}
