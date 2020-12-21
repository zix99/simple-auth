package v1

import (
	"net/http"
	"simple-auth/pkg/appcontext"
	"simple-auth/pkg/routes/common"
	"simple-auth/pkg/routes/middleware/selector/auth"
	"simple-auth/pkg/saerrors"
	"strconv"

	"github.com/labstack/echo/v4"
)

type createAccountRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required" format:"password"`
	Email    string `json:"email" validate:"required"`
}

type createAccountResponse struct {
	ID             string `json:"id"`                       // ID of the created user
	CreatedSession bool   `json:"createdSession,omitempty"` // Did create session
}

const (
	usernameUnavailable saerrors.ErrorCode = "username-unavailable"
)

// RouteCreateAccount creates a new account from echo context
// @Summary Create Account
// @Description Create a new account object
// @Tags Account
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param createRequest body createAccountRequest true "Create request"
// @Param createSession query boolean false "Attempts to create a session on successful login"
// @Success 200 {object} createAccountResponse
// @Failure 400,401,404,500 {object} common.ErrorResponse
// @Router /account [post]
func (env *Environment) RouteCreateAccount(c echo.Context) error {
	logger := appcontext.GetLogger(c)
	loginService := env.localLoginService.WithContext(c)
	accountService := env.accountService.WithContext(c)

	var req createAccountRequest
	if err := c.Bind(&req); err != nil {
		return common.HttpBadRequest(c, err)
	}
	if err := c.Validate(&req); err != nil {
		return common.HttpBadRequest(c, err)
	}

	createSession, _ := strconv.ParseBool(c.QueryParam("createSession"))

	if exists, err := loginService.UsernameExists(req.Username); exists || err != nil {
		return common.HttpError(c, http.StatusConflict, usernameUnavailable.Wrap(err))
	}

	account, err := accountService.CreateAccount(req.Username, req.Email)
	if err != nil {
		return common.HttpError(c, http.StatusBadRequest, err)
	}

	_, err = loginService.Create(account, req.Username, req.Password)
	if err != nil {
		return common.HttpError(c, http.StatusBadGateway, err)
	}

	ret := &createAccountResponse{
		ID: account.UUID,
	}

	if createSession && !accountService.HasUnsatisfiedStipulations(account) {
		if err := auth.CreateSession(c, env.loginConfig, account, auth.SourceLogin); err == nil {
			ret.CreatedSession = true
		} else {
			logger.Warnf("Unable to create session post-login, ignoring: %v", err)
		}
	}

	return c.JSON(http.StatusCreated, ret)
}

type checkUsernameRequest struct {
	Username string `json:"username" validate:"required"`
}

type checkUsernameResponse struct {
	Username string `json:"username"`
	Exists   bool   `json:"exists"`
}

// RouteCheckUsername checks if username is already in use
// @Summary Check username is in use
// @Description Check if username is already in use
// @Tags Account
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param checkUsernameRequest body checkUsernameRequest true "Request"
// @Success 200 {object} checkUsernameResponse
// @Failure 400,401,404,500 {object} common.ErrorResponse
// @Router /account/check [post]
func (env *Environment) RouteCheckUsername(c echo.Context) error {
	var req checkUsernameRequest
	if err := c.Bind(&req); err != nil {
		return common.HttpBadRequest(c, err)
	}
	if err := c.Validate(&req); err != nil {
		return common.HttpBadRequest(c, err)
	}

	exists, _ := env.localLoginService.WithContext(c).UsernameExists(req.Username)

	return c.JSON(http.StatusOK, &checkUsernameResponse{
		Username: req.Username,
		Exists:   exists,
	})
}
