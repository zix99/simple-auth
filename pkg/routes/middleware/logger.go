package middleware

import (
	"simple-auth/pkg/appcontext"

	"github.com/labstack/echo/v4"
)

func NewRequestLoggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logger := appcontext.GetLogger(c)

			req := c.Request()
			logger.Debugf("START: %s %s %s", c.RealIP(), req.Method, req.RequestURI)

			err := next(c)

			resp := c.Response()

			respStatus := resp.Status
			respSize := resp.Size

			if httpErr, ok := err.(*echo.HTTPError); ok {
				respStatus = httpErr.Code
				respSize = int64(len(httpErr.Error()))
			}

			logger.Infof("%s %s %s %d %d", c.RealIP(), req.Method, req.RequestURI, respStatus, respSize)
			return err
		}
	}
}
