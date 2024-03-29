package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"simple-auth/pkg/appcontext"
	"simple-auth/pkg/box"
	"simple-auth/pkg/box/echobox"
	"simple-auth/pkg/config"
	"simple-auth/pkg/db"
	"simple-auth/pkg/routes/api"
	"simple-auth/pkg/routes/api/providers"
	saMiddleware "simple-auth/pkg/routes/middleware"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/acme/autocert"
)

type healthResponse struct {
	Status string
	DB     bool
}

func routeHealth(db db.SADB) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(200, healthResponse{
			Status: "OK",
			DB:     db.IsAlive(),
		})
	}
}

func simpleAuthServer(config *config.Config) error {
	log := logrus.New()

	if config.Help {
		fmt.Println()
		fmt.Println("Help:")
		fmt.Println("http://simple-auth.zdyn.net/")
		fmt.Printf("Version: %s:%s\n", version, buildSha)
		return nil
	}
	if config.Version {
		fmt.Println(version)
		return nil
	}

	if config.Production {
		logrus.Info("Running in production mode")
	}
	if config.Web.Login.Cookie.JWT.SigningKey == "" {
		logrus.Warn("No web.login.cookie.jwt.signingkey is set, user will not be able to login")
	}

	box.Global.Verbose = config.Verbose
	box.Global.CheckDisk = config.StaticFromDisk
	if config.Verbose {
		log.SetLevel(logrus.DebugLevel)
		log.Debugln("Debug enabled")
	}

	// Dependencies
	db := db.New(config.Db.Driver, config.Db.URL)
	db.EnableLogging(config.Db.Debug)

	e := echo.New()
	e.Debug = !config.Production
	e.Validator = NewGoPlaygroundValidator()

	watchForCleanShutdown(e)

	applyHooks(e, config)

	e.Use(middleware.Recover())
	e.Use(appcontext.WithLogger(log).Middleware())
	e.Use(saMiddleware.NewCorrelationMiddleware(false, true))
	e.Use(saMiddleware.NewRequestLoggerMiddleware())
	e.Use(appcontext.WithSADB(db).Middleware())

	// Gateway
	if config.Web.Gateway.Enabled {
		log.Infof("Enabling authentication gateway: %v", config.Web.Gateway.Targets)
		e.Use(saMiddleware.AuthenticationGateway(&config.Web.Gateway, &config.Web.Login.Cookie))
	}

	// Static app router
	e.Renderer = newTemplateRenderer(!config.Production)
	e.GET("/static/*", echobox.Static("./static"))
	e.GET("/dist/*", echobox.Static("./dist"))

	// Health
	e.GET("/health", routeHealth(db))

	// UI
	newUIController(&config.Web, &config.Metadata, &config.Providers).Mount(e.Group(""))

	// API
	api.MountAPI(e.Group("/api"), config, db)

	// OIDC Controllers
	{
		oidcGroup := e.Group("/oidc")
		for _, oidc := range config.Providers.OIDC {
			oidcController := providers.NewOIDCController(config.Web.GetBaseURL()+"/oidc", oidc.ID, &config.Providers.Settings, &config.Web.Login.Settings, oidc, &config.Web.Login.Cookie, db)
			oidcController.Mount(oidcGroup)
		}
	}

	// Well known routes
	e.GET("/onetime", redirectHandler("/api/v1/auth/onetime"))
	e.GET("/oauth2", redirectVue("oauth2"))

	// Start
	log.Infof("Starting server on http://%v", config.Web.Host)
	if config.Web.TLS.Enabled {
		if config.Web.TLS.Auto {
			if len(config.Web.TLS.AutoHosts) > 0 {
				e.AutoTLSManager.HostPolicy = autocert.HostWhitelist(config.Web.TLS.AutoHosts...)
			}
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
	err := simpleAuthServer(config.Load(os.Args[1:]...))
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logrus.Fatal(err)
	}
}
