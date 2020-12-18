package v1

import (
	"errors"
	"net/http"
	"simple-auth/pkg/routes/common"
	"simple-auth/pkg/routes/middleware/selector/auth"
	"simple-auth/pkg/saerrors"

	"github.com/labstack/echo/v4"
)

type getLocalLoginResponse struct {
	Username           string `json:"username"`
	HasTwoFactor       bool   `json:"twofactor"`
	AllowTwoFactor     bool   `json:"twofactorallowed"`
	RequireOldPassword bool   `json:"requireOldPassword"`
}

const localLoginNotFound saerrors.ErrorCode = "local-login-not-found"

// RouteGetLocalLogin gets local auth setup, if any
// @Summary Get Local Auth
// @Tags Local
// @Description Get details about local authentication
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} getLocalLoginResponse
// @Failure 400,401,404,500 {object} common.ErrorResponse
// @Router /local [get]
func (env *Environment) RouteGetLocalLogin(c echo.Context) error {
	authContext := auth.MustGetAuthContext(c)

	if authLocal, err := env.localLoginService.WithContext(c).FindAuthLocal(authContext.UUID); err == nil {
		return c.JSON(http.StatusOK, &getLocalLoginResponse{
			Username:           authLocal.Username(),
			HasTwoFactor:       authLocal.HasTOTP(),
			AllowTwoFactor:     env.localLoginService.AllowTOTP(),
			RequireOldPassword: !allowUnsafePasswordUpdate(authContext),
		})
	}

	return common.HttpError(c, http.StatusNotFound, localLoginNotFound.New())
}

type changePasswordRequest struct {
	OldPassword string `json:"oldpassword"` // Not required if source is one-time (eg reset link)
	NewPassword string `json:"newpassword"`
}

// RouteChangePassword change password for local auth
// @Summary Change Password
// @Tags Local
// @Description Change password for local auth
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param changePasswordRequest body changePasswordRequest true "Change password request"
// @Success 200 {object} getLocalLoginResponse
// @Failure 400,401,404,500 {object} common.ErrorResponse
// @Router /local/password [post]
func (env *Environment) RouteChangePassword(c echo.Context) error {
	authContext := auth.MustGetAuthContext(c)
	loginService := env.localLoginService.WithContext(c)

	var req changePasswordRequest
	if err := c.Bind(&req); err != nil {
		return common.HttpBadRequest(c, err)
	}

	if req.NewPassword == "" {
		return common.HttpBadRequest(c, errors.New("missing newpassword"))
	}

	authLocal, err := loginService.FindAuthLocal(authContext.UUID)
	if err != nil {
		return common.HttpInternalError(c, err)
	}

	if allowUnsafePasswordUpdate(authContext) {
		// Change password, but exempt from the oldPassword requirement
		if err := loginService.UpdatePasswordUnsafe(authLocal, req.NewPassword); err != nil {
			return common.HttpInternalError(c, err)
		}
	} else {
		if err := loginService.UpdatePassword(authLocal, req.OldPassword, req.NewPassword); err != nil {
			return common.HttpError(c, http.StatusUnauthorized, err)
		}
	}

	return common.HttpOK(c)
}

func allowUnsafePasswordUpdate(ctx *auth.AuthContext) bool {
	return ctx.Source == auth.SourceOneTime || ctx.Source == auth.SourceSecret
}
