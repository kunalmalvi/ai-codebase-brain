package errors

import (
	"fmt"
	"strings"
)

type ErrorCode string

const (
	ErrCodeNotFound       ErrorCode = "NOT_FOUND"
	ErrCodeInvalidInput   ErrorCode = "INVALID_INPUT"
	ErrCodeParseError     ErrorCode = "PARSE_ERROR"
	ErrCodeStorageError   ErrorCode = "STORAGE_ERROR"
	ErrCodeIndexError     ErrorCode = "INDEX_ERROR"
	ErrCodeTimeout        ErrorCode = "TIMEOUT"
	ErrCodePermission     ErrorCode = "PERMISSION_DENIED"
	ErrCodeUnknown        ErrorCode = "UNKNOWN"
)

type AppError struct {
	Code    ErrorCode
	Message string
	Cause   error
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Cause
}

func New(code ErrorCode, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

func Wrap(code ErrorCode, message string, cause error) *AppError {
	return &AppError{Code: code, Message: message, Cause: cause}
}

func NotFound(resource string) *AppError {
	return New(ErrCodeNotFound, fmt.Sprintf("%s not found", resource))
}

func InvalidInput(field, reason string) *AppError {
	return New(ErrCodeInvalidInput, fmt.Sprintf("invalid %s: %s", field, reason))
}

func ParseError(file, reason string) *AppError {
	return New(ErrCodeParseError, fmt.Sprintf("failed to parse %s: %s", file, reason))
}

func StorageError(reason string) *AppError {
	return New(ErrCodeStorageError, fmt.Sprintf("storage error: %s", reason))
}

func IndexError(reason string) *AppError {
	return New(ErrCodeIndexError, fmt.Sprintf("index error: %s", reason))
}

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on %s: %s", e.Field, e.Message)
}

func ValidateRequired(field, value string) error {
	if strings.TrimSpace(value) == "" {
		return &ValidationError{Field: field, Message: "is required"}
	}
	return nil
}

func ValidatePositive(field string, value float64) error {
	if value <= 0 {
		return &ValidationError{Field: field, Message: "must be positive"}
	}
	return nil
}

func ValidateInRange(field string, value, min, max float64) error {
	if value < min || value > max {
		return &ValidationError{Field: field, Message: fmt.Sprintf("must be between %.0f and %.0f", min, max)}
	}
	return nil
}