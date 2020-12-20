package v1

import (
	"bytes"
	"net/http"
	"simple-auth/pkg/appcontext"
	"simple-auth/pkg/lib/totp/otpimagery"
	"simple-auth/pkg/routes/common"
	"simple-auth/pkg/routes/middleware/selector/auth"

	"github.com/labstack/echo/v4"
)

type tfaSetupResponse struct {
	Secret string `json:"secret"`
}

// RouteSetup2FA gets parameters for new 2FA setup
// @Summary Get2FA Secret
// @Description Creates a new 2fa secret to use for next steps
// @Tags Local
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} tfaSetupResponse
// @Failure 400,401,404,500 {object} common.ErrorResponse
// @Router /local/2fa [get]
func (env *Environment) RouteSetup2FA(c echo.Context) error {
	secret, err := env.twoFactorService.CreateSecret()
	if err != nil {
		return common.HttpInternalError(c, err)
	}

	return c.JSON(http.StatusOK, tfaSetupResponse{
		Secret: secret,
	})
}

// Route2FAQRCodeImage gets qrcode to display to user
// @Summary Get2FA Secret
// @Description Gets qrcode to display to user
// @Tags Local
// @Security ApiKeyAuth
// @Accept json
// @Produce png
// @Param secret query string true "Secret to generate qrcode for"
// @Success 200
// @Failure 400,401,404,500 {object} common.ErrorResponse
// @Router /local/2fa/qrcode [get]
func (env *Environment) Route2FAQRCodeImage(c echo.Context) error {
	authContext := auth.MustGetAuthContext(c)

	secret := c.QueryParam("secret")
	if secret == "" {
		return common.HttpBadRequestf(c, "missing secret")
	}

	authLocal, err := env.localLoginService.WithContext(c).FindAuthLocal(authContext.UUID)
	if err != nil {
		return common.HttpInternalError(c, err)
	}

	t, err := env.twoFactorService.CreateFullSpecFromSecret(secret, authLocal)
	if err != nil {
		return common.HttpInternalError(c, err)
	}

	png, err := otpimagery.GenerateQRCode(t, 256)
	if err != nil {
		return common.HttpInternalError(c, err)
	}

	return c.Stream(http.StatusOK, "image/png", bytes.NewReader(png))
}

type tfaActivateRequest struct {
	Secret string `json:"secret" format:"password" validate:"required"`
	Code   string `json:"code" validate:"required"`
}

// RouteConfirm2FA confirm 2fa code and activate
// @Summary Setup 2FA
// @Description Activates 2FA
// @Tags Local
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param tfaActivateRequest body tfaActivateRequest true "Body"
// @Success 200 {object} common.OKResponse
// @Failure 400,401,404,500 {object} common.ErrorResponse
// @Router /local/2fa [post]
func (env *Environment) RouteConfirm2FA(c echo.Context) error {
	loginService := env.localLoginService.WithContext(c)

	var req tfaActivateRequest
	if err := c.Bind(&req); err != nil {
		return common.HttpBadRequest(c, err)
	}
	if err := c.Validate(&req); err != nil {
		return common.HttpBadRequest(c, err)
	}

	log := appcontext.GetLogger(c)
	accountUUID := auth.MustGetAccountUUID(c)

	authLocal, err := loginService.FindAuthLocal(accountUUID)
	if err != nil {
		return common.HttpInternalError(c, err)
	}

	t, err := env.twoFactorService.CreateFullSpecFromSecret(req.Secret, authLocal)
	if err != nil {
		return common.HttpInternalError(c, err)
	}

	log.Infof("Setting up TOTP for %s", accountUUID)
	if err := loginService.ActivateTOTP(authLocal, t, req.Code); err != nil {
		return common.HttpError(c, http.StatusForbidden, err)
	}

	return common.HttpOK(c)
}

// RouteDeactivate2FA deactivates 2fa
// @Summary Deactivate 2FA
// @Description Deactivate 2FA
// @Tags Local
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param code query string true "Code to check against before deactivating"
// @Success 200 {object} common.OKResponse
// @Failure 400,401,404,500 {object} common.ErrorResponse
// @Router /local/2fa [delete]
func (env *Environment) RouteDeactivate2FA(c echo.Context) error {
	loginService := env.localLoginService.WithContext(c)
	code := c.QueryParam("code")
	if code == "" {
		return common.HttpBadRequestf(c, "missing field: code")
	}

	uuid := auth.MustGetAccountUUID(c)
	authLocal, err := loginService.FindAuthLocal(uuid)
	if err != nil {
		return common.HttpInternalError(c, err)
	}

	if err := loginService.DeactivateTOTP(authLocal, code); err != nil {
		return common.HttpError(c, http.StatusUnauthorized, err)
	}

	return common.HttpOK(c)
}
