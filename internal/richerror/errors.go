package richerror

import "fmt"

type ErrorCode string

const (
	ErrCodeNotFound     ErrorCode = "NOT_FOUND"
	ErrCodeInvalidInput ErrorCode = "INVALID_INPUT"
	ErrCodeInternal     ErrorCode = "INTERNAL"
)

type AppError struct {
	Code    ErrorCode
	Message string
	Err     error
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// Factory helpers:
func NotFound(msg string, err error) *AppError {
	return &AppError{Code: ErrCodeNotFound, Message: msg, Err: err}
}

func InvalidInput(msg string, err error) *AppError {
	return &AppError{Code: ErrCodeInvalidInput, Message: msg, Err: err}
}

func Internal(msg string, err error) *AppError {
	return &AppError{Code: ErrCodeInternal, Message: msg, Err: err}
}
