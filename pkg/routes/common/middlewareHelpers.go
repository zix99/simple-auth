package common

import (
	"github.com/labstack/echo/v4"
)

// CoalesceMiddleware removes any nil middlewares
func CoalesceMiddleware(middlewares ...echo.MiddlewareFunc) (ret []echo.MiddlewareFunc) {
	for _, m := range middlewares {
		if m != nil {
			ret = append(ret, m)
		}
	}
	return
}
