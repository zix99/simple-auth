package ui

import (
	"net/http"
	"simple-auth/pkg/routes/common"
	"simple-auth/pkg/routes/middleware/selector/auth"

	"github.com/labstack/echo/v4"
)

type changePasswordRequest struct {
	OldPassword string `json:"oldpassword"` // Not required if source is one-time (eg reset link)
	NewPassword string `json:"newpassword"`
}

func (env *environment) routeChangePasswordRequirements(c echo.Context) error {
	authContext := auth.MustGetAuthContext(c)

	return c.JSON(http.StatusOK, common.Json{
		"requireOldPassword": authContext.Source != auth.SourceOneTime,
	})
}

func (env *environment) routeChangePassword(c echo.Context) error {
	authContext := auth.MustGetAuthContext(c)

	var req changePasswordRequest
	if err := c.Bind(&req); err != nil {
		return common.HttpBadRequest(c, err)
	}

	authLocal, err := env.localLoginService.FindAuthLocal(authContext.UUID)
	if err != nil {
		return common.HttpInternalError(c, err)
	}

	if authContext.Source == auth.SourceOneTime {
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
