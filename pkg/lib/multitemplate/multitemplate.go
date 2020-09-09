package multitemplate

import (
	"html/template"
	"io"
	"io/ioutil"
	"path/filepath"
	"simple-auth/pkg/box"
	"time"

	"github.com/sirupsen/logrus"
)

/*
Template loader and auto-reloader when in development mode
*/

type TemplateRenderer interface {
	Render(w io.Writer, name string, data interface{}) error
}

type TemplateSet struct {
	autoreload  bool
	definitions map[string][]string
	templates   map[string]*template.Template
	lastUpdate  map[string]time.Time
	helpers     template.FuncMap
}

func getFileLastUpdate(files ...string) time.Time {
	var maxTime time.Time
	for _, fname := range files {
		if info, err := box.Stat(fname); err == nil {
			if info.ModTime().After(maxTime) {
				maxTime = info.ModTime()
			}
		}
	}
	return maxTime
}

func (s *TemplateSet) compileTemplate(files ...string) (*template.Template, error) {
	basename := filepath.Base(files[0])
	builder := template.New(basename).Funcs(s.helpers)
	for _, fn := range files {
		r, err := box.Read(fn)
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}
		builder.Parse(string(b))
	}
	return builder, nil
}

func (s *TemplateSet) Render(w io.Writer, name string, data interface{}) error {
	if s.autoreload {
		defn := s.definitions[name]
		lastUpdate := s.lastUpdate[name]
		updated := getFileLastUpdate(defn...)

		if updated.After(lastUpdate) {
			logrus.Infof("Reloading template %s...", name)
			s.lastUpdate[name] = updated
			runtimeTemplate, err := s.compileTemplate(defn...)
			if err != nil {
				logrus.Error(err)
			} else {
				s.templates[name] = runtimeTemplate
			}
		}
	}
	return s.templates[name].Execute(w, data)
}

func (s *TemplateSet) Helpers(helpers template.FuncMap) *TemplateSet {
	s.helpers = helpers
	return s
}

func (s *TemplateSet) AutoReload(enabled bool) *TemplateSet {
	s.autoreload = enabled
	return s
}

func (s *TemplateSet) LoadTemplate(name string, templates ...string) *TemplateSet {
	logrus.Infof("Loading template %s...", name)
	s.definitions[name] = templates
	s.templates[name] = template.Must(s.compileTemplate(templates...))
	s.lastUpdate[name] = time.Now()
	return s
}

func (s *TemplateSet) LoadTemplates(definition map[string][]string) *TemplateSet {
	for k, v := range definition {
		s.LoadTemplate(k, v...)
	}
	return s
}

// New creates a template set that has auto-reload capabilities
func New() *TemplateSet {
	return &TemplateSet{
		definitions: map[string][]string{},
		templates:   map[string]*template.Template{},
		lastUpdate:  map[string]time.Time{},
	}
}
