package middleware

import (
	"net/http"
	"net/url"
	"simple-auth/pkg/config"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
)

const GatewayAccountHeader = "X-SA-Account"

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
			claims, err := parseContextSession(&cookieConfig.JWT, c)
			if err != nil {
				// Not logged in, pass-through to self
				return next(c)
			}

			req := c.Request()

			// Special logout path
			if gateway.LogoutPath != "" && req.URL.Path == gateway.LogoutPath {
				ClearSession(c, cookieConfig)
				return c.Redirect(http.StatusTemporaryRedirect, "/")
			}

			// Headers
			req.Header.Set(GatewayAccountHeader, claims.Subject)
			for k, v := range gateway.Headers {
				logrus.Infof("Override header %s = %s", k, v)
				req.Header.Set(k, v)
			}
			if gateway.Host != "" {
				req.Host = gateway.Host
			}

			// Proxy
			ret := balancer(next)(c)
			realTarget := c.Get(targetKey).(*middleware.ProxyTarget)
			logrus.Infof("PROXY %s %s -> %s", req.Method, req.RequestURI, realTarget.URL.String())
			return ret
		}
	}
}