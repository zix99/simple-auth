package ui

import (
	"net/http"
	"simple-auth/pkg/api/common"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (env *environment) routeLogin(c echo.Context) error {
	req := loginRequest{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, common.JsonErrorf("Unable to deserialize request"))
	}

	logrus.Infof("Attempting login for '%s'...", req.Username)

	account, err := env.db.FindAndVerifySimpleAuth(req.Username, req.Password)
	if err != nil {
		logrus.Infof("Login for user '%s' rejected", req.Username)
		return c.JSON(http.StatusUnauthorized, common.JsonError(err))
	}
	logrus.Infof("Login for user '%s' accepted", req.Username)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    env.config.JWT.Issuer,
		Subject:   account.UUID,
		Audience:  req.Username,
		ExpiresAt: time.Now().Add(time.Duration(env.config.JWT.ExpiresMinutes) * time.Minute).Unix(),
	})
	signedToken, err := token.SignedString([]byte(env.config.Secret))
	if err != nil {
		logrus.Warn(err)
		return c.JSON(http.StatusInternalServerError, common.JsonErrorf("Something went wrong signing JWT"))
	}

	c.SetCookie(&http.Cookie{
		Name:     "auth",
		Value:    signedToken,
		HttpOnly: true,
	})
	return c.JSON(http.StatusOK, common.Json{
		"ok": true,
	})
}
