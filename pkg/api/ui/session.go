package ui

import (
	"errors"
	"net/http"
	"simple-auth/pkg/config"
	"simple-auth/pkg/db"
	"time"

	"github.com/labstack/echo"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

const authCookieName = "auth"

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

func createSession(c echo.Context, config *config.ConfigJWT, account *db.Account) error {
	signedToken, err := issueSessionJwt(config, account)
	if err != nil {
		logrus.Warn(err)
		return err
	}

	c.SetCookie(&http.Cookie{
		Name:     authCookieName,
		Value:    signedToken,
		HttpOnly: true,
		Expires:  time.Now().Add(time.Duration(config.ExpiresMinutes) * time.Minute),
	})

	return nil
}

func clearSession(c echo.Context) {
	c.SetCookie(&http.Cookie{
		Name:     authCookieName,
		Value:    "",
		HttpOnly: true,
		Expires:  time.Now(),
	})
}
