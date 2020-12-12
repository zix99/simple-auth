package appcontext

import "github.com/labstack/echo/v4"

func (s ProviderFunc) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if key, val := s(c); val != nil {
				c.Set(key, val)
			}
			return next(c)
		}
	}
}
