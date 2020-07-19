package ui

import (
	"simple-auth/pkg/config"
	"simple-auth/pkg/db"
	"simple-auth/pkg/routes/common"
	"simple-auth/pkg/routes/middleware"

	"github.com/labstack/echo"
	echoMiddleware "github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
)

type environment struct {
	db     db.SADB
	config *config.ConfigWeb
	email  *config.ConfigEmail
}

func NewController(db db.SADB, config *config.ConfigWeb, emailConfig *config.ConfigEmail) common.Controller {
	return &environment{
		db:     db,
		config: config,
		email:  emailConfig,
	}
}

func (env *environment) Mount(group *echo.Group) {
	group.Use(echoMiddleware.CSRF())

	group.POST("/account", env.routeCreateAccount)
	group.POST("/login", env.routeLogin)
	group.POST("/logout", env.routeLogout)

	if env.config.Login.Cookie.JWT.SigningKey != "" {
		manageGroup := group.Group("/manage/")
		manageGroup.Use(middleware.LoggedInMiddleware(&env.config.Login.Cookie.JWT))
		manageGroup.GET("", env.routeAccount)
	} else {
		logrus.Warn("No JWT secret specified, refusing to bind user management endpoints")
	}
}
