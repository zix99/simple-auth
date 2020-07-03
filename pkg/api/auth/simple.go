package auth

import (
	"simple-auth/pkg/api/common"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

/*

Simple authenticator will simply provide an endpoint that will either return a 200 or a 401
depending on whether the username/password has been validated
*/

func setupSimpleAuthenticator(env *environment, group *echo.Group) {
	logrus.Info("Enabling simple auth...")
	group.POST("", env.routeSimpleAuthenticate)
}

func (env *environment) routeSimpleAuthenticate(c echo.Context) error {
	req := struct {
		Username string `json:"username" form:"username"`
		Password string `json:"password" form:"password"`
	}{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, common.JsonError(err))
	}

	account, err := env.db.FindAndVerifySimpleAuth(req.Username, req.Password)
	if err != nil {
		return c.JSON(401, common.JsonError(err))
	}

	return c.JSON(200, map[string]string{
		"id": account.UUID,
	})
}
