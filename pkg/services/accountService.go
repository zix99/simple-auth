package services

import (
	"errors"
	"regexp"
	"simple-auth/pkg/appcontext"
	"simple-auth/pkg/config"
	"simple-auth/pkg/db"
	"simple-auth/pkg/email"
	"unicode/utf8"
)

type AccountService interface {
	WithContext(ctx appcontext.Context) AccountService
	CreateAccount(name, email string) (*db.Account, error)
}

type accountService struct {
	emailService *email.EmailService
	metaConfig   *config.ConfigMetadata
	baseURL      string
	context      appcontext.Context
}

var _ AccountService = &accountService{}

func NewAccountService(config *config.Config, emailService *email.EmailService) AccountService {
	return &accountService{
		emailService,
		&config.Metadata,
		config.Web.GetBaseURL(),
		nil,
	}
}

func (s *accountService) WithContext(ctx appcontext.Context) AccountService {
	copy := *s
	copy.context = ctx
	return &copy
}

func (s *accountService) CreateAccount(name, emailAddress string) (*db.Account, error) {
	db := appcontext.GetSADB(s.context)
	account, err := db.CreateAccount(name, emailAddress)
	if err != nil {
		return nil, err
	}

	go s.emailService.SendWelcomeEmail(emailAddress, &email.WelcomeEmailData{
		EmailData: email.EmailData{
			Company: s.metaConfig.Company,
			BaseURL: s.baseURL,
		},
		Name: name,
	})

	return account, nil
}

var emailValidationRegex = regexp.MustCompile(`(?i)^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$`)

func validateEmail(email string) error {
	elen := utf8.RuneCountInString(email)
	if elen < 5 {
		return errors.New("email too short")
	}
	if !emailValidationRegex.MatchString(email) {
		return errors.New("invalid email")
	}
	return nil
}
