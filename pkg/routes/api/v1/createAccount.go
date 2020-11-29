package v1

import (
	"errors"
	"net/http"
	"simple-auth/pkg/routes/common"
	"simple-auth/pkg/routes/middleware"
	"simple-auth/pkg/routes/middleware/selector/auth"
	"simple-auth/pkg/saerrors"

	"github.com/labstack/echo/v4"
)

type createAccountRequest struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Email       string `json:"email" binding:"required"`
	RecaptchaV2 string `json:"recaptchav2" binding:"required"`
}

type createAccountResponse struct {
	ID string `json:"id"`
}

const (
	usernameUnavailable saerrors.ErrorCode = "username-unavailable"
)

func (env *Environment) RouteCreateAccount(c echo.Context) error {
	logger := middleware.GetLogger(c)

	var req createAccountRequest
	if err := c.Bind(&req); err != nil {
		return common.HttpBadRequest(c, err)
	}

	if exists, err := env.localLoginService.UsernameExists(req.Username); exists || err != nil {
		return common.HttpError(c, http.StatusConflict, usernameUnavailable.Wrap(err))
	}

	account, err := env.accountService.CreateAccount(req.Username, req.Email)
	if err != nil {
		return common.HttpError(c, http.StatusBadRequest, err)
	}

	_, err = env.localLoginService.Create(account, req.Username, req.Password)
	if err != nil {
		return common.HttpError(c, http.StatusBadGateway, err)
	}

	if err := auth.CreateSession(c, env.loginConfig, account, auth.SourceLogin); err != nil {
		logger.Warnf("Unable to create session post-login, ignoring: %v", err)
	}

	return c.JSON(http.StatusCreated, &createAccountResponse{
		ID: account.UUID,
	})
}

type checkUsernameRequest struct {
	Username string `json:"username"`
}

type checkUsernameResponse struct {
	Username string `json:"username"`
	Exists   bool   `json:"exists"`
}

func (env *Environment) RouteCheckUsername(c echo.Context) error {
	var req checkUsernameRequest
	if err := c.Bind(&req); err != nil {
		return common.HttpBadRequest(c, err)
	}
	if req.Username == "" {
		return common.HttpBadRequest(c, errors.New("missing username"))
	}

	exists, _ := env.localLoginService.UsernameExists(req.Username)

	return c.JSON(http.StatusOK, &checkUsernameResponse{
		Username: req.Username,
		Exists:   exists,
	})
}
