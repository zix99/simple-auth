package ui

import (
	"net/http"
	"simple-auth/pkg/routes/common"
	"simple-auth/pkg/routes/middleware"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (env *environment) routeAccount(c echo.Context) error {
	logger := middleware.GetLogger(c)
	accountUUID := c.Get(middleware.ContextAccountUUID).(string)

	logger.Infof("Get account for %s", accountUUID)
	account, err := env.db.FindAccount(accountUUID)
	if err != nil {
		return common.HttpError(c, http.StatusInternalServerError, errorInvalidAccount.Wrapf(err, "Logged in with unknown account"))
	}

	responseAuth := common.Json{}
	response := common.Json{
		"id":      account.UUID,
		"created": account.CreatedAt,
		"email":   account.Email,
		"auth":    responseAuth,
	}

	if authLocal, err := env.localLoginService.FindAuthLocal(account.UUID); err == nil {
		responseAuth["simple"] = common.Json{
			"username":         authLocal.Username(),
			"twofactor":        authLocal.HasTOTP(),
			"twofactorallowed": env.config.Login.TwoFactor.Enabled,
		}
	}

	if providers, err := env.db.FindOIDCForAccount(account); err == nil {
		oidcProviders := make([]common.Json, len(providers))
		for i, oidc := range providers {
			providerConfig := env.config.Login.OIDCByProvider(oidc.Provider)
			oidcProviders[i] = common.Json{
				"provider": oidc.Provider,
				"subject":  oidc.Subject,
				"icon":     providerConfig.Icon,
				"name":     providerConfig.Name,
			}
		}
		responseAuth["oidc"] = oidcProviders
	}

	return c.JSON(http.StatusOK, response)
}

func (env *environment) routeAccountAudit(c echo.Context) error {
	accountUUID := c.Get(middleware.ContextAccountUUID).(string)
	logger := middleware.GetLogger(c)

	logger.Infof("Get account audit for %s", accountUUID)

	account, err := env.db.FindAccount(accountUUID)
	if err != nil {
		return common.HttpError(c, http.StatusInternalServerError, errorInvalidAccount.Wrapf(err, "Logged in with unknown account"))
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
