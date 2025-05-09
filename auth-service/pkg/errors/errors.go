package errors

import (
	"errors"
	"fmt"
)

// Стандартные ошибки
var (
	// Ошибки аутентификации
	ErrNotFound      = errors.New("resource not found")
	ErrAlreadyExists = errors.New("resource already exists")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrInvalidInput  = errors.New("invalid input")

	// Ошибки базы данных
	ErrDatabaseConnection = errors.New("database connection error")
	ErrDatabaseQuery      = errors.New("database query error")

	// Ошибки валидации
	ErrValidation = errors.New("validation error")

	// Системные ошибки
	ErrInternal = errors.New("internal error")
	ErrTimeout  = errors.New("operation timeout")
)

// AppError представляет структурированную ошибку приложения
type AppError struct {
	Err     error
	Message string
	Code    string
}

// Error реализует интерфейс error
func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap позволяет использовать errors.Is и errors.As
func (e *AppError) Unwrap() error {
	return e.Err
}

// Is - используется для проверки типа ошибки
// Позволяет сравнивать ошибки с помощью errors.Is
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As - используется для приведения типа ошибки
// Позволяет приводить ошибки с помощью errors.As
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// Функции для создания новых ошибок

// NewNotFoundError создает новую ошибку "не найдено"
func NewNotFoundError(format string, args ...interface{}) error {
	return &AppError{
		Err:     ErrNotFound,
		Message: fmt.Sprintf(format, args...),
		Code:    "NOT_FOUND",
	}
}

// NewAlreadyExistsError создает новую ошибку "уже существует"
func NewAlreadyExistsError(format string, args ...interface{}) error {
	return &AppError{
		Err:     ErrAlreadyExists,
		Message: fmt.Sprintf(format, args...),
		Code:    "ALREADY_EXISTS",
	}
}

// NewUnauthorizedError создает новую ошибку "не авторизован"
func NewUnauthorizedError(format string, args ...interface{}) error {
	return &AppError{
		Err:     ErrUnauthorized,
		Message: fmt.Sprintf(format, args...),
		Code:    "UNAUTHORIZED",
	}
}

// NewInvalidInputError создает новую ошибку "неверный ввод"
func NewInvalidInputError(format string, args ...interface{}) error {
	return &AppError{
		Err:     ErrInvalidInput,
		Message: fmt.Sprintf(format, args...),
		Code:    "INVALID_INPUT",
	}
}

// NewInternalError создает новую внутреннюю ошибку
func NewInternalError(err error, format string, args ...interface{}) error {
	return &AppError{
		Err:     err,
		Message: fmt.Sprintf(format, args...),
		Code:    "INTERNAL",
	}
}

// GetErrorCode извлекает код ошибки из ошибки приложения
func GetErrorCode(err error) string {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code
	}
	return "UNKNOWN"
}
