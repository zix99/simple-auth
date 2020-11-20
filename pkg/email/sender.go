package email

import (
	"bytes"
	"fmt"
	"html/template"
)

type emailData struct {
	From  template.HTML
	To    template.HTML
	Model interface{}
}

func (s *EmailService) sendEmail(to string, templateName string, data IEmailData) error {
	s.logger.Infof("Sending %s email to %s...", templateName, to)

	templateData := &emailData{
		From:  template.HTML(fmt.Sprintf("%s <%s>", data.Data().Company, s.from)),
		To:    template.HTML(to),
		Model: data,
	}

	var buf bytes.Buffer
	err := templateEngine.Render(&buf, templateName, templateData)
	if err != nil {
		s.logger.Warn(err)
		return err
	}

	err = s.engine.Send(to, s.from, buf.Bytes())
	if err != nil {
		s.logger.Warn(err)
	} else {
		s.logger.Infof("Email sent %d bytes to %s", buf.Len(), to)
	}
	return err
}
