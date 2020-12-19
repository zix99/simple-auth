package services

import (
	"fmt"
	"html/template"
	"simple-auth/pkg/appcontext"
	"simple-auth/pkg/config"
	"simple-auth/pkg/db"
	"simple-auth/pkg/email"
	"simple-auth/pkg/routes/middleware/selector/auth"
	"time"

	"github.com/labstack/echo/v4"
)

type SessionService interface {
	IssueSession(c echo.Context, account db.AccountProvider, source auth.SessionSource) error
	ClearSession(c echo.Context)

	IssueOneTimeToken(acct db.AccountProvider) error
	IssueOneTimeSession(c echo.Context, token string) error

	WithContext(c appcontext.Context) SessionService
}

type sessionService struct {
	emailService  *email.EmailService
	cookieConfig  *config.ConfigLoginCookie
	onetimeConfig *config.OneTimeConfig
	webConfig     *config.ConfigWeb
	metaConfig    *config.ConfigMetadata

	// Contextual vars
	dbOneTime db.AccountAuthOneTime
	context   appcontext.Context
}

var _ SessionService = &sessionService{}

func NewSessionService(emailService *email.EmailService,
	cookieConfig *config.ConfigLoginCookie,
	onetimeConfig *config.OneTimeConfig,
	webConfig *config.ConfigWeb,
	metaConfig *config.ConfigMetadata) SessionService {
	return &sessionService{
		emailService,
		cookieConfig,
		onetimeConfig,
		webConfig,
		metaConfig,
		nil,
		nil,
	}
}

func (s *sessionService) WithContext(c appcontext.Context) SessionService {
	copy := *s
	copy.dbOneTime = appcontext.GetSADB(c)
	copy.context = c
	return &copy
}

func (s *sessionService) IssueSession(c echo.Context, account db.AccountProvider, source auth.SessionSource) error {
	return auth.CreateSession(c, s.cookieConfig, account.Account(), source)
}

func (s *sessionService) ClearSession(c echo.Context) {
	auth.ClearSession(c, s.cookieConfig)
}

func (s *sessionService) IssueOneTimeToken(acct db.AccountProvider) error {
	account := acct.Account()

	duration, err := time.ParseDuration(s.onetimeConfig.TokenDuration)
	if err != nil {
		return fmt.Errorf("invalid onetime duration: %w", err)
	}

	token, err := s.dbOneTime.CreateAccountOneTimeToken(account, duration)
	if err != nil {
		return err
	}

	baseURL := s.webConfig.GetBaseURL()
	go s.emailService.WithContext(s.context).SendForgotPasswordEmail(account.Email, &email.ForgotPasswordData{
		EmailData: email.EmailData{
			Company: s.metaConfig.Company,
			BaseURL: baseURL,
		},
		ResetDuration: s.onetimeConfig.TokenDuration,
		ResetLink:     template.HTML(baseURL + "/onetime?token=" + token),
	})

	return nil
}

func (s *sessionService) IssueOneTimeSession(c echo.Context, token string) error {
	account, err := s.dbOneTime.AssertOneTimeToken(token)
	if err != nil {
		return err
	}

	if err := auth.CreateSession(c, s.cookieConfig, account, auth.SourceOneTime); err != nil {
		return err
	}

	return nil
}
