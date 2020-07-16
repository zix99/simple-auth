package middleware

import (
	"errors"
	"net/http"
	"simple-auth/pkg/config"
	"simple-auth/pkg/db"
	"simple-auth/pkg/routes/common"
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

func CreateSession(c echo.Context, config *config.ConfigJWT, account *db.Account) error {
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

func ClearSession(c echo.Context) {
	c.SetCookie(&http.Cookie{
		Name:     authCookieName,
		Value:    "",
		HttpOnly: true,
		Expires:  time.Now(),
	})
}

func LoggedInMiddleware(key string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie(authCookieName)
			if err != nil || cookie == nil {
				return c.JSON(http.StatusUnauthorized, common.JsonErrorf("Cookie not set"))
			}

			token, err := jwt.ParseWithClaims(cookie.Value, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(key), nil
			})
			if err != nil {
				return c.JSON(http.StatusUnauthorized, common.JsonErrorf("Unable to parse JWT"))
			}

			if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
				c.Set("auth", claims)
				c.Set("accountUUID", claims.Subject)
				return next(c)
			}

			return c.JSON(http.StatusUnauthorized, common.JsonErrorf("Token rejected"))
		}
	}
}
