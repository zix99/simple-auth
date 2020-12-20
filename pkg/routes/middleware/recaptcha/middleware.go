package recaptcha

import (
	"net/http"
	"simple-auth/pkg/routes/common"
	"simple-auth/pkg/saerrors"

	"github.com/labstack/echo/v4"
)

const (
	errorMissingRecaptchaValue saerrors.ErrorCode = "missing-recaptchv2-value"
	errorInvalidRecaptcha      saerrors.ErrorCode = "invalid-recaptcha"
)

// MiddlewareV2 looks for a "recaptchav2" query value in the request, and blocks a response if the value is not valid
// Token must be on query, can't read in body
func MiddlewareV2(secret string) echo.MiddlewareFunc {
	validator := NewValidatorV2(secret)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.QueryParam("recaptchav2")

			if token == "" {
				return common.HttpError(c, http.StatusBadRequest, errorMissingRecaptchaValue.New())
			}

			if err := validator.Validate(token); err != nil {
				return common.HttpError(c, http.StatusForbidden, errorInvalidRecaptcha.Compose(err))
			}

			return next(c)
		}
	}
}
