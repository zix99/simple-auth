package auth

import (
	"simple-auth/pkg/config"
	"simple-auth/pkg/db"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

type environment struct {
	db db.SADB
}

func NewRouter(group *echo.Group, db db.SADB, config *config.ConfigAuthencatorSet) {
	env := &environment{
		db: db,
	}

	logrus.Info("Setting up auth APIs...")

	if config.Token.Enabled {
		setupSessionAuthenticator(env, group.Group("/token"))
	}
	if config.Simple.Enabled {
		setupSimpleAuthenticator(env, group.Group("/simple"))
	}
}
