package email

import "simple-auth/pkg/lib/multitemplate"

var templateDefinitions = map[string][]string{
	"welcome":        {"templates/email/welcome.tmpl"},
	"forgotPassword": {"templates/email/forgotPassword.tmpl"},
	"verification":   {"templates/email/verification.tmpl"},
}
var templateEngine multitemplate.TemplateRenderer

func init() {
	templateEngine = multitemplate.New().LoadTemplates(templateDefinitions)
}
