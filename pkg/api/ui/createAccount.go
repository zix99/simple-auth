package ui

import (
	"errors"
	"simple-auth/pkg/config"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
)

type createAccountRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

func (env *environment) routeCreateAccount(c *gin.Context) {
	req := createAccountRequest{}
	if err := c.Bind(&req); err != nil {
		c.AbortWithStatusJSON(400, hError(errors.New("Unable to deserialize request")))
		return
	}

	if err := validateUsername(req.Username); err != nil {
		c.AbortWithStatusJSON(400, hError(err))
		return
	}
	if err := validatePassword(req.Password); err != nil {
		c.AbortWithStatusJSON(400, hError(err))
		return
	}
	if err := validateEmail(req.Email); err != nil {
		c.AbortWithStatusJSON(400, hError(err))
		return
	}

	account, err := env.db.CreateAccount(req.Email)
	if err != nil {
		c.AbortWithStatusJSON(400, hError(err))
		return
	}

	err2 := env.db.CreateAccountAuthSimple(account, req.Username, req.Password)
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

func validateEmail(email string) error {
	elen := utf8.RuneCountInString(email)
	if elen < 5 { // Must be at least: a@b.c
		return errors.New("Email too short")
	}
	return nil
}
