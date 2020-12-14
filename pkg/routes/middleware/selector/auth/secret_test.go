package auth

import (
	"net/http"
	"net/http/httptest"
	"simple-auth/pkg/routes/middleware/selector"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var secreteAuthMiddleware = selector.NewSelectorMiddleware(
	NewSharedSecretWithAccountAuth("super-secret"),
	selector.HandlerReturns(http.StatusUnauthorized, "invalid"),
)

func TestSecretAuthProviderNoAuth(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	rec, _ := makeMiddlewareRequest(req, secreteAuthMiddleware)

	assert.Equal(t, 401, rec.Code)
}

func TestSecretAuthProviderSuccess(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("Authorization", "SharedKey super-secret")
	req.Header.Add("X-Account-UUID", "woot")

	rec, c := makeMiddlewareRequest(req, secreteAuthMiddleware)

	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, "ok", rec.Body.String())

	accountUUID := c.Get(contextAccountUUID).(string)
	assert.Equal(t, "woot", accountUUID)
}

func TestSecretAuthProviderBadKey(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("Authorization", "SharedKey fake")
	req.Header.Add("X-Account-UUID", "woot")

	rec, _ := makeMiddlewareRequest(req, secreteAuthMiddleware)

	assert.Equal(t, 401, rec.Code)
}

func TestSecretAuthProviderBadKeyType(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("Authorization", "asdf super-secret")
	req.Header.Add("X-Account-UUID", "woot")

	rec, _ := makeMiddlewareRequest(req, secreteAuthMiddleware)

	assert.Equal(t, 401, rec.Code)
}

func mockMiddleware(key, val string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(key, val)
			return next(c)
		}
	}
}

func TestSecretAuthProviderChainedMiddleware(t *testing.T) {
	ap := selector.NewSelectorMiddleware(
		NewSharedSecretWithAccountAuth("super-secret", mockMiddleware("A", "1"), mockMiddleware("A", "2"), mockMiddleware("B", "3")),
		NewSharedSecretWithAccountAuth("other-secret", mockMiddleware("N", "5")),
	)

	{
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Add("Authorization", "SharedKey super-secret")
		req.Header.Add("X-Account-UUID", "woot")

		rec, c := makeMiddlewareRequest(req, ap)

		assert.Equal(t, 200, rec.Code)
		assert.Equal(t, "ok", rec.Body.String())

		assert.Equal(t, "woot", c.Get(contextAccountUUID))
		assert.Equal(t, "2", c.Get("A"))
		assert.Equal(t, "3", c.Get("B"))
		assert.Equal(t, "3", c.Get("B"))
		assert.Nil(t, c.Get("N"))
	}

	{
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Add("Authorization", "SharedKey other-secret")
		req.Header.Add("X-Account-UUID", "huh")

		rec, c := makeMiddlewareRequest(req, ap)

		assert.Equal(t, 200, rec.Code)
		assert.Equal(t, "ok", rec.Body.String())

		assert.Equal(t, "huh", c.Get(contextAccountUUID))
		assert.Equal(t, "5", c.Get("N"))
		assert.Nil(t, c.Get("A"))
		assert.Nil(t, c.Get("B"))
	}
}
