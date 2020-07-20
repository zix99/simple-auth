package main

import (
	"simple-auth/pkg/config"
	"simple-auth/pkg/db"
	"simple-auth/pkg/routes/api/auth"
	"simple-auth/pkg/routes/api/ui"
	logMiddleware "simple-auth/pkg/routes/middleware"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
)

type environment struct {
	db db.SADB
}

type healthResponse struct {
	Status string
	DB     bool
}

func (env *environment) routeHealth(c echo.Context) error {
	return c.JSON(200, healthResponse{
		Status: "OK",
		DB:     env.db.IsAlive(),
	})
}

func simpleAuthServer(config *config.Config) error {
	if config.Production {
		logrus.Info("Running in production mode")
	}
	if config.Web.Login.Cookie.JWT.SigningKey == "" {
		logrus.Warn("No web.login.cookie.jwt.signingkey is set, user will not be able to login")
	}

	// Dependencies
	env := &environment{
		db: db.New(config.Db.Driver, config.Db.URL),
	}

	e := echo.New()
	e.Debug = !config.Production

	e.Use(logMiddleware.NewLoggerMiddleware())
	e.Use(middleware.Recover())
	e.Use(middleware.AddTrailingSlash())

	// Static app router
	e.Renderer = newTemplateRenderer(!config.Production)
	e.Static("/static", "./static")
	e.Static("/dist", "./dist")

	// Health
	e.GET("/health", env.routeHealth)

	// UI
	newUIController(&config.Web, &config.Metadata).Mount(e.Group(""))

	// Attach authenticator routes
	if config.Authenticators.Token.Enabled {
		route := e.Group("/api/v1/auth/token")
		auth.NewTokenAuthController(env.db, &config.Authenticators.Token).Mount(route)
	}
	if config.Authenticators.Simple.Enabled {
		route := e.Group("/api/v1/auth/simple")
		auth.NewSimpleAuthController(env.db).Mount(route)
	}

	// Attach UI/access router
	ui.NewController(env.db, &config.Metadata, &config.Web, &config.Email).Mount(e.Group("/api/ui"))

	// Start
	logrus.Infof("Starting server on http://%v", config.Web.Host)
	return e.Start(config.Web.Host)
}

func main() {
	logrus.Fatal(simpleAuthServer(config.Global))
}
