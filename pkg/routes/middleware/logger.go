package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

const contextLoggerKey = "logger"

func NewLoggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var logger logrus.FieldLogger
			var ok bool
			if logger, ok = GetLoggerOk(c); !ok {
				if correlationID, ok := c.Get(ContextCorrelationIDKey).(string); ok {
					logger = logger.WithField("cid", correlationID)
				}
				c.Set(contextLoggerKey, logger)
			}

			err := next(c)
			req := c.Request()
			resp := c.Response()
			logger.Infof("%s %s %s %d %d", c.RealIP(), req.Method, req.RequestURI, resp.Status, resp.Size)
			return err
		}
	}
}

func GetLoggerOk(c echo.Context) (logrus.FieldLogger, bool) {
	if logger, ok := c.Get(contextLoggerKey).(logrus.FieldLogger); ok {
		return logger, true
	}
	return logrus.StandardLogger(), false
}

func GetLogger(c echo.Context) (logger logrus.FieldLogger) {
	logger, _ = GetLoggerOk(c)
	return
}
