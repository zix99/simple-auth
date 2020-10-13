package common

import (
	"errors"
	"net/http"
	"simple-auth/pkg/routes/middleware"
	"simple-auth/pkg/saerrors"

	"github.com/labstack/echo/v4"
)

const (
	ErrDeserialization saerrors.ErrorCode = "deserialize-failed"
	ErrBadRequest      saerrors.ErrorCode = "bad-request"
	ErrInternal        saerrors.ErrorCode = "internal-error"
)

type ErrorResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Reason  string `json:"reason"`
}

func HttpOK(c echo.Context) error {
	return c.JSON(http.StatusOK, Json{
		"success": true,
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
		return httpErrorCoded(c, code, string(saerr.Code()), saerr.Message())
	}
	return httpErrorCoded(c, code, "no-code", err.Error())
}

func httpErrorCoded(c echo.Context, code int, reason, err string) error {
	log := middleware.GetLogger(c)
	log.Warnf("%d [%s]: %s", code, reason, err)
	return c.JSON(code, ErrorResponse{
		Error:   true,
		Message: err,
		Reason:  reason,
	})
}
