package auth

import (
	"net/http"
	"net/http/httptest"
	"simple-auth/pkg/routes/middleware/selector"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// helper to execute request with a given middleware
func makeMiddlewareRequest(req *http.Request, middleware echo.MiddlewareFunc) (*httptest.ResponseRecorder, echo.Context) {
	e := echo.New()
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	middleware(func(c echo.Context) error {
		return c.HTML(200, "ok")
	})(c)

	return rec, c
}

func TestNoAuthProvider(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	ap := selector.NewSelectorMiddleware()
	ok := false
	ap(func(c echo.Context) error {
		ok = true
		return nil
	})(c)

	assert.False(t, ok)
}

func TestSimpleAuthProvider(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	ap := selector.NewSelectorMiddleware(NewAuthSelectorGroup(selector.SelectorAlways, func(c echo.Context) (*AuthContext, error) {
		return &AuthContext{
			UUID:   "bla",
			Source: SourceLogin,
		}, nil
	}))

	ok := false
	var withContext *AuthContext
	ap(func(c echo.Context) error {
		ok = true
		withContext = MustGetAuthContext(c)
		return nil
	})(c)

	assert.True(t, ok)
	assert.NotNil(t, withContext)
	assert.Equal(t, "bla", withContext.UUID)
}
