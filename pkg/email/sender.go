package email

import (
	"bytes"
	"fmt"
	"html/template"
	"simple-auth/pkg/appcontext"
	"simple-auth/pkg/instrumentation"
)

var emailCounter instrumentation.Counter = instrumentation.NewCounter("sa_email_sends", "Email sending metrics", "template", "success")

type emailData struct {
	From  template.HTML
	To    template.HTML
	Model interface{}
}

func (s *EmailService) sendEmail(to string, templateName string, data IEmailData) error {
	log := appcontext.GetLogger(s.ctx)

	log.Infof("Sending %s email to %s...", templateName, to)

	templateData := &emailData{
		From:  template.HTML(fmt.Sprintf("%s <%s>", data.Data().Company, s.from)),
		To:    template.HTML(to),
		Model: data,
	}

	var buf bytes.Buffer
	err := templateEngine.Render(&buf, templateName, templateData)
	if err != nil {
		log.Warn(err)
		return err
	}

	err = s.engine.Send(to, s.from, buf.Bytes())
	if err != nil {
		log.Warn(err)
	} else {
		log.Infof("Email sent %d bytes to %s", buf.Len(), to)
	}

	emailCounter.Inc(templateName, err == nil)

	return err
}
