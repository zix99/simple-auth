package middleware

import (
	"simple-auth/pkg/appcontext"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/random"
)

const headerCorrelationID = "X-Correlation-ID"
const correlationIDLength = 12

const contextCorrelationIDKey = "correlationID"

func NewCorrelationMiddleware(readHeader, writeHeader bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var cid string
			if readHeader {
				cid = c.Request().Header.Get(headerCorrelationID)
			}
			if cid == "" {
				cid = random.String(correlationIDLength)
			}

			c.Set(contextCorrelationIDKey, cid)

			if logger := appcontext.GetLogger(c); logger != nil {
				appcontext.SetLogger(c, logger.WithField("cid", cid))
			}

			if writeHeader {
				c.Response().Header().Set(headerCorrelationID, cid)
			}

			return next(c)
		}
	}
}

func GetCorrelationId(c appcontext.Context) (string, bool) {
	id, ok := c.Get(contextCorrelationIDKey).(string)
	return id, ok
}
