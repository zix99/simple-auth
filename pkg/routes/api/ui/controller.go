package ui

import (
	"simple-auth/pkg/config"
	"simple-auth/pkg/db"
	"simple-auth/pkg/routes/api/ui/recaptcha"
	"simple-auth/pkg/routes/common"
	"simple-auth/pkg/routes/middleware"
	"simple-auth/pkg/routes/middleware/selector"
	"simple-auth/pkg/routes/middleware/selector/auth"
	"simple-auth/pkg/services"
	"time"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

type environment struct {
	db                db.SADB
	localLoginService services.LocalLoginService

	meta   *config.ConfigMetadata
	config *config.ConfigWeb
	email  *config.ConfigEmail
}

func NewController(db db.SADB, meta *config.ConfigMetadata, config *config.ConfigWeb, emailConfig *config.ConfigEmail) common.Controller {
	return &environment{
		db:                db,
		localLoginService: services.NewLocalLoginService(db, &config.Login.TwoFactor),
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
		if env.config.Login.Settings.CreateAccountEnabled {
			group.POST("/account", env.routeCreateAccount, throttleMiddleware, csrf)
		}
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

		if env.config.Login.TwoFactor.Enabled {
			group.GET("/2fa", env.routeSetup2FA, authProvider)
			group.GET("/2fa/qrcode", env.route2FAQRCodeImage, authProvider)
			group.POST("/2fa", env.routeConfirm2FA, authProvider)
			group.DELETE("/2fa", env.routeDeactivate2FA, authProvider)
		}
	}

}

func buildRecaptchaMiddleware(config *config.ConfigRecaptchaV2) echo.MiddlewareFunc {
	if config == nil || !config.Enabled {
		return nil
	}
	return recaptcha.MiddlewareV2(config.Secret)
}
