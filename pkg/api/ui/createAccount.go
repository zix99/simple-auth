package ui

import (
	"errors"
	"simple-auth/pkg/config"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
)

func (env *environment) routeCreateAccount(c *gin.Context) {
	email := c.PostForm("email")
	username := c.PostForm("username")
	password := c.PostForm("password")

	if err := validateUsername(username); err != nil {
		c.AbortWithStatusJSON(400, hError(err))
		return
	}
	if err := validatePassword(password); err != nil {
		c.AbortWithStatusJSON(400, hError(err))
		return
	}

	account, err := env.db.CreateAccount(email)
	if err != nil {
		c.AbortWithStatusJSON(400, hError(err))
		return
	}

	err2 := env.db.CreateAccountAuthSimple(account, username, password)
	if err2 != nil {
		c.AbortWithStatusJSON(500, hError(err2))
		return
	}

	c.JSON(201, gin.H{
		"id": account.UUID,
	})
}

func validateUsername(username string) error {
	ulen := utf8.RuneCountInString(username)
	if ulen < config.Global.Web.Requirements.UsernameMinLength {
		return errors.New("Username too short")
	}
	if ulen > config.Global.Web.Requirements.UsernameMaxLength {
		return errors.New("Username too long")
	}
	return nil
}

func validatePassword(password string) error {
	plen := utf8.RuneCountInString(password)
	if plen < config.Global.Web.Requirements.PasswordMinLength {
		return errors.New("Password too short")
	}
	if plen > config.Global.Web.Requirements.PasswordMaxLength {
		return errors.New("Password too long")
	}
	return nil
}
