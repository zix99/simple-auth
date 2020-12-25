package auth

import (
	"errors"
	"fmt"
	"net/http"
	"simple-auth/pkg/config"
	"simple-auth/pkg/db"
	"simple-auth/pkg/instrumentation"
	"simple-auth/pkg/routes/middleware/selector"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

var sessionCounter instrumentation.Counter = instrumentation.NewCounter("sa_session_create", "Session creation counter", "source")

const (
	SourceOIDC    SessionSource = "oidc"
	SourceLogin   SessionSource = "login"
	SourceOneTime SessionSource = "onetime"
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
	return nil, fmt.Errorf("unable to parse key for %s", method)
}

func issueSessionJwt(config *config.ConfigJWT, account *db.Account, source SessionSource) (string, error) {
	if len(config.SigningKey) < 8 {
		logrus.Warn("No JWT secret set, or secret too short.  User not able to login")
		return "", errors.New("server needs secret")
	}

	signingMethod := jwt.GetSigningMethod(strings.ToUpper(config.SigningMethod))
	if signingMethod == nil {
		return "", fmt.Errorf("unknown signing method %s, check your config", config.SigningMethod)
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
		Name:     config.Name,
		Value:    signedToken,
		HttpOnly: config.HTTPOnly,
		Secure:   config.SecureOnly,
		Expires:  time.Now().Add(time.Duration(config.JWT.ExpiresMinutes) * time.Minute),
		Domain:   config.Domain,
		Path:     config.Path,
	}

	c.SetCookie(cookie)

	sessionCounter.Inc(source)

	return nil
}

func ClearSession(c echo.Context, config *config.ConfigLoginCookie) {
	c.SetCookie(&http.Cookie{
		Name:     config.Name,
		Value:    "",
		HttpOnly: config.HTTPOnly,
		Secure:   config.SecureOnly,
		Expires:  time.Now(),
		Domain:   config.Domain,
		Path:     config.Path,
	})
}

func ParseContextSession(config *config.ConfigLoginCookie, c echo.Context) (*SimpleAuthClaims, error) {
	cookie, err := c.Cookie(config.Name)
	if err != nil || cookie == nil {
		return nil, errors.New("auth cookie not set")
	}

	token, err := jwt.ParseWithClaims(cookie.Value, &SimpleAuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		return parseSigningKey(config.JWT.SigningMethod, config.JWT.SigningKey, true)
	})
	if err != nil {
		return nil, errors.New("unable to parse JWT")
	}

	if claims, ok := token.Claims.(*SimpleAuthClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("token rejected")
}

func sessionSelector(cookieName string) selector.MiddlewareSelector {
	return func(c echo.Context) error {
		cookie, err := c.Cookie(cookieName)
		if cookie == nil {
			return errors.New("no session cookie")
		}
		if err != nil {
			return err
		}
		return nil
	}
}

func NewSessionAuthHandler(config *config.ConfigLoginCookie) AuthHandler {
	_, parseErr := parseSigningKey(config.JWT.SigningMethod, config.JWT.SigningKey, false)
	if config.JWT.SigningKey == "" || parseErr != nil {
		logrus.Warn("No JWT secret specified, refusing to bind user management endpoints")
		return func(c echo.Context) (*AuthContext, error) {
			return nil, errors.New("server not configured for session api calls")
		}
	}

	return func(c echo.Context) (*AuthContext, error) {
		claims, err := ParseContextSession(config, c)
		if err != nil {
			return nil, err
		}
		return &AuthContext{
			UUID:   claims.Subject,
			Source: claims.Source,
		}, nil
	}
}

func NewSessionAuthProvider(config *config.ConfigLoginCookie, middleware ...echo.MiddlewareFunc) selector.SelectorGroup {
	return NewAuthSelectorGroup(
		sessionSelector(config.Name),
		NewSessionAuthHandler(config),
		middleware...,
	)
}
