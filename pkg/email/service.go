package email

import (
	"errors"
	"simple-auth/pkg/config"
	"simple-auth/pkg/email/engine"

	"github.com/sirupsen/logrus"
)

type EmailService struct {
	logger logrus.FieldLogger
	engine engine.EmailEngine
	from   string
}

func New(logger logrus.FieldLogger, engine engine.EmailEngine, from string) *EmailService {
	return &EmailService{
		logger,
		engine,
		from,
	}
}

func NewFromConfig(logger logrus.FieldLogger, config *config.ConfigEmail) *EmailService {
	if !config.Enabled {
		return New(logger, engine.NewNoopEngine(errors.New("email not enabled")), config.SMTP.From)
	}

	smtp := config.SMTP
	engine := engine.NewSMTPEngine(smtp.Host, smtp.Identity, smtp.Username, smtp.Password)
	return New(logger, engine, smtp.From)
}
