package auth

import (
	"fmt"
	"net/http"
	"simple-auth/pkg/appcontext"
	"simple-auth/pkg/config"
	"simple-auth/pkg/db"
	"simple-auth/pkg/routes/common"
	"simple-auth/pkg/routes/middleware/selector/auth"
	"simple-auth/pkg/services"
	"time"

	"github.com/labstack/echo/v4"
)

type OAuth2Error string

const (
	InvalidRequest       OAuth2Error = "invalid_request"
	InvalidClient        OAuth2Error = "invalid_client"
	InvalidGrant         OAuth2Error = "invalid_grant"
	InvalidScope         OAuth2Error = "invalid_scope"
	UnauthorizedClient   OAuth2Error = "unauthorized_client"
	UnsupportedGrantType OAuth2Error = "unsupported_grant_type"
	InternalError        OAuth2Error = "server_error"
)

type oauth2Error struct {
	Error       OAuth2Error `json:"error"`
	Description string      `json:"error_description"`
}

type OAuth2Controller struct {
	config        *config.ConfigOAuth2
	oauthServices map[string]services.AuthOAuthService
}

func NewOAuth2Controller(config *config.ConfigOAuth2, localLoginService services.LocalLoginService) *OAuth2Controller {
	oauthServices := make(map[string]services.AuthOAuthService)
	for clientID, cfg := range config.Clients {
		oauthServices[clientID] = services.NewAuthOAuthService(clientID, cfg, &config.Settings, localLoginService)
	}

	return &OAuth2Controller{
		config,
		oauthServices,
	}
}

type oauth2ClientResponse struct {
	Name      string `json:"name"`
	Author    string `json:"author"`
	AuthorURL string `json:"author_url"`
}

// @Summary Client Info
// @Description Gets information about oauth2 client
// @Tags Auth
// @Accept json
// @Produce json
// @Param client_id path string true "OAuth2 Client ID"
// @Success 200 {object} oauth2ClientResponse
// @Failure 404,500 {object} oauth2Error
// @Router /auth/oauth2/client/{client_id} [get]
func (s *OAuth2Controller) RouteClientInfo(c echo.Context) error {
	clientID := c.Param("client_id")

	if client, ok := s.config.Clients[clientID]; ok {
		return c.JSON(http.StatusOK, &oauth2ClientResponse{
			Name:      client.Name,
			Author:    client.Author,
			AuthorURL: client.AuthorURL,
		})
	}

	return c.JSON(http.StatusNotFound, &oauth2Error{
		Error:       InvalidClient,
		Description: "Client id not found",
	})
}

type oauth2GetTokenResponseToken struct {
	ClientID   string    `json:"client_id"`
	ClientName string    `json:"client_name"`
	ShortToken string    `json:"short_token"` // A short/obfuscated version of the full token for identification purposes
	Type       string    `json:"type"`
	Created    time.Time `json:"created"`
	Expires    time.Time `json:"expires"`
}

type oauth2GetTokensResponse struct {
	Tokens []*oauth2GetTokenResponseToken `json:"tokens"`
}

// @Summary Get tokens
// @Description Get tokens associated with account
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} oauth2GetTokensResponse
// @Failure 400,404,500 {object} common.ErrorResponse
// @Router /auth/oauth2/client [get]
func (s *OAuth2Controller) RouteGetTokens(c echo.Context) error {
	accountUUID := auth.MustGetAccountUUID(c)
	sadb := appcontext.GetSADB(c)

	account, err := sadb.FindAccount(accountUUID)
	if err != nil {
		return common.HttpInternalError(c, err)
	}

	tokens, err := sadb.GetAllValidOAuthTokens(account)
	if err != nil {
		return common.HttpInternalError(c, err)
	}

	ret := make([]*oauth2GetTokenResponseToken, len(tokens))
	for i, t := range tokens {
		ret[i] = &oauth2GetTokenResponseToken{
			ClientID:   t.ClientID,
			ClientName: s.config.Clients[t.ClientID].Name,
			ShortToken: t.Token[:5],
			Type:       string(t.Type),
			Created:    t.Created,
			Expires:    t.Expires,
		}
	}

	return c.JSON(http.StatusOK, oauth2GetTokensResponse{
		Tokens: ret,
	})
}

