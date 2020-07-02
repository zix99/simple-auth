package main

import (
	"net/http"
	"simple-auth/pkg/api/auth"
	"simple-auth/pkg/api/ui"
	"simple-auth/pkg/config"
	"simple-auth/pkg/db"

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

func buildTemplateContext() map[string]interface{} {
	context := make(map[string]interface{})
	for k, v := range config.Global.Web.Metadata {
		context[k] = v
	}
	context["Requirements"] = config.Global.Web.Requirements
	return context
}

func simpleAuthServer(config *config.Config) error {
	if config.Production {
		// TODO
	}

	// Dependencies
	env := &environment{
		db: db.New(config.Db.Driver, config.Db.URL),
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.AddTrailingSlash())

	// Static app router
	e.Renderer = newTemplateSet()
	e.Static("/static", "./static")
	e.Static("/dist", "./dist")

	// Health
	e.GET("/health", env.routeHealth)

	context := buildTemplateContext()
	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "home", context)
	})
	e.GET("/create", func(c echo.Context) error {
		return c.Render(http.StatusOK, "createAccount", context)
	})

	// Attach routes
	auth.NewRouter(e.Group("/api/v1/auth"), env.db, &config.Authenticators)
	ui.NewRouter(e.Group("/api/ui"), env.db)

	// Start
	logrus.Infof("Starting server on http://%v", config.Web.Host)
	return e.Start(config.Web.Host)
}

func main() {
	logrus.Fatal(simpleAuthServer(config.Global))
}
