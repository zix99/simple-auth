package ui

import (
	"simple-auth/pkg/config"
	"simple-auth/pkg/db"
	"simple-auth/pkg/routes/common"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

type environment struct {
	db     db.SADB
	config *config.ConfigWeb
}

func NewRouter(db db.SADB, config *config.ConfigWeb) common.Controller {
	return &environment{
		db:     db,
		config: config,
	}
}

func (env *environment) Mount(group *echo.Group) {
	group.POST("/account", env.routeCreateAccount)
	group.POST("/login", env.routeLogin)
	group.POST("/logout", env.routeLogout)

	if env.config.JWT.Secret != "" {
		manageGroup := group.Group("/manage/")
		manageGroup.Use(loggedInMiddleware(env.config.JWT.Secret))
		manageGroup.GET("", env.routeAccount)
	} else {
		logrus.Warn("No JWT secret specified, refusing to bind user management endpoints")
	}
}
