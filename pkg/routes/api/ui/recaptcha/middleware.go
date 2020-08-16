package recaptcha

import (
	"net/http"
	"simple-auth/pkg/routes/common"

	"github.com/labstack/echo/v4"
)

// MiddlewareV2 looks for a recaptcha value in the request, and blocks a response if the value is not valid
// Token must be on query, can't read in body
func MiddlewareV2(secret string) echo.MiddlewareFunc {
	validator := NewValidatorV2(secret)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.QueryParam("recaptchav2")

			if token == "" {
				return common.HttpErrorf(c, http.StatusBadRequest, "Missing recaptchav2 value")
			}

			if err := validator.Validate(token); err != nil {
				return common.HttpError(c, http.StatusForbidden, err)
			}

			return next(c)
		}
	}
}
