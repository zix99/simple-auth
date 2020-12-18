// +build prometheus

package main

import (
	"simple-auth/pkg/config"

	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
)

func init() {
	addHook(func(e *echo.Echo, config *config.Config) {
		if config.Web.Prometheus {
			p := prometheus.NewPrometheus("sa", nil)
			p.Use(e)
		}
	})
}
