package auth

import (
	"simple-auth/pkg/config"
	"simple-auth/pkg/db"
	"simple-auth/pkg/routes/common"
	"simple-auth/pkg/routes/middleware/selector/auth"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type VouchAuthController struct {
	db           db.SADB
	config       *config.ConfigVouchAuthenticator
	cookieConfig *config.ConfigLoginCookie
}

func NewVouchAuthController(db db.SADB, config *config.ConfigVouchAuthenticator, cookieConfig *config.ConfigLoginCookie) *VouchAuthController {
	return &VouchAuthController{
		db:           db,
		config:       config,
		cookieConfig: cookieConfig,
	}
}

func (env *VouchAuthController) Mount(group *echo.Group) {
	logrus.Info("Enabling vouch auth...")
	loggedInMiddleware := auth.NewAuthMiddleware(
		auth.NewSessionAuthHandler(env.cookieConfig),
	)
	group.GET("", env.routeVouchAuth, loggedInMiddleware)
}

// @Summary Vouch
// @Description A vouch endpoint that checks if the user is logged in via cookie.  Intended to be used as `auth_request` in nginx
// @Tags Auth
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} common.OKResponse
// @Router /auth/vouch [get]
func (env *VouchAuthController) routeVouchAuth(c echo.Context) error {
	incAuthCounterSuccess("vouch")
	return common.HttpOK(c)
}
