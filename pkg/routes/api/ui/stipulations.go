package ui

import (
	"net/http"
	"simple-auth/pkg/db"
	"simple-auth/pkg/routes/common"

	"github.com/labstack/echo/v4"
)

type tokenStipulationRequest struct {
	Token     string `json:"token"`
	AccountID string `json:"account"`
}

func (env *environment) routeTokenStipulation(c echo.Context) error {
	var req tokenStipulationRequest
	if err := c.Bind(&req); err != nil {
		return common.HttpBadRequest(c, err)
	}

	account, err := env.db.FindAccount(req.AccountID)
	if err != nil {
		return common.HttpInternalError(c, err)
	}

	if err := env.db.SatisfyStipulation(account, &db.TokenStipulation{Code: req.Token}); err != nil {
		return common.HttpError(c, http.StatusUnauthorized, err)
	}

	return common.HttpOK(c)
}
