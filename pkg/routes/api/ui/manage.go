package ui

import (
	"net/http"
	"simple-auth/pkg/routes/common"

	"github.com/sirupsen/logrus"

	"github.com/labstack/echo"
)

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
