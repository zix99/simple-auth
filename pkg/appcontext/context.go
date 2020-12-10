package appcontext

import "github.com/labstack/echo/v4"

type Context interface {
	Get(key string) interface{}
}

type RWContext interface {
	Context
	Set(key string, i interface{})
}

// With is a short-hand to set a key on each echo context
func With(key string, i interface{}) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(key, i)
			return next(c)
		}
	}
}

type ContextualFunc func(c echo.Context) interface{}

// WithContextual providers an injected value that is contextual to the echo context
func WithContextual(key string, provider ContextualFunc) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if val := provider(c); val != nil {
				c.Set(key, val)
			}
			return next(c)
		}
	}
}
