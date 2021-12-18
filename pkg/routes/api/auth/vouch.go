package auth

import (
	"net/http"
	"net/url"
	"simple-auth/pkg/appcontext"
	"simple-auth/pkg/config"
	"simple-auth/pkg/db"
	"simple-auth/pkg/routes/common"
	"simple-auth/pkg/routes/middleware/selector"
	"simple-auth/pkg/routes/middleware/selector/auth"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type VouchAuthController struct {
	db           db.SADB
	config       *config.ConfigVouchAuthenticator
	cookieConfig *config.ConfigLoginCookie
	baseUrl      string
}

func NewVouchAuthController(db db.SADB, config *config.ConfigVouchAuthenticator, cookieConfig *config.ConfigLoginCookie, baseUrl string) *VouchAuthController {
	return &VouchAuthController{
		db:           db,
		config:       config,
		cookieConfig: cookieConfig,
		baseUrl:      baseUrl,
	}
}

func (env *VouchAuthController) Mount(group *echo.Group) {
	logrus.Info("Enabling vouch auth...")
	loggedInMiddleware := selector.NewSelectorMiddleware(
		auth.NewSessionAuthProvider(env.cookieConfig),
		env.authRedirectIfNeeded,
		selector.HandlerUnauthorized(),
	)
	group.GET("", env.routeVouchAuth, loggedInMiddleware)
}

// Resolves URL assuming forward headers (As defined in traefik)
func resolveContinueUrl(c echo.Context) string {
	if continueQuery := c.QueryParam("continue"); continueQuery != "" {
		return continueQuery
	}

	headers := c.Request().Header
	if host := headers.Get("X-Forwarded-Host"); host != "" {
		proto := headers.Get("X-Forwarded-Proto")
		if proto == "" {
			proto = "http"
		}
		return proto + "://" + host + headers.Get("X-Forwarded-Uri")
	}
	return ""
}

// If redirect is requested, send the user to authentication and try to return to the source
func (env *VouchAuthController) authRedirectIfNeeded(next echo.HandlerFunc, c echo.Context) (bool, error) {
	log := appcontext.GetLogger(c)

	if c.QueryParam("forward") != "" {
		continueUrl := resolveContinueUrl(c)

		if baseUrl, err := url.Parse(env.baseUrl); err == nil {
			if continueUrl != "" {
				qp := baseUrl.Query()
				qp.Set("continue", continueUrl)
				baseUrl.RawQuery = qp.Encode()
			}
			return true, c.Redirect(http.StatusTemporaryRedirect, baseUrl.String())
		} else {
			log.Errorf("Error parsing forward url during vouch: %v", err)
			return true, c.Redirect(http.StatusTemporaryRedirect, env.baseUrl)
		}
	}

	return false, nil
}

// @Summary Vouch
// @Description A vouch endpoint that checks if the user is logged in via cookie.  Intended to be used as `auth_request` in nginx for forwardauth in traefik
// @Tags Auth
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param forward query boolean false "If true, will forward to login with a 307 rather than return a 401"
// @Param continue query string false "Will override X-Forward headers to set the continue URL.  Must follow allowedContinueURL settings"
// @Success 200 {object} common.OKResponse
// @Failure 307,401 {object} common.ErrorResponse
// @Router /auth/vouch [get]
func (env *VouchAuthController) routeVouchAuth(c echo.Context) error {
	incAuthCounterSuccess("vouch")

	if env.config.UserHeader != "" {
		accountUUID := auth.MustGetAccountUUID(c)
		c.Response().Header().Set(env.config.UserHeader, accountUUID)
	}

	return common.HttpOK(c)
}
