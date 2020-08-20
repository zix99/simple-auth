package email

import (
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
	return s.sendEmail(cfg, to, "welcome", data)
}

type ForgotPasswordData struct {
	EmailData
	ResetLink     string
	ResetDuration string
}

func (s *EmailService) SendForgotPasswordEmail(cfg *config.ConfigEmail, to string, data *ForgotPasswordData) error {
	return s.sendEmail(cfg, to, "forgotPassword", data)
}
