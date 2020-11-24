package v1

import (
	"errors"
	"net/http"
	"simple-auth/pkg/routes/common"

	"github.com/labstack/echo/v4"
)

type checkUsernameRequest struct {
	Username string `json:"username"`
}

type checkUsernameResponse struct {
	Username string `json:"username"`
	Exists   bool   `json:"exists"`
}

func (env *Environment) RouteCheckUsername(c echo.Context) error {
	var req checkUsernameRequest
	if err := c.Bind(&req); err != nil {
		return common.HttpBadRequest(c, err)
	}
	if req.Username == "" {
		return common.HttpBadRequest(c, errors.New("missing username"))
	}

	exists, _ := env.localLoginService.UsernameExists(req.Username)

	return c.JSON(http.StatusOK, &checkUsernameResponse{
		Username: req.Username,
		Exists:   exists,
	})
}
