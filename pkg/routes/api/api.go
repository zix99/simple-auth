package api

import (
	"simple-auth/pkg/config"
	"simple-auth/pkg/db"
	authAPI "simple-auth/pkg/routes/api/auth"
	"simple-auth/pkg/routes/api/ui"
	v1 "simple-auth/pkg/routes/api/v1"
	saMiddleware "simple-auth/pkg/routes/middleware"
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
		// API
		{
			publicAuth := buildPublicAuthMiddleware(&config.API)
			v1Env := v1.NewEnvironment(&config.Web, db)
			v1api.POST("/account/check", v1Env.RouteCheckUsername, publicAuth)
		}

		// Attach authenticator routes
		if config.Authenticators.Token.Enabled {
			route := v1api.Group("/auth/token")
			authAPI.NewTokenAuthController(db, &config.Authenticators.Token).Mount(route)
		}
		if config.Authenticators.Simple.Enabled {
			route := v1api.Group("/auth/simple")
			authAPI.NewSimpleAuthController(services.NewLocalLoginService(db, &config.Web.Login.TwoFactor), &config.Authenticators.Simple).Mount(route)
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

func buildPublicAuthMiddleware(config *config.ConfigAPI) echo.MiddlewareFunc {
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
	selectorGroups = append(selectorGroups, selector.NewSelectorGroup(
		selector.SelectorAlways,
		middleware.CSRF(),
		throttleMiddleware,
	))

	return selector.NewSelectorMiddleware(selectorGroups...)
}
