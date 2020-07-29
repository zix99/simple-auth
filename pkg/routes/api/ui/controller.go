package ui

import (
	"simple-auth/pkg/config"
	"simple-auth/pkg/db"
	"simple-auth/pkg/routes/common"
	"simple-auth/pkg/routes/middleware"
	"time"

	"github.com/labstack/echo"
	echoMiddleware "github.com/labstack/echo/middleware"
)

type environment struct {
	db     db.SADB
	meta   *config.ConfigMetadata
	config *config.ConfigWeb
	email  *config.ConfigEmail
}

func NewController(db db.SADB, meta *config.ConfigMetadata, config *config.ConfigWeb, emailConfig *config.ConfigEmail) common.Controller {
	return &environment{
		db:     db,
		config: config,
		email:  emailConfig,
		meta:   meta,
	}
}

func (env *environment) Mount(group *echo.Group) {
	group.Use(echoMiddleware.CSRF())

	delayGroup := middleware.NewThrottleGroup(1, 1*time.Second)

	if env.config.Login.CreateAccountEnabled {
		group.POST("/account", env.routeCreateAccount, delayGroup)
	}
	group.POST("/login", env.routeLogin, delayGroup)
	group.POST("/logout", env.routeLogout)

	loggedIn := middleware.LoggedInMiddleware(&env.config.Login.Cookie.JWT)
	group.GET("/account", env.routeAccount, loggedIn)
	group.GET("/account/audit", env.routeAccountAudit, loggedIn)
}
