package auth

import (
	"errors"
	"simple-auth/pkg/routes/middleware/selector"
	"strings"

	"github.com/labstack/echo/v4"
)

const (
	SourceSecret SessionSource = "shared-secret"
)

const (
	SharedSecretHeader          = "Authorization"
	SharedSecretHeaderQualifier = "SharedKey"
	SharedSecretHeaderAccountId = "X-Account-UUID"
)

func checkSecret(secret string, c echo.Context) error {
	authHeader := c.Request().Header.Get(SharedSecretHeader)
	if authHeader == "" {
		return errors.New("missing authorization header")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 {
		return errors.New("missing authorization header qualifier")
	}

	if !strings.EqualFold(parts[0], SharedSecretHeaderQualifier) {
		return errors.New("expected sharedkey qualifier")
	}

	if parts[1] != secret {
		return errors.New("invalid secret")
	}

	return nil
}

func SharedSecretSelector(secret string) selector.MiddlewareSelector {
	return func(c echo.Context) error {
		return checkSecret(secret, c)
	}
}

func SharedSecretWithAccountMiddleware(secret string) AuthHandler {
	return func(c echo.Context) (*AuthContext, error) {
		if err := checkSecret(secret, c); err != nil {
			return nil, err
		}

		accountID := c.Request().Header.Get(SharedSecretHeaderAccountId)
		if accountID == "" {
			return nil, errors.New("missing x-account-uuid header")
		}

		return &AuthContext{
			UUID:   accountID,
			Source: SourceSecret,
		}, nil
	}
}

func NewSharedSecretWithAccountAuth(secret string, middleware ...echo.MiddlewareFunc) selector.SelectorGroup {
	return NewAuthSelectorGroup(
		SharedSecretSelector(secret),
		SharedSecretWithAccountMiddleware(secret),
		middleware...,
	)
}
