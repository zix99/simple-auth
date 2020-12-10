package appcontext

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

const loggerContextKey = "appcontext.logger"

func WithLogger(fl logrus.FieldLogger) echo.MiddlewareFunc {
	return With(loggerContextKey, fl)
}

// GetLogger will always return a logger
func GetLogger(c Context) logrus.FieldLogger {
	if logger, ok := c.Get(loggerContextKey).(logrus.FieldLogger); ok {
		return logger
	}
	return logrus.StandardLogger()
}

func SetLogger(c RWContext, lf logrus.FieldLogger) {
	c.Set(loggerContextKey, lf)
}
