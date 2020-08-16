package providers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"simple-auth/pkg/config"
	"simple-auth/pkg/db"
	"simple-auth/pkg/routes/common"
	"simple-auth/pkg/routes/middleware"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

const (
	oidcStateCookieName = "oidcState"
	continueCookieName  = "continue"
)

type OIDCController struct {
	baseURL      string
	id           string
	loginConfig  *config.ConfigLoginSettings
	oidcConfig   *config.ConfigOIDCProvider
	cookieConfig *config.ConfigLoginCookie
	db           db.SADB
}

func NewOIDCController(baseURL, id string, loginConfig *config.ConfigLoginSettings, oidcConfig *config.ConfigOIDCProvider, cookieConfig *config.ConfigLoginCookie, sadb db.SADB) *OIDCController {
	return &OIDCController{
		baseURL:      baseURL,
		id:           id,
		loginConfig:  loginConfig,
		oidcConfig:   oidcConfig,
		cookieConfig: cookieConfig,
		db:           sadb,
	}
}

func (env *OIDCController) Mount(group *echo.Group) {
	logrus.Infof("Enabling OIDC login for %s", env.id)
	group.GET("/"+env.id, env.routeAuthRedirect)
	group.GET("/"+env.id+"/callback", env.routeAuthCallback)
}

func (env *OIDCController) routeAuthRedirect(c echo.Context) error {
	oidcExpiration := time.Now().Add(5 * time.Minute)

	// Compute continue URL
	continueURL := env.loginConfig.ResolveContinueURL(c.QueryParam("continue"))
	if continueURL != "" {
		c.SetCookie(&http.Cookie{
			Name:     continueCookieName,
			Expires:  oidcExpiration,
			HttpOnly: true,
			Value:    continueURL,
			Path:     "/",
		})
	}

	// Parse and redirect to OIDC provider
	redirectURL, err := url.Parse(env.oidcConfig.AuthURL)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, common.JsonErrorf("Invalid auth URL in config"))
	}

	state := uuid.New().String()
	c.SetCookie(&http.Cookie{
		Name:     oidcStateCookieName,
		Expires:  oidcExpiration,
		HttpOnly: true,
		Value:    state,
		Path:     "/",
	})

	qp := redirectURL.Query()
	qp.Set("response_type", "code")
	qp.Set("client_id", env.oidcConfig.ClientID)
	qp.Set("scope", "openid email")
	qp.Set("redirect_uri", env.buildRedirectUri())
	qp.Set("state", state)
	qp.Set("nonce", uuid.New().String())
	redirectURL.RawQuery = qp.Encode()

	return c.Redirect(http.StatusTemporaryRedirect, redirectURL.String())
}

func (env *OIDCController) routeAuthCallback(c echo.Context) error {
	state := c.QueryParam("state")
	code := c.QueryParam("code")

	if code == "" {
		return c.JSON(http.StatusBadRequest, common.JsonErrorf("Invalid code"))
	}

	stateCookie, err := c.Cookie(oidcStateCookieName)
	if err != nil {
		logrus.Warn("Unable to find state cookie. Forgery?")
		return c.JSON(http.StatusBadRequest, common.JsonError(err))
	}
	if stateCookie.Value != state {
		logrus.Warn("State cookie mismatch. Forgery?")
		return c.JSON(http.StatusUnauthorized, common.JsonErrorf("Invalid state"))
	}

	// Clear the state cookie
	c.SetCookie(&http.Cookie{
		Name:     oidcStateCookieName,
		Expires:  time.Now(),
		Path:     "/",
		HttpOnly: true,
	})

	// Resolve continue-url
	continueURL := env.getContinuationURL(c)

	// Trade in the token
	token, err := env.tradeCodeForToken(code)
	if err != nil {
		logrus.Warnf("Error trading in OIDC code for token: %s", err)
		return c.JSON(http.StatusUnauthorized, common.JsonError(err))
	}

	// And parse the token
	type oidcClaims struct {
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		jwt.StandardClaims
	}
	parsedToken, _ := jwt.ParseWithClaims(token, &oidcClaims{}, nil)
	if parsedToken == nil {
		logrus.Warnf("Error parsing claims")
		return c.JSON(http.StatusBadRequest, common.JsonErrorf("Error parsing claims"))
	}
	claims := parsedToken.Claims.(*oidcClaims)

	// Otherwise, we have a token!! Two possible actions
	// 1) If token is associated with an account, great! we're logged-in
	// 2) If the user is logged-in to SA, associate OIDC provider with their account; or
	// 3) If the user isn't logged in, create a new account and associate
	// WARN: Do not associate token with email.  Someone else may have created simple-auth email fraudulently

	// TODO: Check if user already logged in (and associate)

	// Check if exists
	{
		account, _ := env.db.FindAccountForOIDC(env.id, claims.Subject)
		if account != nil {
			middleware.CreateSession(c, env.cookieConfig, account, middleware.SessionSourceOIDC)
			return c.Redirect(http.StatusTemporaryRedirect, continueURL)
		}
	}

	// If not, try to create it
	if env.loginConfig.CreateAccountEnabled {
		account, err := env.db.CreateAccount(claims.Email)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, common.JsonError(err))
		}
		err2 := env.db.CreateOIDCForAccount(account, env.id, claims.Subject)
		if err2 != nil {
			return c.JSON(http.StatusInternalServerError, common.JsonError(err2))
		}
		middleware.CreateSession(c, env.cookieConfig, account, middleware.SessionSourceOIDC)
		return c.Redirect(http.StatusTemporaryRedirect, continueURL)
	}

	return c.JSON(http.StatusForbidden, common.JsonErrorf("Unable to create new OIDC for user. Account creation disabled."))
}

func (env *OIDCController) tradeCodeForToken(code string) (string, error) {
	form := url.Values{
		"code":          {code},
		"client_id":     {env.oidcConfig.ClientID},
		"client_secret": {env.oidcConfig.ClientSecret},
		"redirect_uri":  {env.buildRedirectUri()},
		"grant_type":    {"authorization_code"},
	}
	resp, err := http.PostForm(env.oidcConfig.TokenURL, form)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Parse it
	var contents map[string]string
	json.Unmarshal(body, &contents)

	if contents["id_token"] == "" {
		return "", errors.New("Invalid id_token")
	}

	return contents["id_token"], nil
}

func (env *OIDCController) buildRedirectUri() string {
	return env.baseURL + "/" + env.id + "/callback"
}

func (env *OIDCController) getContinuationURL(c echo.Context) string {
	if continueCookie, err := c.Cookie(continueCookieName); err == nil && continueCookie.Value != "" {
		c.SetCookie(&http.Cookie{
			Name:     continueCookieName,
			Expires:  time.Now(),
			Path:     "/",
			HttpOnly: true,
		})
		return continueCookie.Value
	}
	if env.loginConfig.RouteOnLogin != "" {
		return env.loginConfig.RouteOnLogin
	}
	return "/"
}
