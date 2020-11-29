package selector

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	MiddlewareSelector func(c echo.Context) error
	SelectorGroup      func(next echo.HandlerFunc, c echo.Context) (handled bool, err error)
)

func NewSelectorGroup(selector MiddlewareSelector, middleware ...echo.MiddlewareFunc) SelectorGroup {
	noNilMiddleware := coalesceMiddleware(middleware...)
	return func(next echo.HandlerFunc, c echo.Context) (bool, error) {
		if err := selector(c); err != nil {
			return false, fmt.Errorf("auth not handled: %w", err)
		}
		return true, chainMiddleware(next, noNilMiddleware...)(c)
	}
}

func chainMiddleware(next echo.HandlerFunc, middleware ...echo.MiddlewareFunc) echo.HandlerFunc {
	if len(middleware) > 0 {
		return middleware[0](func(c echo.Context) error {
			return chainMiddleware(next, middleware[1:]...)(c)
		})
	}
	return next
}

func coalesceMiddleware(middleware ...echo.MiddlewareFunc) []echo.MiddlewareFunc {
	ret := make([]echo.MiddlewareFunc, 0)
	for _, v := range middleware {
		if v != nil {
			ret = append(ret, v)
		}
	}
	return ret
}

func NewSelectorMiddleware(selectorGroups ...SelectorGroup) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			errors := make([]error, 0, len(selectorGroups))
			for _, selector := range selectorGroups {
				if ok, err := selector(next, c); ok {
					return err
				} else if err != nil {
					errors = append(errors, err)
				}
			}
			return c.JSON(http.StatusMethodNotAllowed, jsonErrorf("unhandled-request", errors, "The request could not be handled"))
		}
	}
}

func SelectorAlways(c echo.Context) error {
	return nil
}

func HandlerReturns(code int, i interface{}) SelectorGroup {
	return func(next echo.HandlerFunc, c echo.Context) (bool, error) {
		return true, c.JSON(code, i)
	}
}

func HandlerUnauthorized() SelectorGroup {
	return HandlerReturns(http.StatusUnauthorized, jsonErrorf("unauthorized", nil, "Request did not pass any authentication schemes"))
}
