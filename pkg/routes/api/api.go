package api

import (
	"simple-auth/pkg/appcontext"
	"simple-auth/pkg/config"
	"simple-auth/pkg/db"
	"simple-auth/pkg/email"
	authAPI "simple-auth/pkg/routes/api/auth"
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

// @title Simple Auth API
// @version 1.0
// @description The simple auth API
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host example.com
// @BasePath /api/v1
// @query.collection.format multi

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @securityDefinitions.apikey SessionAuth
// @in cookie
// @name auth

func MountAPI(e *echo.Group, config *config.Config, db db.SADB) {
	transactional := appcontext.Transaction()

	emailService := email.NewFromConfig(&config.Email)
	loginService := services.NewLocalLoginService(emailService, &config.Metadata, &config.Providers.Local, config.Web.GetBaseURL())
	oAuthController := authAPI.NewOAuth2Controller(&config.Authenticators.OAuth2, loginService)

	v1api := e.Group("/v1")
	{
		// Public API (eg. from the UI)
		v1Env := v1.NewEnvironment(config)
		{
			publicAuth := buildPublicAuthMiddleware(&config.API, nil, true)
			publicAuthWithRecaptcha := buildPublicAuthMiddleware(&config.API, &config.Web.RecaptchaV2, true)
			publicAuthNoDelay := buildPublicAuthMiddleware(&config.API, nil, false)

			v1api.POST("/account/check", v1Env.RouteCheckUsername, publicAuth)
			if config.Providers.Settings.CreateAccountEnabled {
				v1api.POST("/account", v1Env.RouteCreateAccount, publicAuthWithRecaptcha, transactional)
			}

			v1api.POST("/stipulation", v1Env.RouteSatisfyTokenStipulation, publicAuth)

			v1api.POST("/auth/session", v1Env.RouteSessionLogin, publicAuth)
			v1api.DELETE("/auth/session", v1Env.RouteSessionLogout, publicAuthNoDelay)

			if config.Web.Login.OneTime.Enabled {
				v1api.GET("/auth/onetime", v1Env.RouteOneTimeAuth, publicAuth)
				if config.Web.Login.OneTime.AllowForgotPassword {
					v1api.POST("/auth/onetime", v1Env.RouteOneTimeCreateToken, publicAuthWithRecaptcha)
				}
			}
		}

		// Private auth
		{
			privateAuth := buildPrivateAuthMiddleware(&config.Web.Login.Cookie, &config.API)
			v1api.GET("/account", v1Env.RouteGetAccount, privateAuth)
			v1api.GET("/account/audit", v1Env.RouteGetAccountAudit, privateAuth)

			v1api.GET("/local", v1Env.RouteGetLocalLogin, privateAuth)
			v1api.POST("/local/password", v1Env.RouteChangePassword, privateAuth)
			if config.Providers.Local.TwoFactor.Enabled {
				v1api.GET("/local/2fa", v1Env.RouteSetup2FA, privateAuth)
				v1api.GET("/local/2fa/qrcode", v1Env.Route2FAQRCodeImage, privateAuth)
				v1api.POST("/local/2fa", v1Env.RouteConfirm2FA, privateAuth)
				v1api.DELETE("/local/2fa", v1Env.RouteDeactivate2FA, privateAuth)
			}

			v1api.GET("/auth/oauth2", oAuthController.RouteGetTokensForUser, privateAuth)
			v1api.DELETE("/auth/oauth2/token", oAuthController.RouteRevokeToken, privateAuth)
			if config.Authenticators.OAuth2.WebGrant {
				v1api.POST("/auth/oauth2/grant", oAuthController.RouteAuthorizedGrantCode, privateAuth, transactional)
			}
		}

		// Attach authenticator routes
		{
			if config.Authenticators.Simple.Enabled {
				route := v1api.Group("/auth/simple")
				authAPI.NewSimpleAuthController(loginService, &config.Authenticators.Simple).Mount(route)
			}
			if config.Authenticators.Vouch.Enabled {
				route := v1api.Group("/auth/vouch")
				authAPI.NewVouchAuthController(db, &config.Authenticators.Vouch, &config.Web.Login.Cookie).Mount(route)
			}
			{
				v1api.GET("/auth/oauth2/client/:client_id", oAuthController.RouteClientInfo)
				v1api.POST("/auth/oauth2/token", oAuthController.RouteTokenGrant, transactional)
				v1api.POST("/auth/oauth2/token_info", oAuthController.RouteIntrospectToken)
			}
		}
	}
}

// Authentication that allows either UI access (with CSRF, throttled, and optional recaptcha), or private-api-key access
func buildPublicAuthMiddleware(config *config.ConfigAPI, recaptcha *config.ConfigRecaptchaV2, throttle bool) echo.MiddlewareFunc {
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
	var throttleMiddleware echo.MiddlewareFunc
	if throttle {
		throttleMiddleware = saMiddleware.NewThrottleGroup(1, throttleDuration)
	}
	recaptchaMiddleware := buildRecaptchaMiddleware(recaptcha)
	selectorGroups = append(selectorGroups, selector.NewSelectorGroup(
		selector.SelectorAlways,
		middleware.CSRF(),
		throttleMiddleware,
		recaptchaMiddleware,
	))

	return selector.NewSelectorMiddleware(selectorGroups...)
}

// Allow a session token, or private key access
func buildPrivateAuthMiddleware(sessionConfig *config.ConfigLoginCookie, apiConfig *config.ConfigAPI) echo.MiddlewareFunc {
	var selectorGroups []selector.SelectorGroup

	csrf := middleware.CSRF()
	selectorGroups = append(selectorGroups, auth.NewSessionAuthProvider(sessionConfig, csrf))

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
