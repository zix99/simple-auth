package v1

import (
	"errors"
	"net/http"
	"simple-auth/pkg/appcontext"
	"simple-auth/pkg/routes/common"
	"strings"

	"github.com/labstack/echo/v4"
)

type oneTimePostRequest struct {
	Email string `json:"email" form:"email" binding:"required" example:"sa@example.com"`
}

// @Summary Create OneTime
// @Description Creates a new onetime token for an email
// @Tags Session
// @Accept json
// @Produce json
// @Param oneTimePostRequest body oneTimePostRequest true "Body"
// @Success 200 {object} common.OKResponse
// @Failure 400,401,500 {object} common.ErrorResponse
// @Router /auth/onetime [post]
func (env *Environment) RouteOneTimeCreateToken(c echo.Context) error {
	logger := appcontext.GetLogger(c)

	var req oneTimePostRequest
	if err := c.Bind(&req); err != nil {
		return common.HttpBadRequest(c, err)
	}

	if req.Email == "" {
		return common.HttpBadRequest(c, common.ErrMissingFields.Newf("Missing email"))
	}

	logger.Infof("Issuing one-time token to email %s...", req.Email)

	account, err := env.accountService.WithContext(c).FindAccountByEmail(req.Email)
	if err != nil {
		logger.Warn("No account found for password reset")
		return common.HttpOK(c) // A mis-direct, to prevent scanning for emails
	}

	if err := env.sessionService.WithContext(c).IssueOneTimeToken(account); err != nil {
		return common.HttpInternalError(c, err)
	}

	return common.HttpOK(c)
}

// @Summary OneTime Authenticate
// @Description Loggin via onetime token and create session
// @Tags Session
// @Accept json
// @Param token query string true "OneTime token to authenticate against"
// @Success 302 {object} common.OKResponse
// @Failure 400,401,500 {object} common.ErrorResponse
// @Router /auth/onetime [get]
func (env *Environment) RouteOneTimeAuth(c echo.Context) error {
	logger := appcontext.GetLogger(c)

	token := strings.TrimSpace(c.QueryParam("token"))
	if token == "" {
		return common.HttpBadRequest(c, errors.New("missing token"))
	}

	logger.Infof("Attemping to one-time signin for token %s...", token)

	if err := env.sessionService.WithContext(c).IssueOneTimeSession(c, token); err != nil {
		return common.HttpError(c, http.StatusUnauthorized, err)
	}

	return c.Redirect(http.StatusTemporaryRedirect, "/")
}
