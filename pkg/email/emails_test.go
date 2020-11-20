package email

import (
	"simple-auth/pkg/email/engine"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestSendWelcomeEmail(t *testing.T) {
	mock := engine.NewMockEngine(nil)
	service := New(logrus.StandardLogger(), mock, "test@test.com")
	service.SendWelcomeEmail("to@to.com", &WelcomeEmailData{
		Name: "Bob",
		EmailData: EmailData{
			Company: "SimpleAuth",
			BaseURL: "http://example.com",
		},
	})

	assert.Equal(t, 1, mock.SendCount())
	assert.Contains(t, mock.LastEmail(), "Bob")
	assert.Contains(t, mock.LastEmail(), "SimpleAuth")
	assert.Contains(t, mock.LastEmail(), "example.com")
	assert.Contains(t, mock.LastEmail(), "<test@test.com>")
	assert.Contains(t, mock.LastEmail(), "to@to.com")
}

func TestVerificationEmail(t *testing.T) {
	mock := engine.NewMockEngine(nil)
	service := New(logrus.StandardLogger(), mock, "test@test.com")
	service.SendVerificationEmail("to@to.com", &VerificationData{
		ActivationLink: "http://bla.com/activate",
		EmailData: EmailData{
			Company: "SimpleAuth",
			BaseURL: "http://example.com",
		},
	})

	assert.Equal(t, 1, mock.SendCount())
	assert.Contains(t, mock.LastEmail(), "http://bla.com/activate")
	assert.Contains(t, mock.LastEmail(), "SimpleAuth")
	assert.Contains(t, mock.LastEmail(), "example.com")
}

func TestChangePasswordEmail(t *testing.T) {
	mock := engine.NewMockEngine(nil)
	service := New(logrus.StandardLogger(), mock, "test@test.com")
	service.SendForgotPasswordEmail("to@to.com", &ForgotPasswordData{
		ResetLink:     "http://bla.com/reset",
		ResetDuration: "10 minutes",
		EmailData: EmailData{
			Company: "SimpleAuth",
			BaseURL: "http://example.com",
		},
	})

	assert.Equal(t, 1, mock.SendCount())
	assert.Contains(t, mock.LastEmail(), "http://bla.com/reset")
	assert.Contains(t, mock.LastEmail(), "10 minutes")
	assert.Contains(t, mock.LastEmail(), "SimpleAuth")
	assert.Contains(t, mock.LastEmail(), "example.com")
}
