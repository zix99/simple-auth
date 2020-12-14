package auth

import (
	"fmt"
	"net/http"
	"simple-auth/pkg/routes/middleware/selector"

	"github.com/labstack/echo/v4"
)

type SessionSource string

type AuthContext struct {
	UUID   string
	Source SessionSource
}

type AuthHandler func(c echo.Context) (*AuthContext, error)

const (
	contextAccountUUID = "accountUUID"
	contextAuth        = "auth"
)

func (s *AuthContext) valid() bool {
	if s.UUID == "" {
		return false
	}
	if s.Source == "" {
		return false
	}
	return true
}

func setAuthContext(c echo.Context, context *AuthContext) {
	c.Set(contextAccountUUID, context.UUID)
	c.Set(contextAuth, context)
}

func NewAuthMiddleware(handler AuthHandler) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			context, err := handler(c)
			if context == nil || !context.valid() || err != nil {
				return c.JSON(http.StatusUnauthorized,
					jsonErrorf("unauthorized", fmt.Sprintf("Unable to authenticate: %s", err.Error())))
			}
			setAuthContext(c, context)
			return next(c)
		}
	}
}

func NewAuthSelectorGroup(mwSelector selector.MiddlewareSelector, handler AuthHandler, middleware ...echo.MiddlewareFunc) selector.SelectorGroup {
	authMiddleware := NewAuthMiddleware(handler)
	middlewareStack := append([]echo.MiddlewareFunc{authMiddleware}, middleware...)
	return selector.NewSelectorGroup(mwSelector, middlewareStack...)
}

func GetAccountUUID(c echo.Context) (string, bool) {
	ret, ok := c.Get(contextAccountUUID).(string)
	return ret, ok
}

func MustGetAccountUUID(c echo.Context) string {
	ret, ok := GetAccountUUID(c)
	if !ok {
		panic("Required auth UUID, bad middleware?")
	}
	return ret
}

func GetAuthContext(c echo.Context) (*AuthContext, bool) {
	ret, ok := c.Get(contextAuth).(*AuthContext)
	return ret, ok
}

func MustGetAuthContext(c echo.Context) *AuthContext {
	ret, ok := GetAuthContext(c)
	if !ok {
		panic("Required auth context, but was missing")
	}
	return ret
}
