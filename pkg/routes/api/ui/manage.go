package ui

import (
	"net/http"
	"simple-auth/pkg/routes/common"
	"simple-auth/pkg/routes/middleware"
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/labstack/echo"
)

func (env *environment) routeAccount(c echo.Context) error {
	accountUUID := c.Get(middleware.ContextAccountUUID).(string)
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

func (env *environment) routeAccountAudit(c echo.Context) error {
	accountUUID := c.Get(middleware.ContextAccountUUID).(string)
	logrus.Infof("Get account audit for %s", accountUUID)

	account, err := env.db.FindAccount(accountUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.JsonErrorf("Logged in with unknown account"))
	}

	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	records, err := env.db.GetAuditTrailForAccount(account, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	ret := make([]common.Json, len(records))
	for i, record := range records {
		ret[i] = common.Json{
			"ts":      record.CreatedAt,
			"module":  record.Module,
			"level":   record.Level,
			"message": record.Message,
		}
	}

	return c.JSON(http.StatusOK, common.Json{
		"records": ret,
	})
}
