package ui

import (
	"net/http"
	"simple-auth/pkg/routes/common"
	"simple-auth/pkg/routes/middleware"

	"github.com/labstack/echo/v4"
)

type changePasswordRequest struct {
	OldPassword string `json:"oldpassword"` // Not required if source is one-time (eg reset link)
	NewPassword string `json:"newpassword"`
}

func (env *environment) routeChangePasswordRequirements(c echo.Context) error {
	claims, ok := middleware.GetSessionClaims(c)
	if !ok {
		return common.HttpInternalErrorf(c, "Claims not present. No session?")
	}

	return c.JSON(http.StatusOK, common.Json{
		"requireOldPassword": claims.Source != middleware.SessionSourceOneTime,
	})
}

func (env *environment) routeChangePassword(c echo.Context) error {
	claims, ok := c.Get(middleware.ContextClaims).(*middleware.SimpleAuthClaims)
	if !ok {
		return common.HttpError(c, http.StatusUnauthorized, errorInvalidClaims.New())
	}

	var req changePasswordRequest
	if err := c.Bind(&req); err != nil {
		return common.HttpBadRequest(c, err)
	}

	accountUUID := c.Get(middleware.ContextAccountUUID).(string)
	account, err := env.db.FindAccount(accountUUID)
	if err != nil {
		return common.HttpError(c, http.StatusInternalServerError, errorInvalidAccount.Wrapf(err, "Unable to look up account"))

	}

	username, err := env.db.FindSimpleAuthUsername(account)
	if err != nil {
		return common.HttpError(c, http.StatusInternalServerError, errorInvalidAccount.Wrapf(err, "Username not associated with account"))
	}

	if claims.Source != middleware.SessionSourceOneTime {
		if _, err := env.db.AssertSimpleAuth(username, req.OldPassword, nil); err != nil {
			return common.HttpError(c, http.StatusUnauthorized, err)
		}
	}

	if err := env.db.UpdatePasswordForUsername(username, req.NewPassword); err != nil {
		return common.HttpError(c, http.StatusInternalServerError, errorInvalidAccount.Wrapf(err, "Unable to update password for user"))
	}

	return common.HttpOK(c)
}
