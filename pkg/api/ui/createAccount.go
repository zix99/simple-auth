package ui

import (
	"errors"
	"simple-auth/pkg/api/common"
	"simple-auth/pkg/config"
	"unicode/utf8"

	"github.com/labstack/echo"
)

type createAccountRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

func (env *environment) routeCreateAccount(c echo.Context) error {
	req := createAccountRequest{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, common.JsonError(errors.New("Unable to deserialize request")))
	}

	if err := validateUsername(req.Username); err != nil {
		return c.JSON(400, common.JsonError(err))
	}
	if err := validatePassword(req.Password); err != nil {
		return c.JSON(400, common.JsonError(err))
	}
	if err := validateEmail(req.Email); err != nil {
		return c.JSON(400, common.JsonError(err))
	}

	account, err := env.db.CreateAccount(req.Email)
	if err != nil {
		return c.JSON(400, common.JsonError(err))
	}

	err2 := env.db.CreateAccountAuthSimple(account, req.Username, req.Password)
	if err2 != nil {
		return c.JSON(500, common.JsonError(err2))
	}

	return c.JSON(201, map[string]string{
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
