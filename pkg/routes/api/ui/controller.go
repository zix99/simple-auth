package ui

import (
	"simple-auth/pkg/config"
	"simple-auth/pkg/db"
	"simple-auth/pkg/routes/api/ui/recaptcha"
	"simple-auth/pkg/routes/common"
	"simple-auth/pkg/routes/middleware"
	"time"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

type environment struct {
	db     db.SADB
	meta   *config.ConfigMetadata
	config *config.ConfigWeb
	email  *config.ConfigEmail
}

func NewController(db db.SADB, meta *config.ConfigMetadata, config *config.ConfigWeb, emailConfig *config.ConfigEmail) common.Controller {
	return &environment{
		db:     db,
		config: config,
		email:  emailConfig,
		meta:   meta,
	}
}

func (env *environment) Mount(group *echo.Group) {
	group.Use(echoMiddleware.CSRF())

	throttleDuration, _ := time.ParseDuration(env.config.Login.Settings.ThrottleDuration)
	throttleMiddleware := middleware.NewThrottleGroup(1, throttleDuration)
	recaptchaMiddleware := buildRecaptchaMiddleware(&env.config.RecaptchaV2)

	if env.config.Login.Settings.CreateAccountEnabled {
		group.POST("/account", env.routeCreateAccount, throttleMiddleware)
	}
	group.POST("/login", env.routeLogin, throttleMiddleware)
	group.POST("/logout", env.routeLogout)

	loggedIn := middleware.LoggedInMiddleware(&env.config.Login.Cookie.JWT)
	group.GET("/account", env.routeAccount, loggedIn)
	group.GET("/account/audit", env.routeAccountAudit, loggedIn)
	group.POST("/account/password", env.routeChangePassword, loggedIn)
	group.GET("/account/password", env.routeChangePasswordRequirements, loggedIn)

	if env.config.Login.OneTime.Enabled {
		group.GET("/onetime", env.routeOneTimeAuth, throttleMiddleware)
		if env.config.Login.OneTime.AllowForgotPassword {
			group.POST("/onetime", env.routeOneTimePost, common.CoalesceMiddleware(throttleMiddleware, recaptchaMiddleware)...)
		}
	}
}

func buildRecaptchaMiddleware(config *config.ConfigRecaptchaV2) echo.MiddlewareFunc {
	if config == nil || !config.Enabled {
		return nil
	}
	return recaptcha.MiddlewareV2(config.Secret)
}
