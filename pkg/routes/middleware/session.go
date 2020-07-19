package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"simple-auth/pkg/config"
	"simple-auth/pkg/db"
	"simple-auth/pkg/routes/common"
	"strings"
	"time"

	"github.com/labstack/echo"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

const authCookieName = "auth"

func parseSigningKey(method, key string) (interface{}, error) {
	lm := strings.ToUpper(method)
	if strings.HasPrefix(lm, "HS") {
		return []byte(key), nil
	}
	if strings.HasPrefix(lm, "RS") {
		return jwt.ParseRSAPrivateKeyFromPEM([]byte(key))
	}
	return nil, fmt.Errorf("Unable to parse key for %s", method)
}

func issueSessionJwt(config *config.ConfigJWT, account *db.Account) (string, error) {
	if len(config.SigningKey) < 8 {
		logrus.Warn("No JWT secret set, or secrete too short.  User not able to login")
		return "", errors.New("Server needs secret")
	}

	signingMethod := jwt.GetSigningMethod(strings.ToUpper(config.SigningMethod))
	if signingMethod == nil {
		return "", fmt.Errorf("Unknown signing method %s, check your config", config.SigningMethod)
	}

	decodedKey, err := parseSigningKey(config.SigningMethod, config.SigningKey)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(signingMethod, jwt.StandardClaims{
		Issuer:    config.Issuer,
		Subject:   account.UUID,
		Audience:  "simple-auth",
		ExpiresAt: time.Now().Add(time.Duration(config.ExpiresMinutes) * time.Minute).Unix(),
	})
	return token.SignedString(decodedKey)
}

func CreateSession(c echo.Context, config *config.ConfigLoginCookie, account *db.Account) error {
	signedToken, err := issueSessionJwt(&config.JWT, account)
	if err != nil {
		logrus.Warn(err)
		return err
	}

	cookie := &http.Cookie{
		Name:     authCookieName,
		Value:    signedToken,
		HttpOnly: config.HTTPOnly,
		Secure:   config.SecureOnly,
		Expires:  time.Now().Add(time.Duration(config.JWT.ExpiresMinutes) * time.Minute),
		Domain:   config.Domain,
		Path:     config.Path,
	}

	c.SetCookie(cookie)

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

func LoggedInMiddleware(config *config.ConfigJWT) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie(authCookieName)
			if err != nil || cookie == nil {
				return c.JSON(http.StatusUnauthorized, common.JsonErrorf("Cookie not set"))
			}

			token, err := jwt.ParseWithClaims(cookie.Value, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
				return parseSigningKey(config.SigningMethod, config.SigningKey)
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
