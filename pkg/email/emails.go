package email

import (
	"errors"
	"simple-auth/pkg/config"
)

type IEmailData interface {
	Data() *EmailData
}

func (s EmailData) Data() *EmailData {
	return &s
}

type EmailData struct {
	IEmailData
	Company string
	BaseURL string
}

type WelcomeEmailData struct {
	EmailData
	Name string
}

func (s *EmailService) SendWelcomeEmail(cfg *config.ConfigEmail, to string, data *WelcomeEmailData) error {
	if !cfg.Enabled {
		s.logger.Infof("Skipping sending welcome to %s, disabled", to)
		return errors.New("email disabled")
	}
	return s.sendEmail(&cfg.SMTP, to, "welcome", data)
}

type ForgotPasswordData struct {
	EmailData
	ResetLink     string
	ResetDuration string
}

func (s *EmailService) SendForgotPasswordEmail(cfg *config.ConfigEmail, to string, data *ForgotPasswordData) error {
	if !cfg.Enabled {
		s.logger.Infof("Skipping sending email forgot-password to %s, disabled", to)
		return errors.New("email disabled")
	}
	return s.sendEmail(&cfg.SMTP, to, "forgotPassword", data)
}

type VerificationData struct {
	EmailData
	ActivationLink string
}

func (s *EmailService) SendVerificationEmail(cfg *config.ConfigEmail, to string, data *VerificationData) error {
	if !cfg.Enabled {
		s.logger.Infof("Skipping verification email to %s, disabled", to)
		return errors.New("email disabled")
	}
	return s.sendEmail(&cfg.SMTP, to, "verification", data)
}
