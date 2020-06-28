package ui

import (
	"simple-auth/pkg/db"

	"github.com/labstack/echo"
)

type environment struct {
	db db.SADB
}

func NewRouter(group *echo.Group, db db.SADB) {
	env := &environment{
		db: db,
	}

	group.POST("/account", env.routeCreateAccount)
}
