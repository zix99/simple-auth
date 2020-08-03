package main

import (
	"net/http"
	"simple-auth/pkg/config"
	"simple-auth/pkg/routes/common"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

func buildTemplateContext(c echo.Context, meta *config.ConfigMetadata, web *config.ConfigWeb) map[string]interface{} {
	continueURL := web.Login.Settings.ResolveContinueURL(c.QueryParam("continue"))

	app := common.Json{
		"company": meta.Company,
		"footer":  meta.Footer,
		"csrf":    c.Get("csrf"),
		"login": common.Json{
			"createAccount": web.Login.Settings.CreateAccountEnabled,
			"continue":      continueURL,
		},
		"requirements": web.Requirements,
		"recaptchav2": common.Json{
			"enabled": web.RecaptchaV2.Enabled,
			"sitekey": web.RecaptchaV2.SiteKey,
			"theme":   web.RecaptchaV2.Theme,
		},
		"oidc": buildOIDCContext(web.Login.OIDC),
	}

	for k, v := range meta.Bucket {
		app[k] = v
	}

	return app
}

func buildOIDCContext(oidc []*config.ConfigOIDCProvider) []common.Json {
	ret := []common.Json{}
	for _, config := range oidc {
		oidcConfig := common.Json{
			"id":   config.ID,
			"name": config.Name,
			"icon": config.Icon,
		}
		ret = append(ret, oidcConfig)
	}
	return ret
}
