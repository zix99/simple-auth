package ui

import "github.com/gin-gonic/gin"

func hError(err error) gin.H {
	return gin.H{"message": err.Error()}
}
