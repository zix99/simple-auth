package email

import (
	"errors"
	"simple-auth/pkg/appcontext"
	"simple-auth/pkg/config"
	"simple-auth/pkg/email/engine"
	"strings"
)

type EmailService struct {
	engine engine.EmailEngine
	from   string
	ctx    appcontext.Context
}

func New(engine engine.EmailEngine, from string) *EmailService {
	return &EmailService{
		engine,
		from,
		nil,
	}
}

func (s *EmailService) WithContext(c appcontext.Context) *EmailService {
	copy := *s
	copy.ctx = c
	return &copy
}

func NewFromConfig(config *config.ConfigEmail) *EmailService {
	engine := engineFromConfig(config)
	return New(engine, config.From)
}

func engineFromConfig(config *config.ConfigEmail) engine.EmailEngine {
	switch strings.ToLower(config.Engine) {
	case "smtp":
		smtp := config.SMTP
		return engine.NewSMTPEngine(smtp.Host, smtp.Identity, smtp.Username, smtp.Password)
	case "noop":
		return engine.NewNoopEngine(nil)
	case "stdout":
		return engine.NewStdoutEngine()
	}
	return engine.NewNoopEngine(errors.New("engine not specified"))
}
