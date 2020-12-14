//+build ignore

package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"simple-auth/pkg/box/boxutil"
	"strings"
	"text/template"

	"github.com/sirupsen/logrus"
)

const relativePath = "../../"

type TemplateEmbedFile struct {
	Name       string
	Bytes      []byte
	Compressed bool
	Size       int64
	ModTime    int64
}
type TemplateEmbedBox struct {
	BoxName          string
	BuildConstraints string
	Files            map[string]*TemplateEmbedFile
}

var tmplHelpers = map[string]interface{}{
	"bytes": fmtByteSlice,
}
var tmpl = template.Must(template.New("").Funcs(tmplHelpers).Parse(`//+build {{ $.BuildConstraints }}

package box

// Code generated by go generate; DO NOT EDIT.

import "simple-auth/pkg/box/boxutil"

func init() {
    {{- range $name, $f := .Files }}
			{{ $.BoxName }}.Add("{{ $name }}", &EmbedFile{
					name:    "{{ $f.Name }}",
					size:    {{ $f.Size }},
					modtime: {{ $f.ModTime }},
					bytes: {{if $f.Compressed}}boxutil.Must(boxutil.Decompress([]byte{ {{ bytes $f.Bytes }} })){{else}}[]byte{ {{ bytes $f.Bytes }} }{{end}},
					compressed: {{ $f.Compressed }},
				})
    {{- end }}
}`),
)

func fmtByteSlice(s []byte) string {
	builder := strings.Builder{}

	for _, v := range s {
		builder.WriteString(fmt.Sprintf("%d,", int(v)))
	}

	return builder.String()
}

func collectFiles(boxpath string, compressEnabled bool) (*TemplateEmbedBox, error) {
	box := &TemplateEmbedBox{
		Files: make(map[string]*TemplateEmbedFile),
	}
	realPath := filepath.Join(relativePath, boxpath)

	err := filepath.Walk(realPath, func(wpath string, info os.FileInfo, err error) error {
		if info == nil || info.IsDir() {
			return nil // Skip
		}

		b, err := ioutil.ReadFile(wpath)
		if err != nil {
			return err
		}

		relPath := filepath.Clean(wpath[len(relativePath):])

		if compressEnabled {
			b, err = boxutil.Compress(b)
			if err != nil {
				return err
			}
		}

		logrus.Infof("Adding file '%s' (%d bytes, %d compressed)...", relPath, info.Size(), len(b))

		box.Files[relPath] = &TemplateEmbedFile{
			Name:       filepath.Base(relPath),
			Size:       info.Size(),
			ModTime:    info.ModTime().Unix(),
			Bytes:      b,
			Compressed: compressEnabled,
		}
		return nil
	})

	return box, err
}

func writeBox(boxname string, box *TemplateEmbedBox) error {
	builder := &bytes.Buffer{}
	if err := tmpl.Execute(builder, box); err != nil {
		return err
	}

	data, err := format.Source(builder.Bytes())
	if err != nil {
		return err
	}

	outFilename := boxname + ".gen.go"
	if err := ioutil.WriteFile(outFilename, data, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func main() {
	compress := flag.Bool("compress", false, "Should compress data")
	codename := flag.String("codename", "Global", "Box name to add files to")
	constraints := flag.String("constraints", "box", "Set build constraints for box")
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		logrus.Fatal("Invalid number of arguments")
	}

	boxname := args[0]
	boxpath := filepath.Clean(args[1])

	logrus.Infof("Generating box '%s' at '%s' (Compress:%v)...\n", boxname, boxpath, *compress)
	cwd, _ := os.Getwd()
	logrus.Infof("Current dir: %s\n", cwd)

	box, err := collectFiles(boxpath, *compress)
	if err != nil {
		logrus.Fatal("Error building box: ", err)
	}

	box.BoxName = *codename
	box.BuildConstraints = *constraints

	err = writeBox(boxname, box)
	if err != nil {
		logrus.Fatal(err)
	}
}
