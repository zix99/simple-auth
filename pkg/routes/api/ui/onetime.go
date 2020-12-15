package ui

import (
	"errors"
	"html/template"
	"net/http"
	"simple-auth/pkg/appcontext"
	"simple-auth/pkg/email"
	"simple-auth/pkg/routes/common"
	"simple-auth/pkg/routes/middleware/selector/auth"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type oneTimePostRequest struct {
	Email string `json:"email" form:"email" binding:"required"`
}

func (env *environment) routeOneTimePost(c echo.Context) error {
	logger := appcontext.GetLogger(c)

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
	go email.NewFromConfig(env.email).WithContext(c).SendForgotPasswordEmail(req.Email, &email.ForgotPasswordData{
		EmailData: email.EmailData{
			Company: env.meta.Company,
			BaseURL: baseURL,
		},
		ResetDuration: env.config.Login.OneTime.TokenDuration,
		ResetLink:     template.HTML(baseURL + "/onetime?token=" + token),
	})

	return c.JSON(http.StatusOK, common.Json{"status": true})
}

func (env *environment) routeOneTimeAuth(c echo.Context) error {
	logger := appcontext.GetLogger(c)

	token := strings.TrimSpace(c.QueryParam("token"))
	if token == "" {
		return common.HttpBadRequest(c, errors.New("missing token"))
	}

	logger.Infof("Attemping to one-time signin for token %s...", token)

	account, err := env.db.AssertOneTimeToken(token)
	if err != nil {
		return common.HttpError(c, http.StatusUnauthorized, err)
	}

	err = auth.CreateSession(c, &env.config.Login.Cookie, account, auth.SourceOneTime)
	if err != nil {
		return common.HttpInternalError(c, err)
	}

	return c.Redirect(http.StatusTemporaryRedirect, "/")
}
