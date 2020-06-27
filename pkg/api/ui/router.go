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
