package main

import (
	"net/http"
	"simple-auth/pkg/config"
	"simple-auth/pkg/routes/common"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type uiController struct {
	config *config.ConfigWeb
	meta   *config.ConfigMetadata
}

func newUIController(config *config.ConfigWeb, meta *config.ConfigMetadata) *uiController {
	return &uiController{config, meta}
}

func (s *uiController) Mount(group *echo.Group) {
	group.Use(middleware.CSRF())

	group.GET("/", func(c echo.Context) error {
		context := buildTemplateContext(c, s.meta, s.config)
		return c.Render(http.StatusOK, "home", context)
	})
}

func stringInList(val string, lst []string) bool {
	for _, item := range lst {
		if item == val {
			return true
		}
	}
	return false
}

func buildTemplateContext(c echo.Context, meta *config.ConfigMetadata, web *config.ConfigWeb) map[string]interface{} {
	continueURL := web.Login.RouteOnLogin
	queryContinue := c.QueryParam("continue")
	if queryContinue != "" && stringInList(queryContinue, web.Login.AllowedContinueUrls) {
		continueURL = queryContinue
	}

	app := common.Json{
		"company": meta.Company,
		"footer":  meta.Footer,
		"csrf":    c.Get("csrf"),
		"login": common.Json{
			"createAccount": web.Login.CreateAccountEnabled,
			"continue":      continueURL,
		},
		"requirements": web.Requirements,
		"recaptchav2": common.Json{
			"enabled": web.RecaptchaV2.Enabled,
			"sitekey": web.RecaptchaV2.SiteKey,
			"theme":   web.RecaptchaV2.Theme,
		},
	}

	for k, v := range meta.Bucket {
		app[k] = v
	}

	return app
}
