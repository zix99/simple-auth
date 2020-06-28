package auth

import "github.com/gin-gonic/gin"

func hError(err error) gin.H {
	return gin.H {
		"error": true,
		"message": err.Error(),
	}
}
