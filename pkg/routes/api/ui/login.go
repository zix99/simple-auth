package ui

import (
	"net/http"
	"simple-auth/pkg/routes/common"
	"simple-auth/pkg/routes/middleware"

	"github.com/labstack/echo/v4"
)

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (env *environment) routeLogin(c echo.Context) error {
	logger := middleware.GetLogger(c)
	req := loginRequest{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, common.JsonErrorf("Unable to deserialize request"))
	}

	logger.Infof("Attempting login for '%s'...", req.Username)

	account, err := env.db.AssertSimpleAuth(req.Username, req.Password, nil)
	if err != nil {
		logger.Infof("Login for user '%s' rejected", req.Username)
		return c.JSON(http.StatusUnauthorized, common.JsonError(err))
	}
	logger.Infof("Login for user '%s' accepted", req.Username)

	err = middleware.CreateSession(c, &env.config.Login.Cookie, account, middleware.SessionSourceLogin)
	if err != nil {
		logger.Error(err)
		return c.JSON(http.StatusInternalServerError, common.JsonErrorf("Something went wrong signing JWT"))
	}

	return c.JSON(http.StatusOK, common.Json{
		"ok": true,
	})
}

func (env *environment) routeLogout(c echo.Context) error {
	middleware.ClearSession(c, &env.config.Login.Cookie)
	return c.JSON(http.StatusOK, common.Json{
		"ok": true,
	})
}
