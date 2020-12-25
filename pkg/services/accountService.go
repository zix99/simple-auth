package services

import (
	"errors"
	"regexp"
	"simple-auth/pkg/appcontext"
	"simple-auth/pkg/config"
	"simple-auth/pkg/db"
	"simple-auth/pkg/email"
	"simple-auth/pkg/instrumentation"
	"simple-auth/pkg/saerrors"
	"unicode/utf8"
)

var accountCreateCounter instrumentation.Counter = instrumentation.NewCounter("sa_account_create", "Account creation counter")

type AccountService interface {
	WithContext(ctx appcontext.Context) AccountService
	CreateAccount(name, email string) (*db.Account, error)
	FindAccountByEmail(email string) (*db.Account, error)

	HasUnsatisfiedStipulations(account *db.Account) bool
}

type accountService struct {
	emailService   *email.EmailService
	metaConfig     *config.ConfigMetadata
	baseURL        string
	dbAccount      db.AccountStore
	dbStipulations db.AccountStipulations
}

var _ AccountService = &accountService{}

func NewAccountService(configMeta *config.ConfigMetadata, configWeb *config.ConfigWeb, emailService *email.EmailService) AccountService {
	return &accountService{
		emailService,
		configMeta,
		configWeb.GetBaseURL(),
		nil,
		nil,
	}
}

func (s *accountService) WithContext(ctx appcontext.Context) AccountService {
	copy := *s
	sadb := appcontext.GetSADB(ctx)
	copy.dbAccount = sadb
	copy.dbStipulations = sadb
	return &copy
}

const (
	AccountEmailExists saerrors.ErrorCode = "account-email-exists"
)

func (s *accountService) CreateAccount(name, emailAddress string) (*db.Account, error) {
	if account, _ := s.dbAccount.FindAccountByEmail(emailAddress); account != nil {
		return nil, AccountEmailExists.New()
	}

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

	accountCreateCounter.Inc()

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

func (s *accountService) HasUnsatisfiedStipulations(account *db.Account) bool {
	return s.dbStipulations.AccountHasUnsatisfiedStipulations(account)
}
