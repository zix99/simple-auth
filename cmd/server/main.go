package main

import (
	"simple-auth/pkg/api/auth"
	"simple-auth/pkg/db"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type environment struct {
	db *db.DB
}

func (env *environment) routeHealth(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "OK",
		"db":     env.db.IsAlive(),
	})
}

func simpleAuthServer(config *config) error {
	if config.Production {
		gin.SetMode("release")
	}

	env := &environment{
		db: db.New(config.Db.Driver, config.Db.URL),
	}

	r := gin.Default()
	r.Static("/static", "./static")

	// Static app router
	r.Static("img", "./ui/dist/img")
	r.Static("js", "./ui/dist/js")
	r.StaticFile("favicon.ico", "./ui/dist/favicon.ico")
	r.StaticFile("index.html", "./ui/dist/index.html")
	r.StaticFile("/", "./ui/dist/index.html")

	r.GET("/health", env.routeHealth)

	auth.NewRouter(r.Group("/api/v1/auth"), env.db)

	logrus.Infof("Starting server on http://%v", config.Web.Host)
	r.Run(config.Web.Host)
	return nil
}

func main() {
	config := readConfig()
	simpleAuthServer(config)
}
