package ui

import (
	"errors"
	"simple-auth/pkg/api/common"
	"simple-auth/pkg/api/ui/recaptcha"
	"unicode/utf8"

	"github.com/labstack/echo"
)

type createAccountRequest struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Email       string `json:"email" binding:"required"`
	RecaptchaV2 string `json:"recaptchav2" binding:"required"`
}

func (env *environment) routeCreateAccount(c echo.Context) error {
	req := createAccountRequest{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, common.JsonError(errors.New("Unable to deserialize request")))
	}

	if err := env.validateUsername(req.Username); err != nil {
		return c.JSON(400, common.JsonError(err))
	}
	if err := env.validatePassword(req.Password); err != nil {
		return c.JSON(400, common.JsonError(err))
	}
	if err := env.validateEmail(req.Email); err != nil {
		return c.JSON(400, common.JsonError(err))
	}

	// Validate recaptcha if needed
	if err := env.validateRecaptchaV2(req.RecaptchaV2); err != nil {
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

func (env *environment) validateUsername(username string) error {
	ulen := utf8.RuneCountInString(username)
	if ulen < env.config.Requirements.UsernameMinLength {
		return errors.New("Username too short")
	}
	if ulen > env.config.Requirements.UsernameMaxLength {
		return errors.New("Username too long")
	}
	return nil
}

func (env *environment) validatePassword(password string) error {
	plen := utf8.RuneCountInString(password)
	if plen < env.config.Requirements.PasswordMinLength {
		return errors.New("Password too short")
	}
	if plen > env.config.Requirements.PasswordMaxLength {
		return errors.New("Password too long")
	}
	return nil
}

func (env *environment) validateEmail(email string) error {
	elen := utf8.RuneCountInString(email)
	if elen < 5 { // Must be at least: a@b.c
		return errors.New("Email too short")
	}
	return nil
}

func (env *environment) validateRecaptchaV2(code string) error {
	if !env.config.RecaptchaV2.Enabled {
		return nil
	}

	validator := recaptcha.NewValidatorV2(env.config.RecaptchaV2.Secret)
	return validator.Validate(code)
}
