package auth

import (
	"simple-auth/pkg/config"
	"simple-auth/pkg/routes/common"
	"time"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

/*
This authentication schema is for authentication when there are 2 parties
that want to validate, with a reduction of trust in each party. (eg, a game launcher -> a game -> a game server)

Based on a simple authentication scheme (Must have password)

In this case, the following will happen:
	1. The game launcher retrieves a "session" token using the auth information for whatever schema (eg simple).  The session token is passed to the client
	   NOTE: Only one session token can be activated a time.  If another session token is claimed, all existing tokens become invalid
	2. Upon joining the server, the client will attempt to trade the session token for a verification token.  The verification token is
	   useless except for a 3rd party to verify that it is "valid"
	3. The server, having the userId and verification token can validate that the two belong to each other, and are current and valid
*/

type responseToken struct {
	Token string `json:"token"`
}

// routeUser validates a user and issues a account-token
// only one session can be active at a given time
func (env *environment) routeIssueSessionToken(c echo.Context) error {
	req := struct {
		Username string `json:"username" form:"username"`
		Password string `json:"password" form:"password"`
	}{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, common.JsonError(err))
	}

	expireDuration := time.Duration(config.Global.Authenticators.Token.SessionExpiresMinutes) * time.Minute
	token, err := env.db.AssertCreateSessionToken(req.Username, req.Password, expireDuration)
	if err != nil {
		return c.JSON(401, common.JsonErrorf("Unable to create session token"))
	}

	return c.JSON(200, responseToken{
		Token: token,
	})
}

func (env *environment) routeIssueVerificationToken(c echo.Context) error {
	req := struct {
		Username string `json:"username" form:"username"`
		Token    string `json:"token" form:"token"`
	}{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, common.JsonError(err))
	}

	vToken, err := env.db.CreateVerificationToken(req.Username, req.Token)
	if err != nil {
		logrus.Error(err)
		return c.JSON(401, common.JsonErrorf("Unable to create verification token"))
	}

	logrus.Infof("Issuing verification token for %s", req.Username)
	return c.JSON(200, responseToken{
		Token: vToken,
	})
}

func (env *environment) routeVerifyToken(c echo.Context) error {
	req := struct {
		Username string `json:"username" form:"username"`
		Token    string `json:"token" form:"token"`
	}{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, common.JsonError(err))
	}

	account, err := env.db.AssertVerificationToken(req.Username, req.Token)
	if err != nil {
		return c.JSON(401, common.JsonErrorf("Verification token not found"))
	}
	return c.JSON(200, common.Json{
		"username":   req.Username,
		"account_id": account.UUID,
	})
}
