package auth

import (
	"net/http"
	"simple-auth/pkg/config"
	"simple-auth/pkg/db"
	"simple-auth/pkg/routes/common"
	"strings"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

/*
Simple authenticator will simply provide an endpoint that will either return a 200 or a 401
depending on whether the username/password has been validated
*/

type SimpleAuthController struct {
	db     db.SADB
	config *config.ConfigSimpleAuthenticator
}

func NewSimpleAuthController(db db.SADB, config *config.ConfigSimpleAuthenticator) *SimpleAuthController {
	return &SimpleAuthController{
		db,
		config,
	}
}

func (env *SimpleAuthController) Mount(group *echo.Group) {
	logrus.Info("Enabling simple auth...")
	group.POST("", env.routeSimpleAuthenticate)
}

func (env *SimpleAuthController) routeSimpleAuthenticate(c echo.Context) error {
	if env.config.SharedSecret != "" {
		authHeader := c.Request().Header.Get(echo.HeaderAuthorization)
		if authHeader == "" {
			return c.JSON(http.StatusBadRequest, common.JsonErrorf("Authorization header required"))
		}
		authParts := strings.Split(authHeader, " ")
		if len(authParts) != 2 || authParts[0] != "Bearer" {
			return c.JSON(http.StatusBadRequest, common.JsonErrorf("Expected 'Bearer' on authHeader"))
		}
		if authParts[1] != env.config.SharedSecret {
			return c.JSON(http.StatusUnauthorized, common.JsonErrorf("Invalid authorization token"))
		}
	}

	req := struct {
		Username string `json:"username" form:"username"`
		Password string `json:"password" form:"password"`
	}{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, common.JsonError(err))
	}

	account, err := env.db.FindAndVerifySimpleAuth(req.Username, req.Password)
	if err != nil {
		return c.JSON(http.StatusForbidden, common.JsonError(err))
	}

	return c.JSON(http.StatusOK, map[string]string{
		"id": account.UUID,
	})
}
