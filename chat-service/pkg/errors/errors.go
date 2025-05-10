package errors

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound      = errors.New("resource not found")
	ErrAlreadyExists = errors.New("resource already exists")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrInvalidInput  = errors.New("invalid input")
	ErrForbidden     = errors.New("forbidden")

	ErrDatabaseConnection = errors.New("database connection error")
	ErrDatabaseQuery      = errors.New("database query error")

	ErrValidation = errors.New("validation error")

	ErrTimeout = errors.New("operation timeout")

	ErrTokenExpired = errors.New("token expired")
	ErrInvalidToken = errors.New("invalid token")
)

const (
	CodeNotFound      = "NOT_FOUND"
	CodeAlreadyExists = "ALREADY_EXISTS"
	CodeForbidden     = "FORBIDDEN"
	CodeUnauthorized  = "UNAUTHORIZED"
	CodeInvalidInput  = "INVALID_INPUT"
	CodeValidation    = "VALIDATION"
	CodeInternal      = "INTERNAL"
	CodeTimeout       = "TIMEOUT"
	CodeToken         = "TOKEN"
	CodeService       = "SERVICE"
)

type AppError struct {
	Err       error
	Message   string
	Code      string
	Details   map[string]interface{}
	Retriable bool
}

func (e *AppError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	if e.Err != nil {
		return e.Err.Error()
	}
	return "unknown error"
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) WithDetails(key string, value interface{}) *AppError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return e
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func NewNotFoundError(format string, args ...interface{}) *AppError {
	return &AppError{
		Err:     ErrNotFound,
		Message: fmt.Sprintf(format, args...),
		Code:    CodeNotFound,
	}
}

func NewAlreadyExistsError(format string, args ...interface{}) *AppError {
	return &AppError{
		Err:     ErrAlreadyExists,
		Message: fmt.Sprintf(format, args...),
		Code:    CodeAlreadyExists,
	}
}

func NewUnauthorizedError(format string, args ...interface{}) *AppError {
	return &AppError{
		Err:     ErrUnauthorized,
		Message: fmt.Sprintf(format, args...),
		Code:    CodeUnauthorized,
	}
}

func NewForbiddenError(format string, args ...interface{}) *AppError {
	return &AppError{
		Err:     ErrForbidden,
		Message: fmt.Sprintf(format, args...),
		Code:    CodeForbidden,
	}
}

func NewInvalidInputError(format string, args ...interface{}) *AppError {
	return &AppError{
		Err:     ErrInvalidInput,
		Message: fmt.Sprintf(format, args...),
		Code:    CodeInvalidInput,
	}
}

func NewValidationError(format string, args ...interface{}) *AppError {
	return &AppError{
		Err:     ErrValidation,
		Message: fmt.Sprintf(format, args...),
		Code:    CodeValidation,
	}
}

func NewInternalError(err error, format string, args ...interface{}) *AppError {
	return &AppError{
		Err:     err,
		Message: fmt.Sprintf(format, args...),
		Code:    CodeInternal,
	}
}

func NewTimeoutError(format string, args ...interface{}) *AppError {
	return &AppError{
		Err:     ErrTimeout,
		Message: fmt.Sprintf(format, args...),
		Code:    CodeTimeout,
	}
}

func NewTokenError(err error, format string, args ...interface{}) *AppError {
	return &AppError{
		Err:     err,
		Message: fmt.Sprintf(format, args...),
		Code:    CodeToken,
	}
}

func NewServiceError(err error, format string, args ...interface{}) *AppError {
	return &AppError{
		Err:     err,
		Message: fmt.Sprintf(format, args...),
		Code:    CodeService,
	}
}
