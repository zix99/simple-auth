package ui

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"simple-auth/pkg/db"
	"simple-auth/pkg/email"
	"simple-auth/pkg/routes/api/ui/recaptcha"
	"simple-auth/pkg/routes/common"
	"simple-auth/pkg/routes/middleware"
	"simple-auth/pkg/routes/middleware/selector/auth"
	"unicode/utf8"

	"github.com/labstack/echo/v4"
)

type createAccountRequest struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Email       string `json:"email" binding:"required"`
	RecaptchaV2 string `json:"recaptchav2" binding:"required"`
}

func (env *environment) routeCreateAccount(c echo.Context) error {
	logger := middleware.GetLogger(c)

	req := createAccountRequest{}
	if err := c.Bind(&req); err != nil {
		return common.HttpBadRequest(c, err)
	}

	if err := env.validateUsername(req.Username); err != nil {
		return common.HttpBadRequest(c, err)
	}
	if err := env.validatePassword(req.Password); err != nil {
		return common.HttpBadRequest(c, err)
	}
	if err := env.validateEmail(req.Email); err != nil {
		return common.HttpBadRequest(c, err)
	}

	// Validate recaptcha if needed
	if err := env.validateRecaptchaV2(req.RecaptchaV2); err != nil {
		return common.HttpError(c, http.StatusBadRequest, errorInvalidRecaptcha.Wrap(err))
	}

	account, err := env.db.CreateAccount(req.Email)
	if err != nil {
		return common.HttpBadRequest(c, err)
	}

	_, err2 := env.db.CreateAuthLocal(account, req.Username, req.Password)
	if err2 != nil {
		return common.HttpError(c, http.StatusInternalServerError, err2)
	}

	// trigger email
	go email.NewFromConfig(logger, env.email).SendWelcomeEmail(req.Email, &email.WelcomeEmailData{
		EmailData: email.EmailData{
			Company: env.meta.Company,
			BaseURL: env.config.GetBaseURL(),
		},
		Name: req.Username,
	})

	// Activation
	if env.config.Login.Settings.EmailValidationRequired {
		stip := db.NewTokenStipulation()
		env.db.AddStipulation(account, stip)

		baseURL := env.config.GetBaseURL()
		go email.NewFromConfig(logger, env.email).SendVerificationEmail(req.Email, &email.VerificationData{
			EmailData: email.EmailData{
				Company: env.meta.Company,
				BaseURL: baseURL,
			},
			ActivationLink: template.HTML(fmt.Sprintf("%s/#/activate?account=%s&token=%s", baseURL, account.UUID, stip.Code)),
		})
	}

	// log the user in to a session
	err3 := auth.CreateSession(c, &env.config.Login.Cookie, account, auth.SourceLogin)
	if err3 != nil {
		logger.Warnf("Unable to create session post-login, ignoring: %v", err3)
	}

	return c.JSON(201, common.Json{
		"id": account.UUID,
	})
}

func (env *environment) validateUsername(username string) error {
	ulen := utf8.RuneCountInString(username)
	if ulen < env.config.Requirements.UsernameMinLength {
		return errors.New("username too short")
	}
	if ulen > env.config.Requirements.UsernameMaxLength {
		return errors.New("username too long")
	}

	if env.config.Requirements.UsernameRegex != "" {
		re, err := regexp.Compile(env.config.Requirements.UsernameRegex)
		if err != nil {
			return errors.New("unable to parse valid username regex, ask your server admin to fix this")
		}
		if !re.MatchString(username) {
			return errors.New("invalid username characters")
		}
	}

	return nil
}

func (env *environment) validatePassword(password string) error {
	plen := utf8.RuneCountInString(password)
	if plen < env.config.Requirements.PasswordMinLength {
		return errors.New("password too short")
	}
	if plen > env.config.Requirements.PasswordMaxLength {
		return errors.New("password too long")
	}
	return nil
}

func (env *environment) validateEmail(email string) error {
	elen := utf8.RuneCountInString(email)
	if elen < 5 { // Must be at least: a@b.c
		return errors.New("email too short")
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
