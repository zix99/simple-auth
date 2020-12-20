package main

import (
	"simple-auth/pkg/routes/common"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type validateWrapper struct {
	cv *validator.Validate
}

func NewGoPlaygroundValidator() echo.Validator {
	return &validateWrapper{
		validator.New(),
	}
}

func (s *validateWrapper) Validate(i interface{}) error {
	if err := s.cv.Struct(i); err != nil {
		return common.ErrBadRequest.Compose(err)
	}
	return nil
}
