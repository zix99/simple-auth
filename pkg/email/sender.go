package email

import (
	"bytes"
	"errors"
	"fmt"
	"net/smtp"
	"simple-auth/pkg/config"
	"strings"
)

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

func (s *EmailService) sendEmail(cfg *config.ConfigEmail, to string, templateName string, data IEmailData) error {
	if !cfg.Enabled {
		s.logger.Infof("Skipping sending email %s to %s, disabled", templateName, to)
		return errors.New("Email disabled")
	}

	s.logger.Infof("Sending %s email to %s...", templateName, to)
	auth := smtp.PlainAuth(cfg.Identity, cfg.Username, cfg.Password, extractHostname(cfg.Host))

	templateData := &emailData{
		From:  fmt.Sprintf("%s <%s>", data.Data().Company, cfg.From),
		To:    to,
		Model: data,
	}

	var buf bytes.Buffer
	err := templateEngine.Render(&buf, templateName, templateData)
	if err != nil {
		s.logger.Warn(err)
		return err
	}

	err = smtp.SendMail(cfg.Host, auth, cfg.From, []string{to}, buf.Bytes())
	if err != nil {
		s.logger.Warn(err)
	} else {
		s.logger.Infof("Email sent %d bytes to %s", buf.Len(), to)
	}
	return err
}
