package email

import (
	"os"
	"path/filepath"
	"simple-auth/pkg/lib/multitemplate"
	"strings"
)

var templateDefinitions = map[string][]string{
	"welcome":        {"templates/email/welcome.tmpl"},
	"forgotPassword": {"templates/email/forgotPassword.tmpl"},
	"verification":   {"templates/email/verification.tmpl"},
}
var templateEngine multitemplate.TemplateRenderer

func init() {
	// HACK: When testing, make sure we're in the correct path to build the templates
	if strings.HasSuffix(os.Args[0], ".test") {
		cwd, _ := os.Getwd()
		if strings.HasSuffix(cwd, "/pkg/email") {
			os.Chdir(filepath.Join(cwd, "../../"))
		}
	}
	templateEngine = multitemplate.New().LoadTemplates(templateDefinitions)
}
