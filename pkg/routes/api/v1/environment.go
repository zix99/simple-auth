package v1

import (
	"simple-auth/pkg/config"
	"simple-auth/pkg/email"
	"simple-auth/pkg/services"
)

type Environment struct {
	accountService    services.AccountService
	localLoginService services.LocalLoginService
	twoFactorService  services.TwoFactorService
	oidcService       services.OIDCService
	sessionService    services.SessionService
	loginConfig       *config.ConfigLoginCookie
}

func NewEnvironment(config *config.Config) *Environment {
	emailService := email.NewFromConfig(&config.Email)
	return &Environment{
		services.NewAccountService(&config.Metadata, &config.Web, emailService),
		services.NewLocalLoginService(emailService, &config.Metadata, &config.Web.Login.TwoFactor, &config.Web.Requirements, config.Web.GetBaseURL()),
		services.NewTwoFactorService(&config.Web.Login.TwoFactor),
		services.NewOIDCService(config.Web.Login.OIDC),
		services.NewSessionService(emailService, &config.Web.Login.Cookie, &config.Web.Login.OneTime, &config.Web, &config.Metadata),
		&config.Web.Login.Cookie,
	}
}
