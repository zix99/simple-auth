//+build swagger

package main

import (
	"simple-auth/pkg/config"

	_ "simple-auth/pkg/swagdocs"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func init() {
	addHook(func(e *echo.Echo, config *config.Config) {
		if config.Web.Swagger {
			e.GET("/swagger", redirectHandler("/swagger/index.html"))
			e.GET("/swagger/*", echoSwagger.WrapHandler)
		}
	})
}
