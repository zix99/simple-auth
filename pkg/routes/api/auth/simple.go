package auth

import (
	"errors"
	"net/http"
	"simple-auth/pkg/config"
	"simple-auth/pkg/db"
	"simple-auth/pkg/routes/common"
	"simple-auth/pkg/saerrors"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

/*
Simple authenticator will simply provide an endpoint that will either return a 200 or a 401
depending on whether the username/password has been validated
*/

const (
	invalidAuthorizationToken saerrors.ErrorCode = "invalid-auth-token"
)

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
			return common.HttpBadRequest(c, errors.New("Authorization header required"))
		}
		authParts := strings.Split(authHeader, " ")
		if len(authParts) != 2 || authParts[0] != "Bearer" {
			return common.HttpBadRequest(c, errors.New("Expected 'Bearer' on authHeader"))
		}
		if authParts[1] != env.config.SharedSecret {
			return common.HttpError(c, http.StatusUnauthorized, invalidAuthorizationToken.New())
		}
	}

	req := struct {
		Username string `json:"username" form:"username"`
		Password string `json:"password" form:"password"`
	}{}
	if err := c.Bind(&req); err != nil {
		return common.HttpBadRequest(c, err)
	}

	account, err := env.db.AssertSimpleAuth(req.Username, req.Password, nil)
	if err != nil {
		return common.HttpError(c, http.StatusForbidden, err)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"id": account.UUID,
	})
}
