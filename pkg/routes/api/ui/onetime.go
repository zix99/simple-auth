package ui

import (
	"errors"
	"net/http"
	"simple-auth/pkg/email"
	"simple-auth/pkg/routes/common"
	"simple-auth/pkg/routes/middleware"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type oneTimePostRequest struct {
	Email string `json:"email" form:"email" binding:"required"`
}

func (env *environment) routeOneTimePost(c echo.Context) error {
	logger := middleware.GetLogger(c)

	var req oneTimePostRequest
	if err := c.Bind(&req); err != nil {
		return common.HttpBadRequest(c, err)
	}

	logger.Infof("Issuing one-time token to email %s...", req.Email)

	account, err := env.db.FindAccountByEmail(req.Email)
	if err != nil {
		logger.Warn("No account found for password reset")
		return c.JSON(http.StatusOK, common.Json{"status": true}) // A mis-direct, to prevent scanning for tokens
	}

	duration, err := time.ParseDuration(env.config.Login.OneTime.TokenDuration)
	if err != nil {
		return common.HttpInternalErrorf(c, "Invalid token duration. Config error")
	}

	token, err := env.db.CreateAccountOneTimeToken(account, duration)
	if err != nil {
		return common.HttpError(c, http.StatusUnauthorized, err)
	}

	baseURL := env.config.GetBaseURL()
	err = email.New(logger).SendForgotPasswordEmail(env.email, req.Email, &email.ForgotPasswordData{
		EmailData: email.EmailData{
			Company: env.meta.Company,
			BaseURL: baseURL,
		},
		ResetDuration: env.config.Login.OneTime.TokenDuration,
		ResetLink:     baseURL + "/onetime?token=" + token,
	})
	if err != nil {
		return common.HttpInternalError(c, errorEmailSend.Compose(err))
	}

	return c.JSON(http.StatusOK, common.Json{"status": true})
}

func (env *environment) routeOneTimeAuth(c echo.Context) error {
	logger := middleware.GetLogger(c)

	token := strings.TrimSpace(c.QueryParam("token"))
	if token == "" {
		return common.HttpBadRequest(c, errors.New("Missing token"))
	}

	logger.Infof("Attemping to one-time signin for token %s...", token)

	account, err := env.db.AssertOneTimeToken(token)
	if err != nil {
		return common.HttpError(c, http.StatusUnauthorized, err)
	}

	err = middleware.CreateSession(c, &env.config.Login.Cookie, account, middleware.SessionSourceOneTime)
	if err != nil {
		return common.HttpInternalError(c, err)
	}

	return c.Redirect(http.StatusTemporaryRedirect, "/")
}
