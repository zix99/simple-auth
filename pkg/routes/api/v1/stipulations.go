package v1

import (
	"net/http"
	"simple-auth/pkg/appcontext"
	"simple-auth/pkg/db"
	"simple-auth/pkg/routes/common"
	"simple-auth/pkg/routes/middleware/selector/auth"

	"github.com/labstack/echo/v4"
)

type tokenStipulationRequest struct {
	Token     string `json:"token"`
	AccountID string `json:"account"`
}

func (env *Environment) RouteSatisfyTokenStipulation(c echo.Context) error {
	var req tokenStipulationRequest
	if err := c.Bind(&req); err != nil {
		return common.HttpBadRequest(c, err)
	}

	sadb := appcontext.GetSADB(c)

	if req.AccountID == "" { // Might be in session
		req.AccountID, _ = auth.GetAccountUUID(c)
	}

	account, err := sadb.FindAccount(req.AccountID)
	if err != nil {
		return common.HttpInternalError(c, err)
	}

	if err := sadb.SatisfyStipulation(account, &db.TokenStipulation{Code: req.Token}); err != nil {
		return common.HttpError(c, http.StatusUnauthorized, err)
	}

	return common.HttpOK(c)
}