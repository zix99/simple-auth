package ui

import (
	"simple-auth/pkg/config"
	"simple-auth/pkg/db"

	"github.com/labstack/echo"
)

type environment struct {
	db     db.SADB
	config *config.ConfigWeb
}

func NewRouter(group *echo.Group, db db.SADB, config *config.ConfigWeb) {
	env := &environment{
		db:     db,
		config: config,
	}

	group.POST("/account", env.routeCreateAccount)
}
