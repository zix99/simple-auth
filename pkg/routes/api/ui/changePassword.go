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
		return common.HttpErrorf(c, http.StatusInternalServerError, "Claims not present. No session?")
	}

	return c.JSON(http.StatusOK, common.Json{
		"requireOldPassword": claims.Source != middleware.SessionSourceOneTime,
	})
}

func (env *environment) routeChangePassword(c echo.Context) error {
	logger := middleware.GetLogger(c)

	claims, ok := c.Get(middleware.ContextClaims).(*middleware.SimpleAuthClaims)
	if !ok {
		return common.HttpErrorf(c, http.StatusUnauthorized, "Invalid claims")
	}

	var req changePasswordRequest
	if err := c.Bind(&req); err != nil {
		return common.HttpError(c, http.StatusBadRequest, err)
	}

	accountUUID := c.Get(middleware.ContextAccountUUID).(string)
	account, err := env.db.FindAccount(accountUUID)
	if err != nil {
		logger.Error(err)
		return common.HttpErrorf(c, http.StatusInternalServerError, "Unable to look up account")
	}

	username, err := env.db.FindSimpleAuthUsername(account)
	if err != nil {
		logger.Error(err)
		return common.HttpErrorf(c, http.StatusInternalServerError, "Username not associated with account")
	}

	if claims.Source != middleware.SessionSourceOneTime {
		if _, err := env.db.FindAndVerifySimpleAuth(username, req.OldPassword); err != nil {
			return common.HttpErrorf(c, http.StatusUnauthorized, "Not allowed to update password")
		}
	}

	if err := env.db.UpdatePasswordForUsername(username, req.NewPassword); err != nil {
		logger.Error(err)
		return common.HttpErrorf(c, http.StatusInternalServerError, "Unable to update password for user")
	}

	return c.JSON(http.StatusOK, common.Json{
		"success": true,
	})
}
