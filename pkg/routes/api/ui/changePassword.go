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
	claims := middleware.MustGetSessionClaims(c)

	var req changePasswordRequest
	if err := c.Bind(&req); err != nil {
		return common.HttpBadRequest(c, err)
	}

	authLocal, err := env.localLoginService.FindAuthLocal(claims.Subject)
	if err != nil {
		return common.HttpInternalError(c, err)
	}

	if claims.Source == middleware.SessionSourceOneTime {
		// Change password, but exempt from the oldPassword requirement
		if err := env.localLoginService.UpdatePasswordUnsafe(authLocal, req.NewPassword); err != nil {
			return common.HttpInternalError(c, err)
		}
	} else {
		if err := env.localLoginService.UpdatePassword(authLocal, req.OldPassword, req.NewPassword); err != nil {
			return common.HttpError(c, http.StatusUnauthorized, err)
		}
	}

	return common.HttpOK(c)
}
