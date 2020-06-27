package main

import (
	"simple-auth/pkg/api/auth"
	"simple-auth/pkg/api/ui"
	"simple-auth/pkg/config"
	"simple-auth/pkg/db"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type environment struct {
	db db.SADB
}

func (env *environment) routeHealth(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "OK",
		"db":     env.db.IsAlive(),
	})
}

func buildTemplateContext() gin.H {
	context := make(map[string]interface{})
	for k, v := range config.Global.Web.Metadata {
		context[k] = v
	}
	context["Requirements"] = config.Global.Web.Requirements
	return context
}

func simpleAuthServer(config *config.Config) error {
	if config.Production {
		gin.SetMode("release")
	}

	// Dependencies
	env := &environment{
		db: db.New(config.Db.Driver, config.Db.URL),
	}

	r := gin.Default()

	// Static app router
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")
	r.Static("/dist", "./dist")

	// Health
	r.GET("/health", env.routeHealth)

	context := buildTemplateContext()
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "createAccount.tmpl", context)
	})

	// Attach routes
	auth.NewRouter(r.Group("/api/v1/auth"), env.db, &config.Authenticators)
	ui.NewRouter(r.Group("/api/ui"), env.db)

	// Start
	logrus.Infof("Starting server on http://%v", config.Web.Host)
	r.Run(config.Web.Host)
	return nil
}

func main() {
	simpleAuthServer(config.Global)
}
