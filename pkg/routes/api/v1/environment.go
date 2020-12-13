package v1

import (
	"simple-auth/pkg/config"
	"simple-auth/pkg/db"
	"simple-auth/pkg/email"
	"simple-auth/pkg/services"

	"github.com/sirupsen/logrus"
)

type Environment struct {
	accountService    services.AccountService
	localLoginService services.LocalLoginService
	twoFactorService  services.TwoFactorService
	loginConfig       *config.ConfigLoginCookie
}

func NewEnvironment(config *config.Config, db db.SADB) *Environment {
	emailService := email.NewFromConfig(logrus.StandardLogger(), &config.Email)
	return &Environment{
		services.NewAccountService(&config.Metadata, &config.Web, emailService),
		services.NewLocalLoginService(emailService, &config.Metadata, &config.Web.Login.TwoFactor, &config.Web.Requirements, config.Web.GetBaseURL()),
		services.NewTwoFactorService(&config.Web.Login.TwoFactor),
		&config.Web.Login.Cookie,
	}
}
