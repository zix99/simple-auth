package ui

import (
	"bytes"
	"errors"
	"net/http"
	"simple-auth/pkg/lib/totp"
	"simple-auth/pkg/lib/totp/otpimagery"
	"simple-auth/pkg/routes/common"
	"simple-auth/pkg/routes/middleware"
	"simple-auth/pkg/saerrors"

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

const (
	twoFactorInvalidCode saerrors.ErrorCode = "invalid-code"
)

func (env *environment) routeConfirm2FA(c echo.Context) error {
	var req tfaActivateRequest
	if err := c.Bind(&req); err != nil {
		return common.HttpBadRequest(c, err)
	}

	log := middleware.GetLogger(c)
	claims, ok := middleware.GetSessionClaims(c)
	if !ok {
		return common.HttpInternalErrorf(c, "No session")
	}

	account, err := env.db.FindAccount(claims.Subject)
	if err != nil {
		return common.HttpInternalError(c, err)
	}

	log.Infof("Setting up TOTP for %s", claims.Subject)

	t, err := totp.FromSecret(req.Secret, env.config.Login.TwoFactor.Issuer, claims.Subject)
	if err != nil {
		return common.HttpInternalError(c, err)
	}

	if !t.Validate(req.Code, env.config.Login.TwoFactor.Drift) {
		return common.HttpError(c, http.StatusForbidden, twoFactorInvalidCode.New())
	}

	tStr := t.String()
	if err := env.db.SetAuthSimpleTOTP(account, &tStr); err != nil {
		return common.HttpInternalError(c, err)
	}

	return common.HttpOK(c)
}

func (env *environment) routeDeactivate2FA(c echo.Context) error {
	code := c.QueryParam("code")

	uuid := c.Get(middleware.ContextAccountUUID).(string)
	account, err := env.db.FindAccount(uuid)
	if err != nil {
		return common.HttpInternalError(c, err)
	}

	if !env.db.ValidateTOTP(account, code) {
		return common.HttpError(c, http.StatusUnauthorized, twoFactorInvalidCode.New())
	}

	if err := env.db.SetAuthSimpleTOTP(account, nil); err != nil {
		return common.HttpInternalError(c, err)
	}

	return common.HttpOK(c)
}
