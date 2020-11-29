package ui

import (
	"simple-auth/pkg/config"
	"simple-auth/pkg/db"
	"simple-auth/pkg/email"
	"simple-auth/pkg/routes/common"
	"simple-auth/pkg/routes/middleware"
	"simple-auth/pkg/routes/middleware/recaptcha"
	"simple-auth/pkg/routes/middleware/selector"
	"simple-auth/pkg/routes/middleware/selector/auth"
	"simple-auth/pkg/services"
	"time"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

type environment struct {
	db                db.SADB
	localLoginService services.LocalLoginService

	meta   *config.ConfigMetadata
	config *config.ConfigWeb
	email  *config.ConfigEmail
}

func NewController(db db.SADB, meta *config.ConfigMetadata, config *config.ConfigWeb, emailConfig *config.ConfigEmail) common.Controller {
	emailService := email.NewFromConfig(logrus.StandardLogger(), emailConfig)
	return &environment{
		db:                db,
		localLoginService: services.NewLocalLoginService(db, emailService, meta, &config.Login.TwoFactor, &config.Requirements, config.GetBaseURL()),
		config:            config,
		email:             emailConfig,
		meta:              meta,
	}
}

func (env *environment) Mount(group *echo.Group) {
	csrf := echoMiddleware.CSRF()
	throttleDuration, _ := time.ParseDuration(env.config.Login.Settings.ThrottleDuration)
	throttleMiddleware := middleware.NewThrottleGroup(1, throttleDuration)
	recaptchaMiddleware := buildRecaptchaMiddleware(&env.config.RecaptchaV2)

	{ // Insecure routes
		group.POST("/login", env.routeLogin, throttleMiddleware, csrf)
		group.POST("/logout", env.routeLogout, csrf)

		group.POST("/stipulation", env.routeTokenStipulation, throttleMiddleware, csrf)

		if env.config.Login.OneTime.Enabled {
			group.GET("/onetime", env.routeOneTimeAuth, throttleMiddleware)
			if env.config.Login.OneTime.AllowForgotPassword {
				group.POST("/onetime", env.routeOneTimePost, common.CoalesceMiddleware(throttleMiddleware, recaptchaMiddleware, csrf)...)
			}
		}
	}

	{ // Secure routes
		authProvider := selector.NewSelectorMiddleware(
			auth.NewSessionAuthProvider(&env.config.Login.Cookie.JWT, csrf),
			selector.HandlerUnauthorized(),
		)

		group.GET("/account", env.routeAccount, authProvider)
		group.GET("/account/audit", env.routeAccountAudit, authProvider)
		group.POST("/account/password", env.routeChangePassword, authProvider)
		group.GET("/account/password", env.routeChangePasswordRequirements, authProvider)
	}

}

func buildRecaptchaMiddleware(config *config.ConfigRecaptchaV2) echo.MiddlewareFunc {
	if config == nil || !config.Enabled {
		return nil
	}
	return recaptcha.MiddlewareV2(config.Secret)
}
