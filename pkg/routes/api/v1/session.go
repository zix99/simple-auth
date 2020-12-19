package v1

import (
	"net/http"
	"simple-auth/pkg/appcontext"
	"simple-auth/pkg/routes/common"
	"simple-auth/pkg/routes/middleware/selector/auth"

	"github.com/labstack/echo/v4"
)

type loginRequest struct {
	Username string  `json:"username" binding:"required"`
	Password string  `json:"password" binding:"required"`
	Totp     *string `json:"totp"`
}

// @Summary Session Login
// @Description Login to a session with username and password, and set cookie
// @Tags Session
// @Accept json
// @Produce json
// @Param loginRequest body loginRequest true "Body"
// @Success 200 {object} common.OKResponse
// @Failure 400,401,500 {object} common.ErrorResponse
// @Router /auth/session [post]
func (env *Environment) RouteSessionLogin(c echo.Context) error {
	logger := appcontext.GetLogger(c)
	req := loginRequest{}
	if err := c.Bind(&req); err != nil {
		return common.HttpBadRequest(c, err)
	}

	if req.Username == "" || req.Password == "" {
		return common.HttpBadRequest(c, common.ErrMissingFields.Newf("missing email or password"))
	}

	logger.Infof("Attempting login for '%s'...", req.Username)

	authLocal, err := env.localLoginService.WithContext(c).AssertLogin(req.Username, req.Password, req.Totp)
	if err != nil {
		logger.Infof("Login for user '%s' rejected: %v", req.Username, err)
		return common.HttpError(c, http.StatusUnauthorized, err)
	}
	logger.Infof("Login for user '%s' accepted", req.Username)

	if err := env.sessionService.IssueSession(c, authLocal, auth.SourceLogin); err != nil {
		return common.HttpError(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, common.Json{
		"ok": true,
	})
}

// @Summary Session Logout
// @Description Logout session (clear cookie)
// @Tags Session
// @Accept json
// @Produce json
// @Success 200 {object} common.OKResponse
// @Failure 400,401,500 {object} common.ErrorResponse
// @Router /auth/session [delete]
func (env *Environment) RouteSessionLogout(c echo.Context) error {
	env.sessionService.ClearSession(c)
	return common.HttpOK(c)
}
