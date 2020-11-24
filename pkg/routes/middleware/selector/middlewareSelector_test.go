package selector

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func makeMiddlewareRequest(req *http.Request, middleware echo.MiddlewareFunc) (*httptest.ResponseRecorder, echo.Context) {
	e := echo.New()
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	middleware(func(c echo.Context) error {
		return c.HTML(200, "ok")
	})(c)

	return rec, c
}

func TestSelectorWithNoGroups(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	mw := NewSelectorMiddleware()
	resp, _ := makeMiddlewareRequest(req, mw)
	assert.Equal(t, 405, resp.Code)
}

func TestSelectorWithPassThrough(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	mw := NewSelectorMiddleware(NewPassthruHandler())
	resp, _ := makeMiddlewareRequest(req, mw)
	assert.Equal(t, 200, resp.Code)
}

func mockSelectIfHasContext(c echo.Context) error {
	_, ok := c.Get("mock-key").(string)
	if ok {
		return nil
	}
	return errors.New("didnt pass")
}

func TestSelectorWithFailedSelector(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	mw := NewSelectorMiddleware(
		NewSelectorGroup(mockSelectIfHasContext),
		HandlerReturns(401, "Denied"),
	)
	resp, _ := makeMiddlewareRequest(req, mw)
	assert.Equal(t, 401, resp.Code)
}

func TestSelectorWithGoodSelector(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	mw := NewSelectorMiddleware(
		NewSelectorGroup(mockSelectIfHasContext),
		HandlerReturns(401, "Denied"),
	)

	e := echo.New()
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("mock-key", "heyo")

	mw(func(c echo.Context) error {
		return c.HTML(200, "ok")
	})(c)

	assert.Equal(t, 200, rec.Code)
}
