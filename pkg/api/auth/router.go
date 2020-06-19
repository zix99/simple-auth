package auth

import (
	"simple-auth/pkg/db"

	"github.com/gin-gonic/gin"
)

func NewRouter(group *gin.RouterGroup, db db.AccountStore) {
	group.POST("/user/simple/authenticate", routeSimpleAuthenticate)
}

func routeSimpleAuthenticate(c *gin.Context) {
	c.Status(200)
}
