package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

/*

Simple authenticator will simply provide an endpoint that will either return a 200 or a 401
depending on whether the username/password has been validated
*/

func setupSimpleAuthenticator(env *environment, group *gin.RouterGroup) {
	logrus.Info("Enabling simple auth...")
	group.POST("/", env.routeSimpleAuthenticate)
}

func (env *environment) routeSimpleAuthenticate(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	account, err := env.db.FindAndVerifySimpleAuth(username, password)
	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{
			"message": err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"id": account.UUID,
		})
	}
}
