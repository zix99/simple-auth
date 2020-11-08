package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCorrelationIDMiddleware(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	h := NewCorrelationMiddleware(false, true)(routePassthru)

	if assert.NoError(t, h(c)) {
		assert.Equal(t, 200, rec.Code)
		assert.NotEmpty(t, rec.Header().Get("x-correlation-id"))
	}
}

func routePassthru(c echo.Context) error {
	return c.HTML(200, "ok")
}
