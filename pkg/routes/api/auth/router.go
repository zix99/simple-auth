package auth

import (
	"simple-auth/pkg/config"
	"simple-auth/pkg/db"
	"simple-auth/pkg/routes/common"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

type environment struct {
	db     db.SADB
	config *config.ConfigAuthencatorSet
}

func NewController(db db.SADB, config *config.ConfigAuthencatorSet) common.Controller {
	return &environment{
		db:     db,
		config: config,
	}
}

func (env *environment) Mount(group *echo.Group) {
	logrus.Info("Setting up auth APIs...")

	if env.config.Token.Enabled {
		logrus.Info("Enabling session auth...")
		group.POST("/token", env.routeIssueSessionToken)
		group.POST("/token/session", env.routeIssueVerificationToken)
		group.POST("/token/session/verify", env.routeVerifyToken)
	}
	if env.config.Simple.Enabled {
		logrus.Info("Enabling simple auth...")
		group.POST("/simple", env.routeSimpleAuthenticate)
	}
}
