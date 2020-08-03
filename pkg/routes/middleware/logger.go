package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func NewLoggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			req := c.Request()
			resp := c.Response()
			logrus.Infof("%s %s %s %d %d", c.RealIP(), req.Method, req.RequestURI, resp.Status, resp.Size)
			return err
		}
	}
}