// @Summary Revoke Token
// @Description Revoke tokens for a given client_id
// @Tags Auth
// @Param client_id query string true "Client ID to revoke"
// @Param token query string false "Optional specific token to revoke"
// @Success 200 {object} oauth2GetTokensResponse
// @Failure 400,404,500 {object} common.ErrorResponse
// @Router /auth/oauth2/token [delete]
func (s *OAuth2Controller) RouteRevokeToken(c echo.Context) error {
	clientID := c.QueryParam("client_id")
	token := c.QueryParam("token")
	if clientID == "" {
		return oauthError(c, InvalidRequest, "Missing client_id")
	}

	accountUUID := auth.MustGetAccountUUID(c)
	sadb := appcontext.GetSADB(c)

	account, err := sadb.FindAccount(accountUUID)
	if err != nil {
		return common.HttpInternalError(c, err)
	}

	if token != "" {
		if err := sadb.InvalidateToken(clientID, account, token); err != nil {
			return common.HttpInternalError(c, err)
		}
	} else {
		if err := sadb.InvalidateAllOAuth(clientID, account); err != nil {
			return common.HttpInternalError(c, err)
		}
	}

	return common.HttpOK(c)
}

type authorizedGrantRequest struct {
	ClientID     string `json:"client_id" validate:"required"`
	ResponseType string `json:"response_type" validate:"required"`
	Scope        string `json:"scope"`
	RedirectURI  string `json:"redirect_uri" validate:"required"`
	State        string `json:"state"`
	Auto         bool   `json:"auto"` // If it's an auto-grant request
}

type authorizedGrantResponse struct {
	Code  string `json:"code"`
	State string `json:"state,omitempty"`
}

// @Summary Authentication Grant Code
// @Description Called by UI to authorized a grant token. MUST pass CSRF
//              If you need an alternative, use the password grant type to obtain access token
// @Tags Auth
// @Accept json
// @Produce json
// @Param authorizedGrantRequest body authorizedGrantRequest true "body"
// @Success 200 {object} authorizedGrantResponse
// @Failure 400,404,500 {object} oauth2Error
// @Router /auth/oauth2/grant [post]
func (s *OAuth2Controller) RouteAuthorizedGrantCode(c echo.Context) error {
	log := appcontext.GetLogger(c)

	var req authorizedGrantRequest
	if err := c.Bind(&req); err != nil {
		return oauthError(c, InvalidRequest, err.Error())
	}
	if err := c.Validate(&req); err != nil {
		return oauthError(c, InvalidRequest, err.Error())
	}
	if req.ResponseType != "code" {
		return oauthError(c, InvalidRequest, "Expected response_type to be code")
	}

	if client, ok := s.config.Clients[req.ClientID]; ok {
		if client.RedirectURI != req.RedirectURI {
			return oauthError(c, InvalidRequest, "Unknown redirect URI")
		}
	} else {
		return oauthError(c, InvalidClient, "Unknown client id %s", req.ClientID)
	}

	uuid := auth.MustGetAccountUUID(c)
	sadb := appcontext.GetSADB(c)

	account, err := sadb.FindAccount(uuid)
	if err != nil {
		return oauthError(c, InvalidRequest, "No session")
	}

	oauthService := s.oauthServices[req.ClientID].WithContext(c)
	scopes := db.NewOAuthScope(req.Scope)

	if !oauthService.ValidateScopes(scopes) {
		return oauthError(c, InvalidScope, "Invalid scopes")
	}

	if req.Auto {
		if !s.config.Settings.AllowAutoGrant {
			return oauthError(c, InvalidRequest, "Auto granting not allowed")
		}
		_, err := oauthService.FindExistingToken(account, db.OAuthTypeAccessToken, scopes)
		if err != nil {
			return oauthError(c, InvalidGrant, "no existing token")
		}
		// Otherwise, fall-through to generate access token
		log.Infof("Allowing auto-grant for client %s, account %s...", req.ClientID, account.UUID)
	}

	code, err := oauthService.CreateAccessCode(account, scopes)
	if err != nil {
		return oauthError(c, InternalError, err.Error())
	}

	log.Infof("Issued code to %s", account.UUID)
	return c.JSON(http.StatusOK, &authorizedGrantResponse{
		Code:  code,
		State: req.State,
	})
}

type grantTokenRequest struct {
	GrantType string `form:"grant_type" json:"grant_type" query:"grant_type" validate:"required"`

	// grantType == authorization_code
	Code        string `form:"code" json:"code"`
	RedirectURI string `form:"redirect_uri" json:"redirect_uri"`

	// grantType == "password"
	Username string  `form:"username" json:"username"`
	Password string  `form:"password" json:"password"`
	Totp     *string `form:"totp" json:"totp"`
	Scope    string  `form:"scope" json:"scope"`

	// grantType == "refresh_token"
	RefreshToken string `form:"refresh_token" json:"refresh_token"`

	// General
	ClientID     string `form:"client_id" json:"client_id" validate:"required"`
	ClientSecret string `form:"client_secret" json:"client_secret"`
}

type grantTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in,omitempty"` // Seconds access token expires in
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

// @Summary Grant Token
// @Description Grants a token given a various grant_type
// @Tags Auth
// @Accept json
// @Produce json
// @Param grantTokenRequest body grantTokenRequest true "body"
// @Success 200 {object} grantTokenResponse
// @Failure 400,404,500 {object} oauth2Error
// @Router /auth/oauth2/token [post]
func (s *OAuth2Controller) RouteTokenGrant(c echo.Context) error {
	var req grantTokenRequest
	if err := c.Bind(&req); err != nil {
		return oauthError(c, InvalidRequest, err.Error())
	}
	if err := c.Validate(&req); err != nil {
		return oauthError(c, InvalidRequest, err.Error())
	}

	clientService, hasClientService := s.oauthServices[req.ClientID]
	if !hasClientService {
		return oauthError(c, InvalidClient, "Unknown client id %s", req.ClientID)
	}
	clientService = clientService.WithContext(c)

	switch req.GrantType {
	case "password":
		return s.routeTokenGrantPassword(c, clientService, &req)
	case "authorization_code": // code code for token
		return s.routeTokenGrantAuthorizationCode(c, clientService, &req)
	case "refresh_token": // re-issue new access token with refresh
		return s.routeTokenGrantRefreshToken(c, clientService, &req)
	}

	return oauthError(c, UnsupportedGrantType, "Unknown grant type: %s", req.GrantType)
}

// routeTokenGrantPassword allows oauth2 grant via username, password (and totp)
func (s *OAuth2Controller) routeTokenGrantPassword(c echo.Context, clientService services.AuthOAuthService, req *grantTokenRequest) error {
	scopes := db.NewOAuthScope(req.Scope)

	if !clientService.ValidateScopes(db.NewOAuthScope(req.Scope)) {
		return oauthError(c, InvalidScope, "Invald scopes")
	}

	retToken, err := clientService.TradeCredentialsForToken(req.ClientSecret, req.Username, req.Password, req.Totp, scopes)
	if err != nil {
		return oauthError(c, InvalidRequest, err.Error())
	}

	return c.JSON(http.StatusOK, &grantTokenResponse{
		AccessToken:  retToken.AccessToken,
		RefreshToken: retToken.RefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    retToken.Expires,
		Scope:        retToken.Scope.String(),
	})
}

// routeTokenGrantAuthorizationCode grants a token (and maybe refresh token) by trading a code from the auth flow
func (s *OAuth2Controller) routeTokenGrantAuthorizationCode(c echo.Context, clientService services.AuthOAuthService, req *grantTokenRequest) error {
	if !clientService.ValidateRedirectURI(req.RedirectURI) {
		return oauthError(c, UnauthorizedClient, "Invalid redirect_uri")
	}

	retToken, err := clientService.TradeCodeForToken(req.ClientSecret, req.Code)
	if err != nil {
		return oauthError(c, InternalError, err.Error())
	}

	return c.JSON(http.StatusOK, &grantTokenResponse{
		AccessToken:  retToken.AccessToken,
		RefreshToken: retToken.RefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    retToken.Expires,
		Scope:        retToken.Scope.String(),
	})
}

// routeTokenGrantRefreshToken trades the refresh token for a new access token
func (s *OAuth2Controller) routeTokenGrantRefreshToken(c echo.Context, clientService services.AuthOAuthService, req *grantTokenRequest) error {
	retToken, err := clientService.TradeRefreshTokenForAccessToken(req.ClientSecret, req.RefreshToken)
	if err != nil {
		return oauthError(c, InternalError, err.Error())
	}

	return c.JSON(http.StatusOK, &grantTokenResponse{
		AccessToken: retToken.AccessToken,
		TokenType:   "Bearer",
		ExpiresIn:   retToken.Expires,
		Scope:       retToken.Scope.String(),
	})
}

func oauthError(c echo.Context, code OAuth2Error, msg string, args ...interface{}) error {
	log := appcontext.GetLogger(c)
	fullMsg := fmt.Sprintf(msg, args...)
	log.Warnf("Error issuing OAuth2 token '%s': %s", code, fullMsg)
	return c.JSON(http.StatusBadRequest, &oauth2Error{
		Error:       code,
		Description: fullMsg,
	})
}
