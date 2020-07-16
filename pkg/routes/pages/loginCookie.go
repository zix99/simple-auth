package pages

import (
	"simple-auth/pkg/config"
	"simple-auth/pkg/db"

	"github.com/labstack/echo"
)

type LoginCookieController struct {
	db     db.SADB
	config *config.ConfigLoginCookie
}

func NewLoginCookieController(db db.SADB, config *config.ConfigLoginCookie) *LoginCookieController {
	return &LoginCookieController{db, config}
}

func (s *LoginCookieController) Mount(group *echo.Group) {
	group.GET("", s.routeLoginPage)
}

func (s *LoginCookieController) routeLoginPage(c echo.Context) error {
	return nil
}
