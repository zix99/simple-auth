package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"simple-auth/pkg/config"
	"simple-auth/pkg/db"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

const (
	authCookieName = "auth"
)

const (
	ContextClaims      = "auth"
	ContextAccountUUID = "accountUUID"
)

type SessionSource string

const (
	SessionSourceOIDC    SessionSource = "oidc"
	SessionSourceLogin                 = "login"
	SessionSourceOneTime               = "onetime"
)

type SimpleAuthClaims struct {
	jwt.StandardClaims
	Source SessionSource `json:"src,omitempty"`
}

func parseSigningKey(method, key string, verifying bool) (interface{}, error) {
	lm := strings.ToUpper(method)
	if strings.HasPrefix(lm, "HS") {
		return []byte(key), nil
	}
	if strings.HasPrefix(lm, "RS") {
		key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(key))
		if err != nil {
			return nil, err
		}
		if verifying {
			return &key.PublicKey, nil
		}
		return key, nil
	}
	return nil, fmt.Errorf("Unable to parse key for %s", method)
}

func issueSessionJwt(config *config.ConfigJWT, account *db.Account, source SessionSource) (string, error) {
	if len(config.SigningKey) < 8 {
		logrus.Warn("No JWT secret set, or secret too short.  User not able to login")
		return "", errors.New("Server needs secret")
	}

	signingMethod := jwt.GetSigningMethod(strings.ToUpper(config.SigningMethod))
	if signingMethod == nil {
		return "", fmt.Errorf("Unknown signing method %s, check your config", config.SigningMethod)
	}

	decodedKey, err := parseSigningKey(config.SigningMethod, config.SigningKey, false)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(signingMethod, SimpleAuthClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    config.Issuer,
			Subject:   account.UUID,
			Audience:  "simple-auth",
			ExpiresAt: time.Now().Add(time.Duration(config.ExpiresMinutes) * time.Minute).Unix(),
		},
		Source: source,
	})
	return token.SignedString(decodedKey)
}

func CreateSession(c echo.Context, config *config.ConfigLoginCookie, account *db.Account, source SessionSource) error {
	signedToken, err := issueSessionJwt(&config.JWT, account, source)
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

func ClearSession(c echo.Context, config *config.ConfigLoginCookie) {
	c.SetCookie(&http.Cookie{
		Name:     authCookieName,
		Value:    "",
		HttpOnly: config.HTTPOnly,
		Secure:   config.SecureOnly,
		Expires:  time.Now(),
		Domain:   config.Domain,
		Path:     config.Path,
	})
}

func parseContextSession(config *config.ConfigJWT, c echo.Context) (*SimpleAuthClaims, error) {
	cookie, err := c.Cookie(authCookieName)
	if err != nil || cookie == nil {
		return nil, errors.New("Auth cookie not set")
	}

	token, err := jwt.ParseWithClaims(cookie.Value, &SimpleAuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		return parseSigningKey(config.SigningMethod, config.SigningKey, true)
	})
	if err != nil {
		return nil, errors.New("Unable to parse JWT")
	}

	if claims, ok := token.Claims.(*SimpleAuthClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("Token rejected")
}

func LoggedInMiddleware(config *config.ConfigJWT) echo.MiddlewareFunc {
	_, parseErr := parseSigningKey(config.SigningMethod, config.SigningKey, false)
	if config.SigningKey == "" || parseErr != nil {
		logrus.Warn("No JWT secret specified, refusing to bind user management endpoints")
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				return c.JSON(http.StatusMethodNotAllowed, jsonErrorf("Server not configured to allow session API calls"))
			}
		}
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, err := parseContextSession(config, c)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, jsonErrorf(err.Error()))
			}
			c.Set(ContextClaims, claims)
			c.Set(ContextAccountUUID, claims.Subject)
			return next(c)
		}
	}
}

func GetSessionClaims(c echo.Context) (*SimpleAuthClaims, bool) {
	ret, ok := c.Get(ContextClaims).(*SimpleAuthClaims)
	return ret, ok
}

func jsonErrorf(s string, args ...interface{}) map[string]interface{} {
	return map[string]interface{}{
		"error":   true,
		"message": fmt.Sprintf(s, args...),
	}
}
