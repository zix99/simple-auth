package main

import (
	"simple-auth/pkg/config"

	"github.com/labstack/echo/v4"
)

type pluginHook func(e *echo.Echo, cfg *config.Config)

var plugins []pluginHook

//lint:ignore U1000 It's possible to not have any activated tags to add a hook
func addHook(hook pluginHook) {
	plugins = append(plugins, hook)
}

func applyHooks(e *echo.Echo, cfg *config.Config) {
	for _, v := range plugins {
		v(e, cfg)
	}
}
