package v1

import (
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
// @Security SessionAuth
// @Accept json
// @Produce json
// @Success 200 {object} getLocalLoginResponse
// @Failure 400,401,404,500 {object} common.ErrorResponse
// @Router /local [get]
func (env *Environment) RouteGetLocalLogin(c echo.Context) error {
	resp, err := env.getLocalLoginResponse(c)
	if err != nil {
		return common.HttpError(c, http.StatusNotFound, localLoginNotFound.New())
	}

	return c.JSON(http.StatusOK, resp)
}

func (env *Environment) getLocalLoginResponse(c echo.Context) (*getLocalLoginResponse, error) {
	authContext := auth.MustGetAuthContext(c)

	authLocal, err := env.localLoginService.WithContext(c).FindAuthLocal(authContext.UUID)
	if err != nil {
		return nil, err
	}

	return &getLocalLoginResponse{
		Username:           authLocal.Username(),
		HasTwoFactor:       authLocal.HasTOTP(),
		AllowTwoFactor:     env.localLoginService.AllowTOTP(),
		RequireOldPassword: !allowUnsafePasswordUpdate(authContext),
	}, nil
}

type changePasswordRequest struct {
	OldPassword string `json:"oldpassword"` // Not required if source is one-time (eg reset link)
	NewPassword string `json:"newpassword" validate:"required"`
}

// RouteChangePassword change password for local auth
// @Summary Change Password
// @Tags Local
// @Description Change password for local auth
// @Security ApiKeyAuth
// @Security SessionAuth
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
	if err := c.Validate(&req); err != nil {
		return common.HttpBadRequest(c, err)
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
