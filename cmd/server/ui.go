package main

import (
	"net/http"
	"simple-auth/pkg/config"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type uiController struct {
	config *config.ConfigWeb
}

func newUIController(config *config.ConfigWeb) *uiController {
	return &uiController{config}
}

func (s *uiController) Mount(group *echo.Group) {
	group.Use(middleware.CSRF())

	group.GET("/", func(c echo.Context) error {
		context := buildTemplateContext(c, s.config)
		return c.Render(http.StatusOK, "home", context)
	})
	group.GET("/create", func(c echo.Context) error {
		context := buildTemplateContext(c, s.config)
		return c.Render(http.StatusOK, "createAccount", context)
	})
	group.GET("/manage", func(c echo.Context) error {
		context := buildTemplateContext(c, s.config)
		return c.Render(http.StatusOK, "manageAccount", context)
	})
	group.GET("/login", func(c echo.Context) error {
		context := buildTemplateContext(c, s.config)
		return c.Render(http.StatusOK, "login", context)
	})
}

func buildTemplateContext(c echo.Context, web *config.ConfigWeb) map[string]interface{} {
	context := make(map[string]interface{})
	for k, v := range web.Metadata {
		context[k] = v
	}
	context["Requirements"] = web.Requirements
	context["RecaptchaV2"] = struct {
		Enabled bool
		SiteKey string
		Theme   string
	}{web.RecaptchaV2.Enabled, web.RecaptchaV2.SiteKey, web.RecaptchaV2.Theme}

	context["csrf"] = c.Get("csrf")

	return context
}
