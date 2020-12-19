package common

import (
	"errors"
	"net/http"
	"simple-auth/pkg/appcontext"
	"simple-auth/pkg/saerrors"

	"github.com/labstack/echo/v4"
)

const (
	ErrDeserialization saerrors.ErrorCode = "deserialize-failed"
	ErrBadRequest      saerrors.ErrorCode = "bad-request"
	ErrInternal        saerrors.ErrorCode = "internal-error"
	ErrMissingFields   saerrors.ErrorCode = "missing-fields"
)

type ErrorResponse struct {
	Error   bool   `json:"error" default:"true"`
	Message string `json:"message" example:"A human-readable message"`
	Reason  string `json:"reason" example:"machine-code"`
}

type OKResponse struct {
	Success bool `json:"success" example:"true"`
}

func HttpOK(c echo.Context) error {
	return c.JSON(http.StatusOK, OKResponse{
		Success: true,
	})
}

func HttpBadRequest(c echo.Context, err error) error {
	return HttpError(c, http.StatusBadRequest, ErrBadRequest.Compose(err))
}

func HttpInternalError(c echo.Context, err error) error {
	return HttpError(c, http.StatusInternalServerError, ErrInternal.Wrap(err))
}

func HttpInternalErrorf(c echo.Context, err string) error {
	return HttpError(c, http.StatusInternalServerError, ErrInternal.Newf(err))
}

func HttpError(c echo.Context, code int, err error) error {
	var saerr saerrors.CodedError
	if errors.As(err, &saerr) {
		return httpErrorCoded(c, code, string(saerr.Code()), saerr.Message(), saerr.Error())
	}
	return httpErrorCoded(c, code, "no-code", err.Error(), err.Error())
}

func httpErrorCoded(c echo.Context, code int, reason, message, fullError string) error {
	log := appcontext.GetLogger(c)
	log.Warnf("%d [%s]: %s", code, reason, fullError)
	return c.JSON(code, ErrorResponse{
		Error:   true,
		Message: message,
		Reason:  reason,
	})
}
