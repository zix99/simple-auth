package ui

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

	group.POST("/account", env.routeCreateAccount)
}

type createAccountModel struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (env *environment) routeCreateAccount(c *gin.Context) {
	var model createAccountModel
	c.BindJSON(&model)

	account, err := env.db.CreateAccount(model.Email)
	if err != nil {
		c.AbortWithStatusJSON(400, hError(err))
		return
	}

	err2 := env.db.CreateAccountAuthSimple(account, model.Username, model.Password)
	if err2 != nil {
		c.AbortWithStatusJSON(500, hError(err2))
		return
	}

	c.JSON(201, gin.H{
		"id": account.UUID,
	})
}

func hError(err error) gin.H {
	return gin.H{"message": err.Error()}
}
