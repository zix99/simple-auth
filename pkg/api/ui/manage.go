package ui

import (
	"net/http"
	"simple-auth/pkg/api/common"

	"github.com/sirupsen/logrus"

	"github.com/labstack/echo"
)

func newManagementRouter(env *environment, group *echo.Group) {
	if env.config.JWT.Secret == "" {
		logrus.Warn("No JWT secret specified, refusing to bind user management endpoints")
		return
	}

	group.Use(loggedInMiddleware(env.config.JWT.Secret))
	group.GET("", env.routeAccount)
}

func (env *environment) routeAccount(c echo.Context) error {
	accountUUID := c.Get("accountUUID").(string)
	logrus.Infof("Get account for %s", accountUUID)
	account, err := env.db.FindAccount(accountUUID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, common.JsonErrorf("Logged in with unknown account"))
	}

	username, _ := env.db.FindSimpleAuthUsername(account)

	return c.JSON(http.StatusOK, common.Json{
		"id":      account.UUID,
		"created": account.CreatedAt,
		"email":   account.Email,
		"auth": common.Json{
			"simple": common.Json{
				"username": username,
			},
		},
	})
}
