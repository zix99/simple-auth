package ui

import (
	"errors"
	"net/http"
	"simple-auth/pkg/api/common"
	"simple-auth/pkg/config"
	"simple-auth/pkg/db"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func issueSessionJwt(config *config.ConfigJWT, account *db.Account) (string, error) {
	if len(config.Secret) < 8 {
		logrus.Warn("No JWT secret set, or secrete too short.  User not able to login")
		return "", errors.New("Server needs secret")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    config.Issuer,
		Subject:   account.UUID,
		Audience:  "simple-auth",
		ExpiresAt: time.Now().Add(time.Duration(config.ExpiresMinutes) * time.Minute).Unix(),
	})
	return token.SignedString([]byte(config.Secret))
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

	signedToken, err := issueSessionJwt(&env.config.JWT, account)
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
