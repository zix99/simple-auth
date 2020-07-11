package email

import (
	"bytes"
	"fmt"
	"net/smtp"
	"simple-auth/pkg/config"
	"simple-auth/pkg/lib/multitemplate"
	"strings"

	"github.com/sirupsen/logrus"
)

var templateDefinitions = map[string][]string{
	"welcome": {"templates/email/welcome.tmpl"},
}
var templateEngine multitemplate.TemplateRenderer

func init() {
	templateEngine = multitemplate.New().LoadTemplates(templateDefinitions)
}

type WelcomeEmailData struct {
	Company   string
	Name      string
	WebHost   string
	AccountID string
}

type emailData struct {
	From  string
	To    string
	Model interface{}
}

func extractHostname(host string) string {
	idx := strings.Index(host, ":")
	if idx < 0 {
		return host
	}
	return host[:idx]
}

func SendWelcomeEmail(cfg *config.ConfigEmail, to string, data *WelcomeEmailData) error {
	if !cfg.Enabled {
		return nil
	}

	logrus.Infof("Sending welcome email to %s...", to)
	auth := smtp.PlainAuth(cfg.Identity, cfg.Username, cfg.Password, extractHostname(cfg.Host))

	templateData := &emailData{
		From:  fmt.Sprintf("%s <%s>", data.Company, cfg.From),
		To:    to,
		Model: data,
	}

	var buf bytes.Buffer
	err := templateEngine.Render(&buf, "welcome", templateData)
	if err != nil {
		logrus.Warn(err)
		return err
	}

	err = smtp.SendMail(cfg.Host, auth, cfg.From, []string{to}, buf.Bytes())
	if err != nil {
		logrus.Warn(err)
	} else {
		logrus.Infof("Email sent %d bytes to %s", buf.Len(), to)
	}
	return err
}
