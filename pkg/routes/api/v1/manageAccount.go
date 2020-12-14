package v1

import (
	"net/http"
	"simple-auth/pkg/appcontext"
	"simple-auth/pkg/routes/common"
	"simple-auth/pkg/routes/middleware/selector/auth"
	"simple-auth/pkg/saerrors"
	"strconv"

	"github.com/labstack/echo/v4"
)

const (
	errorInvalidAccount saerrors.ErrorCode = "invalid-account"
)

func (env *Environment) RouteGetAccount(c echo.Context) error {
	logger := appcontext.GetLogger(c)
	sadb := appcontext.GetSADB(c)
	accountUUID := auth.MustGetAccountUUID(c)

	logger.Infof("Get account for %s", accountUUID)
	account, err := sadb.FindAccount(accountUUID)
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

	if authLocal, err := env.localLoginService.WithContext(c).FindAuthLocal(account.UUID); err == nil {
		responseAuth["simple"] = common.Json{
			"username":         authLocal.Username(),
			"twofactor":        authLocal.HasTOTP(),
			"twofactorallowed": env.localLoginService.AllowTOTP(),
		}
	}

	if providers, err := sadb.FindOIDCForAccount(account); err == nil {
		oidcProviders := make([]common.Json, len(providers))
		for i, oidc := range providers {
			providerConfig := env.oidcService.GetProvider(oidc.Provider)
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

func (env *Environment) RouteGetAccountAudit(c echo.Context) error {
	accountUUID := auth.MustGetAccountUUID(c)
	logger := appcontext.GetLogger(c)
	sadb := appcontext.GetSADB(c)

	logger.Infof("Get account audit for %s", accountUUID)

	account, err := sadb.FindAccount(accountUUID)
	if err != nil {
		return common.HttpError(c, http.StatusInternalServerError, errorInvalidAccount.Wrapf(err, "Logged in with unknown account"))
	}

	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	records, err := sadb.GetAuditTrailForAccount(account, offset, limit)
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
