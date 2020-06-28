package common

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
