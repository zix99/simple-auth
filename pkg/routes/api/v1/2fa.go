package v1

import (
	"bytes"
	"errors"
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

func (env *Environment) RouteSetup2FA(c echo.Context) error {
	secret, err := env.twoFactorService.CreateSecret()
	if err != nil {
		return common.HttpInternalError(c, err)
	}

	return c.JSON(http.StatusOK, tfaSetupResponse{
		Secret: secret,
	})
}

func (env *Environment) Route2FAQRCodeImage(c echo.Context) error {
	authContext := auth.MustGetAuthContext(c)

	secret := c.QueryParam("secret")
	if secret == "" {
		return common.HttpBadRequest(c, errors.New("missing secret"))
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
	Secret string `json:"secret"`
	Code   string `json:"code"`
}

func (env *Environment) RouteConfirm2FA(c echo.Context) error {
	loginService := env.localLoginService.WithContext(c)

	var req tfaActivateRequest
	if err := c.Bind(&req); err != nil {
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

func (env *Environment) RouteDeactivate2FA(c echo.Context) error {
	loginService := env.localLoginService.WithContext(c)
	code := c.QueryParam("code")

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
