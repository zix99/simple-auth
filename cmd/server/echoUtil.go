package main

import (
	"fmt"
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

func redirectVue(vueRoute string) echo.HandlerFunc {
	return func(c echo.Context) error {
		url := c.Request().URL
		vueURL := fmt.Sprintf("/#/%s?%s", vueRoute, url.Query().Encode())
		return c.Redirect(http.StatusTemporaryRedirect, vueURL)
	}
}
