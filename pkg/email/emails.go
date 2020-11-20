package email

import "html/template"

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

func (s *EmailService) SendWelcomeEmail(to string, data *WelcomeEmailData) error {
	return s.sendEmail(to, "welcome", data)
}

type ForgotPasswordData struct {
	EmailData
	ResetLink     template.HTML
	ResetDuration string
}

func (s *EmailService) SendForgotPasswordEmail(to string, data *ForgotPasswordData) error {
	return s.sendEmail(to, "forgotPassword", data)
}

type VerificationData struct {
	EmailData
	ActivationLink template.HTML
}

func (s *EmailService) SendVerificationEmail(to string, data *VerificationData) error {
	return s.sendEmail(to, "verification", data)
}
