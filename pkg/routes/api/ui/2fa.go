package ui

import (
	"bytes"
	"errors"
	"net/http"
	"simple-auth/pkg/lib/totp"
	"simple-auth/pkg/lib/totp/otpimagery"
	"simple-auth/pkg/routes/common"
	"simple-auth/pkg/routes/middleware"

	"github.com/labstack/echo/v4"
)

type tfaSetupResponse struct {
	Secret string `json:"secret"`
}

func (env *environment) routeSetup2FA(c echo.Context) error {
	claims, ok := middleware.GetSessionClaims(c)
	if !ok {
		return common.HttpInternalErrorf(c, "No session")
	}

	config := env.config.Login.TwoFactor
	t, err := totp.NewTOTP(config.KeyLength, config.Issuer, claims.Subject)
	if err != nil {
		return common.HttpInternalError(c, err)
	}

	return c.JSON(http.StatusOK, tfaSetupResponse{
		Secret: t.Secret(),
	})
}

func (env *environment) route2FAQRCodeImage(c echo.Context) error {
	config := env.config.Login.TwoFactor
	claims, ok := middleware.GetSessionClaims(c)
	if !ok {
		return common.HttpInternalErrorf(c, "No session")
	}

	secret := c.QueryParam("secret")
	if secret == "" {
		return common.HttpBadRequest(c, errors.New("missing secret"))
	}

	t, err := totp.FromSecret(secret, config.Issuer, claims.Subject)
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

func (env *environment) routeConfirm2FA(c echo.Context) error {
	var req tfaActivateRequest
	if err := c.Bind(&req); err != nil {
		return common.HttpBadRequest(c, err)
	}

	log := middleware.GetLogger(c)
	accountUUID := middleware.MustGetSessionAccountUUID(c)

	authLocal, err := env.localLoginService.FindAuthLocal(accountUUID)
	if err != nil {
		return common.HttpInternalError(c, err)
	}

	log.Infof("Setting up TOTP for %s", accountUUID)
	if err := env.localLoginService.ActivateTOTP(authLocal, req.Secret, req.Code); err != nil {
		return common.HttpError(c, http.StatusForbidden, err)
	}

	return common.HttpOK(c)
}

func (env *environment) routeDeactivate2FA(c echo.Context) error {
	code := c.QueryParam("code")

	uuid := c.Get(middleware.ContextAccountUUID).(string)
	authLocal, err := env.localLoginService.FindAuthLocal(uuid)
	if err != nil {
		return common.HttpInternalError(c, err)
	}

	if err := env.localLoginService.DeactivateTOTP(authLocal, code); err != nil {
		return common.HttpError(c, http.StatusUnauthorized, err)
	}

	return common.HttpOK(c)
}
