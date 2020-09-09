package main

import (
	"simple-auth/pkg/box"
	"simple-auth/pkg/box/echobox"
	"simple-auth/pkg/config"
	"simple-auth/pkg/db"
	"simple-auth/pkg/routes/api/auth"
	"simple-auth/pkg/routes/api/providers"
	"simple-auth/pkg/routes/api/ui"
	saMiddleware "simple-auth/pkg/routes/middleware"

	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/acme/autocert"
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

	box.Global.Verbose = config.Verbose
	box.Global.CheckDisk = config.StaticFromDisk

	// Dependencies
	env := &environment{
		db: db.New(config.Db.Driver, config.Db.URL),
	}
	env.db.EnableLogging(config.Db.Debug)

	e := echo.New()
	e.Debug = !config.Production

	e.Use(saMiddleware.NewCorrelationMiddleware(false, true))
	e.Use(saMiddleware.NewLoggerMiddleware())
	e.Use(middleware.Recover())

	// Prometheus
	if config.Web.Prometheus {
		p := prometheus.NewPrometheus("sa", nil)
		p.Use(e)
	}

	// Gateway
	if config.Web.Gateway.Enabled {
		logrus.Infof("Enabling authentication gateway: %v", config.Web.Gateway.Targets)
		e.Use(saMiddleware.AuthenticationGateway(&config.Web.Gateway, &config.Web.Login.Cookie))
	}

	// Static app router
	e.Renderer = newTemplateRenderer(!config.Production)
	e.GET("/static/*", echobox.Static("./static"))
	e.GET("/dist/*", echobox.Static("./dist"))

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
		auth.NewSimpleAuthController(env.db, &config.Authenticators.Simple).Mount(route)
	}

	// Attach UI/access router
	ui.NewController(env.db, &config.Metadata, &config.Web, &config.Email).Mount(e.Group("/api/ui"))

	// OIDC Controllers
	{
		oidcGroup := e.Group("/oidc")
		for _, oidc := range config.Web.Login.OIDC {
			oidcController := providers.NewOIDCController(config.Web.GetBaseURL()+"/oidc", oidc.ID, &config.Web.Login.Settings, oidc, &config.Web.Login.Cookie, env.db)
			oidcController.Mount(oidcGroup)
		}
	}

	// Well known routes
	e.GET("/onetime", redirectHandler("/api/ui/onetime"))

	// Start
	logrus.Infof("Starting server on http://%v", config.Web.Host)
	if config.Web.TLS.Enabled {
		if config.Web.TLS.Auto {
			e.AutoTLSManager.Cache = autocert.DirCache(config.Web.TLS.Cache)
			return e.StartAutoTLS(config.Web.Host)
		} else {
			return e.StartTLS(config.Web.Host, config.Web.TLS.CertFile, config.Web.TLS.KeyFile)
		}
	} else {
		return e.Start(config.Web.Host)
	}
}

func main() {
	logrus.Fatal(simpleAuthServer(config.Load(true)))
}
