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
	db        db.SADB
	config    *config.ConfigVouchAuthenticator
	jwtConfig *config.ConfigJWT
}

func NewVouchAuthController(db db.SADB, config *config.ConfigVouchAuthenticator, jwtConfig *config.ConfigJWT) *VouchAuthController {
	return &VouchAuthController{
		db:        db,
		config:    config,
		jwtConfig: jwtConfig,
	}
}

func (env *VouchAuthController) Mount(group *echo.Group) {
	logrus.Info("Enabling vouch auth...")
	loggedInMiddleware := auth.NewAuthMiddleware(
		auth.NewSessionAuthHandler(env.jwtConfig),
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
	return common.HttpOK(c)
}
