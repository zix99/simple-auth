package email

import "github.com/sirupsen/logrus"

type EmailService struct {
	logger logrus.FieldLogger
}

func New(logger logrus.FieldLogger) *EmailService {
	return &EmailService{
		logger: logger,
	}
}
