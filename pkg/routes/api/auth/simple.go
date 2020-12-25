package auth

import (
	"errors"
	"net/http"
	"simple-auth/pkg/config"
	"simple-auth/pkg/routes/common"
	"simple-auth/pkg/saerrors"
	"simple-auth/pkg/services"
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
	config     *config.ConfigSimpleAuthenticator
	localLogin services.LocalLoginService
}

func NewSimpleAuthController(localLoginService services.LocalLoginService, config *config.ConfigSimpleAuthenticator) *SimpleAuthController {
	return &SimpleAuthController{
		config,
		localLoginService,
	}
}

func (env *SimpleAuthController) Mount(group *echo.Group) {
	logrus.Info("Enabling simple auth...")
	group.POST("", env.routeSimpleAuthenticate)
}

type simpleAuthRequest struct {
	Username string  `json:"username" form:"username"`
	Password string  `json:"password" form:"password"`
	TOTP     *string `json:"totp" form:"totp"`
}

type simpleAuthResponse struct {
	ID string `json:"id"` // Account ID
}

// @Summary Authenticate
// @Description Authenticate username and password. 200 on success, otherwise 403
// @Tags Auth
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param simpleAuthRequest body simpleAuthRequest true "Credentials"
// @Success 200 {object} simpleAuthResponse
// @Failure 400,401,403,404,500 {object} common.ErrorResponse
// @Router /auth/simple [post]
func (env *SimpleAuthController) routeSimpleAuthenticate(c echo.Context) error {
	const metricName = "simple"

	if env.config.SharedSecret != "" {
		authHeader := c.Request().Header.Get(echo.HeaderAuthorization)
		if authHeader == "" {
			return common.HttpBadRequest(c, errors.New("authorization header required"))
		}
		authParts := strings.Split(authHeader, " ")
		if len(authParts) != 2 || authParts[0] != "Bearer" {
			return common.HttpBadRequest(c, errors.New("expected 'Bearer' on authHeader"))
		}
		if authParts[1] != env.config.SharedSecret {
			return common.HttpError(c, http.StatusUnauthorized, invalidAuthorizationToken.New())
		}
	}

	var req simpleAuthRequest
	if err := c.Bind(&req); err != nil {
		return common.HttpBadRequest(c, err)
	}

	authLocal, err := env.localLogin.WithContext(c).AssertLogin(req.Username, req.Password, req.TOTP)
	if err != nil {
		incAuthCounterError(metricName, err)
		return common.HttpError(c, http.StatusForbidden, err)
	}

	incAuthCounterSuccess(metricName)

	return c.JSON(http.StatusOK, simpleAuthResponse{
		ID: authLocal.Account().UUID,
	})
}
