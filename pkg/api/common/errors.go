package common

import "fmt"

type ErrorResponse struct {
	Fatal   bool   `json:"fatal"`
	Message string `json:"message"`
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
