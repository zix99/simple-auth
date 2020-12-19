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
	FindAccountByEmail(email string) (*db.Account, error)
}

type accountService struct {
	emailService *email.EmailService
	metaConfig   *config.ConfigMetadata
	baseURL      string
	dbAccount    db.AccountStore
}

var _ AccountService = &accountService{}

func NewAccountService(configMeta *config.ConfigMetadata, configWeb *config.ConfigWeb, emailService *email.EmailService) AccountService {
	return &accountService{
		emailService,
		configMeta,
		configWeb.GetBaseURL(),
		nil,
	}
}

func (s *accountService) WithContext(ctx appcontext.Context) AccountService {
	copy := *s
	copy.dbAccount = appcontext.GetSADB(ctx)
	return &copy
}

func (s *accountService) CreateAccount(name, emailAddress string) (*db.Account, error) {
	account, err := s.dbAccount.CreateAccount(name, emailAddress)
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

func (s *accountService) FindAccountByEmail(email string) (*db.Account, error) {
	account, err := s.dbAccount.FindAccountByEmail(email)
	return account, err
}
