package v1

import (
	"net/http"
	"simple-auth/pkg/appcontext"
	"simple-auth/pkg/db"
	"simple-auth/pkg/routes/common"
	"simple-auth/pkg/routes/middleware/selector/auth"
	"simple-auth/pkg/saerrors"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	errorInvalidAccount saerrors.ErrorCode = "invalid-account"
)

type (
	getAccountOIDCAuthResponse struct {
		Provider string `json:"provider"`
		Subject  string `json:"subject"`
		Icon     string `json:"icon"`
		Name     string `json:"name"`
	}
	getAccountAuthProviderResponse struct {
		Local *getLocalLoginResponse        `json:"local,omitempty"`
		OIDC  *[]getAccountOIDCAuthResponse `json:"oidc,omitempty"`
	}
	getAccountResponse struct {
		ID      string                         `json:"id" example:"00000000-0000-0000-0000-000000000000"`
		Created time.Time                      `json:"created"`
		Email   string                         `json:"email" example:"sa@example.com"`
		Name    string                         `json:"name" example:"John Smith"`
		Auth    getAccountAuthProviderResponse `json:"auth"`
	}
)

// RouteGetAccount gets account info
// @Summary Get Account
// @Tags Account
// @Description Get details about account
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} getAccountResponse
// @Failure 400,401,404,500 {object} common.ErrorResponse
// @Router /account [get]
func (env *Environment) RouteGetAccount(c echo.Context) error {
	logger := appcontext.GetLogger(c)
	sadb := appcontext.GetSADB(c)
	accountUUID := auth.MustGetAccountUUID(c)

	logger.Infof("Get account for %s", accountUUID)
	account, err := sadb.FindAccount(accountUUID)
	if err != nil {
		return common.HttpError(c, http.StatusInternalServerError, errorInvalidAccount.Wrapf(err, "Logged in with unknown account"))
	}

	response := getAccountResponse{
		ID:      account.UUID,
		Created: account.CreatedAt,
		Email:   account.Email,
		Name:    account.Name,
	}

	if localLoginResponse, err := env.getLocalLoginResponse(c); err != nil {
		logger.Warnln(err)
	} else {
		response.Auth.Local = localLoginResponse
	}

	if providers, err := sadb.FindOIDCForAccount(account); err == nil && len(providers) > 0 {
		oidcProviders := make([]getAccountOIDCAuthResponse, len(providers))
		for i, oidc := range providers {
			providerConfig := env.oidcService.GetProvider(oidc.Provider)
			oidcProviders[i] = getAccountOIDCAuthResponse{
				Provider: oidc.Provider,
				Subject:  oidc.Subject,
				Icon:     providerConfig.Icon,
				Name:     providerConfig.Name,
			}
		}
		response.Auth.OIDC = &oidcProviders
	}

	return c.JSON(http.StatusOK, response)
}

type (
	getAccountAuditRecordResponse struct {
		Timestamp time.Time      `json:"ts"`
		Module    db.AuditModule `json:"module"`
		Level     db.AuditLevel  `json:"level"`
		Message   string         `json:"message"`
	}
	getAccountAuditResponse struct {
		Records []*getAccountAuditRecordResponse `json:"records"`
	}
)

func (env *Environment) RouteGetAccountAudit(c echo.Context) error {
	accountUUID := auth.MustGetAccountUUID(c)
	logger := appcontext.GetLogger(c)
	sadb := appcontext.GetSADB(c)

	logger.Infof("Get account audit for %s", accountUUID)

	account, err := sadb.FindAccount(accountUUID)
	if err != nil {
		return common.HttpError(c, http.StatusInternalServerError, errorInvalidAccount.Wrapf(err, "Logged in with unknown account"))
	}

	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	records, err := sadb.GetAuditTrailForAccount(account, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	ret := make([]*getAccountAuditRecordResponse, len(records))
	for i, record := range records {
		ret[i] = &getAccountAuditRecordResponse{
			Timestamp: record.CreatedAt,
			Module:    record.Module,
			Level:     record.Level,
			Message:   record.Message,
		}
	}

	return c.JSON(http.StatusOK, getAccountAuditResponse{
		Records: ret,
	})
}
