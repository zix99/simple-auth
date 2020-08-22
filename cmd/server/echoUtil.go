package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func redirectHandler(to string) echo.HandlerFunc {
	return func(c echo.Context) error {
		newURL := *c.Request().URL
		newURL.Path = to
		return c.Redirect(http.StatusTemporaryRedirect, newURL.String())
	}
}
