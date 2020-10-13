package saerrors

import (
	"errors"
	"fmt"
)

type ErrorCode string

const Undefined ErrorCode = "undefined"

type codedError struct {
	cause   error
	message string
	code    ErrorCode
}

type CodedError interface {
	error
	Unwrap() error
	Code() ErrorCode
	Message() string
}

func (s *codedError) Error() string {
	if s.cause != nil {
		return s.Message() + ": " + s.cause.Error()
	}
	return s.cause.Error()
}

func (s *codedError) Unwrap() error {
	return s.cause
}

// Code returns the error-code
func (s *codedError) Code() ErrorCode {
	return s.code
}

// Message is the raw message without concatenating the error
func (s *codedError) Message() string {
	if s.message != "" {
		return s.message
	}
	return string(s.code)
}

func (s ErrorCode) New() error {
	return &codedError{
		code: s,
	}
}

func (s ErrorCode) Newf(msg string, args ...interface{}) error {
	return &codedError{
		message: fmt.Sprintf(msg, args...),
		code:    s,
	}
}

// Wrap an error with a code
func (s ErrorCode) Wrap(err error) error {
	return &codedError{
		cause: err,
		code:  s,
	}
}

// Compose an error (rather than wrap it, become that error)
func (s ErrorCode) Compose(err error) error {
	return &codedError{
		code:    s,
		message: err.Error(),
	}
}

// Wrapf takes an error, and wraps it with an additional message
func (s ErrorCode) Wrapf(err error, msg string, args ...interface{}) error {
	return &codedError{
		message: fmt.Sprintf(msg, args...),
		cause:   err,
		code:    s,
	}
}

// UnwrapCode finds the closet error in the wrapped error that has a code
func UnwrapCode(err error) ErrorCode {
	var codedErr CodedError
	if errors.As(err, &codedErr) {
		return codedErr.Code()
	}
	return Undefined
}
