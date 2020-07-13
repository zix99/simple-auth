package main

import (
	"net/http"

	"github.com/labstack/echo"
)

// Setup router on /oidc
func setupOidcRouter(group *echo.Group) {
	group.GET("", routeGetOidc)
}

func routeGetOidc(c echo.Context) error {
	return c.HTML(http.StatusOK, "Hi")
}
