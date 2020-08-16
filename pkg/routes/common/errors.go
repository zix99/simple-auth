package common

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type ErrorResponse struct {
	Fatal   bool   `json:"fatal"`
	Message string `json:"message"`
}

func HttpError(c echo.Context, code int, err error) error {
	return HttpErrorf(c, code, err.Error())
}

func HttpErrorf(c echo.Context, code int, err string, args ...interface{}) error {
	msg := fmt.Sprintf(err, args...)
	logrus.Warn(msg)
	return c.JSON(code, JsonErrorf(msg))
}

func JsonError(err error) ErrorResponse {
	return JsonErrorf(err.Error())
}

func JsonErrorf(s string, args ...interface{}) ErrorResponse {
	return ErrorResponse{
		Fatal:   true,
		Message: fmt.Sprintf(s, args...),
	}
}
