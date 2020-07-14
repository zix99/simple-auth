package main

import (
	"net/http"
	"simple-auth/pkg/config"
	"simple-auth/pkg/db"
	"simple-auth/pkg/routes/api/auth"
	"simple-auth/pkg/routes/api/ui"

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

func buildTemplateContext(web *config.ConfigWeb) map[string]interface{} {
	context := make(map[string]interface{})
	for k, v := range web.Metadata {
		context[k] = v
	}
	context["Requirements"] = web.Requirements
	context["RecaptchaV2"] = struct {
		Enabled bool
		SiteKey string
		Theme   string
	}{web.RecaptchaV2.Enabled, web.RecaptchaV2.SiteKey, web.RecaptchaV2.Theme}
	return context
}

func simpleAuthServer(config *config.Config) error {
	if config.Production {
		logrus.Info("Running in production mode")
	}
	if config.Web.JWT.Secret == "" {
		logrus.Warn("No web.jwt.secret is set, user will not be able to login")
	}

	// Dependencies
	env := &environment{
		db: db.New(config.Db.Driver, config.Db.URL),
	}

	e := echo.New()
	e.Debug = !config.Production

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.AddTrailingSlash())

	// Static app router
	e.Renderer = newTemplateRenderer()
	e.Static("/static", "./static")
	e.Static("/dist", "./dist")

	// Health
	e.GET("/health", env.routeHealth)

	context := buildTemplateContext(&config.Web)
	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "home", context)
	})
	e.GET("/create", func(c echo.Context) error {
		return c.Render(http.StatusOK, "createAccount", context)
	})
	e.GET("/manage", func(c echo.Context) error {
		return c.Render(http.StatusOK, "manageAccount", context)
	})

	// Attach routes
	auth.NewController(env.db, &config.Authenticators).Mount(e.Group("/api/v1/auth"))
	ui.NewController(env.db, &config.Web).Mount(e.Group("/api/ui"))

	// Start
	logrus.Infof("Starting server on http://%v", config.Web.Host)
	return e.Start(config.Web.Host)
}

func main() {
	logrus.Fatal(simpleAuthServer(config.Global))
}
