package ui

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

func (env *environment) routeLogin(c echo.Context) error {
	logger := appcontext.GetLogger(c)
	req := loginRequest{}
	if err := c.Bind(&req); err != nil {
		return common.HttpBadRequest(c, err)
	}

	logger.Infof("Attempting login for '%s'...", req.Username)

	authLocal, err := env.localLoginService.WithContext(c).AssertLogin(req.Username, req.Password, req.Totp)
	if err != nil {
		logger.Infof("Login for user '%s' rejected: %v", req.Username, err)
		return common.HttpError(c, http.StatusUnauthorized, err)
	}
	logger.Infof("Login for user '%s' accepted", req.Username)

	err = auth.CreateSession(c, &env.config.Login.Cookie, authLocal.Account(), auth.SourceLogin)
	if err != nil {
		return common.HttpError(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, common.Json{
		"ok": true,
	})
}

func (env *environment) routeLogout(c echo.Context) error {
	auth.ClearSession(c, &env.config.Login.Cookie)
	return c.JSON(http.StatusOK, common.Json{
		"ok": true,
	})
}
