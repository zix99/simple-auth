package middleware

import (
	"encoding/base64"
	"errors"
	"net/http"
	"net/url"
	"simple-auth/pkg/appcontext"
	"simple-auth/pkg/config"
	"simple-auth/pkg/db"
	"simple-auth/pkg/routes/middleware/selector/auth"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

const (
	gatewayAccountHeader = "X-SA-Account"
	authorizationHeader  = "Authorization"
	authorizationBasic   = "basic"
)

func newRoundRobinBalancer(targets ...string) middleware.ProxyBalancer {
	var proxyTargets []*middleware.ProxyTarget

	for _, target := range targets {
		url, err := url.Parse(target)
		if err != nil {
			logrus.Warnf("Unable to parse proxy target %s: %v", target, err)
		} else {
			pt := &middleware.ProxyTarget{
				URL: url,
			}
			proxyTargets = append(proxyTargets, pt)
		}
	}

	return middleware.NewRoundRobinBalancer(proxyTargets)
}

func AuthenticationGateway(gateway *config.ConfigLoginGateway, cookieConfig *config.ConfigLoginCookie) echo.MiddlewareFunc {
	const targetKey = "target"

	proxyConfig := middleware.ProxyConfig{
		ContextKey: targetKey,
		Balancer:   newRoundRobinBalancer(gateway.Targets...),
		Rewrite:    gateway.Rewrite,
		Skipper:    middleware.DefaultSkipper,
	}
	balancer := middleware.ProxyWithConfig(proxyConfig)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			log := appcontext.GetLogger(c)
			req := c.Request()

			var subject string
			if authHeader := req.Header.Get(authorizationHeader); authHeader != "" && gateway.BasicAuth {
				// Check authorization header (if allowed) for API access
				uname, pass, err := decodeBasicAuth(authHeader)
				if err != nil {
					return c.HTML(http.StatusBadRequest, err.Error())
				}
				auth, err := validateBasicCredentials(c, uname, pass)
				if err != nil {
					return c.HTML(http.StatusUnauthorized, err.Error())
				}

				req.Header.Del(authorizationHeader)
				subject = auth.Account().UUID
			} else {
				// Check for session
				claims, err := auth.ParseContextSession(cookieConfig, c)
				if err != nil {
					// Not logged in, pass-through to self
					return next(c)
				}
				subject = claims.Subject
			}

			// Special logout path
			if gateway.LogoutPath != "" && req.URL.Path == gateway.LogoutPath {
				auth.ClearSession(c, cookieConfig)
				return c.Redirect(http.StatusTemporaryRedirect, "/")
			}

			// Headers
			req.Header.Set(gatewayAccountHeader, subject)
			for k, v := range gateway.Headers {
				log.Debugf("Override header %s = %s", k, v)
				req.Header.Set(k, v)
			}
			if gateway.Host != "" {
				req.Host = gateway.Host
			}

			// Try to bust the cache
			resp := c.Response()
			if gateway.NoCache {
				resp.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			}

			// Proxy
			ret := balancer(next)(c)
			realTarget := c.Get(targetKey).(*middleware.ProxyTarget)
			log.Infof("PROXY %s %s -> %s", req.Method, req.RequestURI, realTarget.URL.String())
			return ret
		}
	}
}

func decodeBasicAuth(headerValue string) (uname, pass string, err error) {
	authHeaderParts := strings.Fields(headerValue)
	if len(authHeaderParts) != 2 {
		return "", "", errors.New("expected auth type and value")
	}
	if strings.ToLower(authHeaderParts[0]) != authorizationBasic {
		return "", "", errors.New("expected basic auth")
	}

	decoded, err := base64.StdEncoding.DecodeString(authHeaderParts[1])
	if err != nil {
		return "", "", errors.New("unable to base64 decode credentials")
	}

	decodedParts := strings.SplitN(string(decoded), ":", 2)
	if len(decodedParts) != 2 {
		return "", "", errors.New("malformed username and password")
	}

	return decodedParts[0], decodedParts[1], nil
}

func validateBasicCredentials(c appcontext.Context, username, password string) (*db.AuthLocal, error) {
	sadb := appcontext.GetSADB(c)

	errInvalidCredentials := errors.New("invalid username or password")

	authLocal, err := sadb.FindAuthLocalByUsername(username)
	if err != nil {
		return nil, errInvalidCredentials
	}

	if !authLocal.VerifyPassword(password) {
		return nil, errInvalidCredentials
	}

	if authLocal.HasTOTP() {
		return nil, errors.New("unable to validate 2fa")
	}

	return authLocal, nil
}
