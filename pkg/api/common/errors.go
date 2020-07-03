package common

import "fmt"

type ErrorResponse struct {
	Fatal   bool
	Message string
}

func JsonError(err error) ErrorResponse {
	return ErrorResponse{
		Fatal:   true,
		Message: err.Error(),
	}
}

func JsonErrorf(s string, args ...interface{}) ErrorResponse {
	return ErrorResponse{
		Fatal:   true,
		Message: fmt.Sprintf(s, args...),
	}
}
