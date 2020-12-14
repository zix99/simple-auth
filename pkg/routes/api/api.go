package api

import (
	"simple-auth/pkg/config"
	"simple-auth/pkg/db"
	"simple-auth/pkg/email"
	authAPI "simple-auth/pkg/routes/api/auth"
	"simple-auth/pkg/routes/api/ui"
	v1 "simple-auth/pkg/routes/api/v1"
	saMiddleware "simple-auth/pkg/routes/middleware"
	"simple-auth/pkg/routes/middleware/recaptcha"
	"simple-auth/pkg/routes/middleware/selector"
	"simple-auth/pkg/routes/middleware/selector/auth"
	"simple-auth/pkg/services"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

func MountAPI(e *echo.Group, config *config.Config, db db.SADB) {
	v1api := e.Group("/v1")
	{
		// Public API (eg. from the UI)
		v1Env := v1.NewEnvironment(config, db)
		{
			publicAuth := buildPublicAuthMiddleware(&config.API, nil)
			publicAuthWithRecaptcha := buildPublicAuthMiddleware(&config.API, &config.Web.RecaptchaV2)
			v1api.POST("/account/check", v1Env.RouteCheckUsername, publicAuth)
			if config.Web.Login.Settings.CreateAccountEnabled {
				v1api.POST("/account", v1Env.RouteCreateAccount, publicAuthWithRecaptcha)
			}

			v1api.POST("/stipulation", v1Env.RouteSatisfyTokenStipulation, publicAuth)
		}
		{
			privateAuth := buildPrivateAuthMiddleware(&config.Web.Login.Cookie, &config.API)
			v1api.GET("/account", v1Env.RouteGetAccount, privateAuth)
			v1api.GET("/account/audit", v1Env.RouteGetAccountAudit, privateAuth)
			if config.Web.Login.TwoFactor.Enabled {
				v1api.GET("/2fa", v1Env.RouteSetup2FA, privateAuth)
				v1api.GET("/2fa/qrcode", v1Env.Route2FAQRCodeImage, privateAuth)
				v1api.POST("/2fa", v1Env.RouteConfirm2FA, privateAuth)
				v1api.DELETE("/2fa", v1Env.RouteDeactivate2FA, privateAuth)
			}
		}

		// Attach authenticator routes
		if config.Authenticators.Token.Enabled {
			route := v1api.Group("/auth/token")
			authAPI.NewTokenAuthController(db, &config.Authenticators.Token).Mount(route)
		}
		if config.Authenticators.Simple.Enabled {
			route := v1api.Group("/auth/simple")
			emailService := email.NewFromConfig(logrus.StandardLogger(), &config.Email)
			loginService := services.NewLocalLoginService(emailService, &config.Metadata, &config.Web.Login.TwoFactor, &config.Web.Requirements, config.Web.GetBaseURL())
			authAPI.NewSimpleAuthController(loginService, &config.Authenticators.Simple).Mount(route)
		}
		if config.Authenticators.Vouch.Enabled {
			route := v1api.Group("/auth/vouch")
			authAPI.NewVouchAuthController(db, &config.Authenticators.Vouch, &config.Web.Login.Cookie.JWT).Mount(route)
		}
	}

	// LEGACY: Attach UI/access router
	{
		uiGroup := e.Group("/ui")
		controller := ui.NewController(db, &config.Metadata, &config.Web, &config.Email)
		controller.Mount(uiGroup)
	}
}

func buildPublicAuthMiddleware(config *config.ConfigAPI, recaptcha *config.ConfigRecaptchaV2) echo.MiddlewareFunc {
	var selectorGroups []selector.SelectorGroup

	if config.External {
		logrus.Info("Enabling external API...")
		if config.SharedSecret == "" {
			logrus.Fatal("Invalid shared-secret for external api")
		}
		selectorGroups = append(selectorGroups,
			selector.NewSelectorGroup(auth.SharedSecretSelector(config.SharedSecret)),
		)
	}

	throttleDuration, _ := time.ParseDuration(config.ThrottleDuration)
	throttleMiddleware := saMiddleware.NewThrottleGroup(1, throttleDuration)
	recaptchaMiddleware := buildRecaptchaMiddleware(recaptcha)
	selectorGroups = append(selectorGroups, selector.NewSelectorGroup(
		selector.SelectorAlways,
		middleware.CSRF(),
		throttleMiddleware,
		recaptchaMiddleware,
	))

	return selector.NewSelectorMiddleware(selectorGroups...)
}

func buildPrivateAuthMiddleware(sessionConfig *config.ConfigLoginCookie, apiConfig *config.ConfigAPI) echo.MiddlewareFunc {
	var selectorGroups []selector.SelectorGroup

	csrf := middleware.CSRF()
	selectorGroups = append(selectorGroups, auth.NewSessionAuthProvider(&sessionConfig.JWT, csrf))

	if apiConfig.External {
		selectorGroups = append(selectorGroups, auth.NewSharedSecretWithAccountAuth(apiConfig.SharedSecret))
	}

	selectorGroups = append(selectorGroups, selector.HandlerUnauthorized())

	return selector.NewSelectorMiddleware(selectorGroups...)
}

func buildRecaptchaMiddleware(config *config.ConfigRecaptchaV2) echo.MiddlewareFunc {
	if config == nil || !config.Enabled {
		return nil
	}
	return recaptcha.MiddlewareV2(config.Secret)
}
