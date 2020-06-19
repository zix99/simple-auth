package auth

import (
	"simple-auth/pkg/db"

	"github.com/gin-gonic/gin"
)

type environment struct {
	db db.SADB
}

func NewRouter(group *gin.RouterGroup, db db.SADB) {
	env := &environment{
		db: db,
	}

	group.POST("/simple", env.routeSimpleAuthenticate)

	setupSessionAuthenticator(env, group.Group("/session"))
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
			"id": account.ID,
		})
	}
}
