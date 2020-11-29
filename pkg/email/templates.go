package email

import (
	"simple-auth/pkg/lib/multitemplate"
	"simple-auth/pkg/testutil"
)

var templateDefinitions = map[string][]string{
	"welcome":        {"templates/email/welcome.tmpl"},
	"forgotPassword": {"templates/email/forgotPassword.tmpl"},
	"verification":   {"templates/email/verification.tmpl"},
}
var templateEngine multitemplate.TemplateRenderer

func init() {
	// HACK: When testing, make sure we're in the correct path to build the templates
	if testutil.IsTesting() {
		testutil.SetRootWorkDir()
	}
	templateEngine = multitemplate.New().LoadTemplates(templateDefinitions)
}
