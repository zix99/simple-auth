package auth

import (
	"github.com/gin-gonic/gin"
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

func setupSessionAuthenticator(env *environment, g *gin.RouterGroup) {
	logrus.Info("Enabling session auth...")
	g.POST("/", env.routeIssueAccountToken)
	g.POST("/session", env.routeIssueSessionToken)
	g.POST("/session/verify", env.routeSessionVerify)
}

// routeUser validates a user and issues a account-token
// only one session can be active at a given time
func (env *environment) routeIssueAccountToken(c *gin.Context) {
	req := struct {
		username string
		password string
	}{}
	if err := c.Bind(&req); err != nil {
		c.AbortWithStatusJSON(400, hError(err))
		return
	}

}

func (env *environment) routeIssueSessionToken(c *gin.Context) {
	// accountToken := c.PostForm("account-token")

}

func (env *environment) routeSessionVerify(c *gin.Context) {

}