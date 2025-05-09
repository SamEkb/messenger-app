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

	ErrTokenExpired = errors.New("token expired")
	ErrInvalidToken = errors.New("invalid token")

	ErrServiceUnavailable = errors.New("service unavailable")
)

const (
	CodeNotFound      = "NOT_FOUND"
	CodeAlreadyExists = "ALREADY_EXISTS"
	CodeUnauthorized  = "UNAUTHORIZED"
	CodeForbidden     = "FORBIDDEN"
	CodeInvalidInput  = "INVALID_INPUT"
	CodeValidation    = "VALIDATION"
	CodeInternal      = "INTERNAL"
	CodeDatabase      = "DATABASE"
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

func As(err error, target interface{}) bool {
	return errors.As(err, &target)
}

func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}

func New(message string) error {
	return errors.New(message)
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

func NewDatabaseError(err error, format string, args ...interface{}) *AppError {
	return &AppError{
		Err:     err,
		Message: fmt.Sprintf(format, args...),
		Code:    CodeDatabase,
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

func GetErrorCode(err error) string {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code
	}
	return "UNKNOWN"
}

func GetErrorDetails(err error) map[string]interface{} {
	var appErr *AppError
	if errors.As(err, &appErr) && appErr.Details != nil {
		return appErr.Details
	}
	return nil
}
