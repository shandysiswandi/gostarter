// Package errs provides a structured way to define, create, and handle errors in a Go application.
// It categorizes errors by type (e.g., validation, business, server) and code (e.g., specific error codes),
// facilitating easier error handling and debugging.
package gerr

import (
	"net/http"
)

// Type represents the domain or category of an error.
type Type int

// Constants for different error types.
const (
	TypeValidation Type = iota // Validation errors (e.g., input validation failures).
	TypeBusiness               // Business logic errors (e.g., domain rule violations).
	TypeServer                 // Server-side errors (e.g., database or network issues).
)

// Code represents a specific error code for identifying the exact error.
type Code int

// Constants for common error codes. Extend these based on application needs.
const (
	CodeUnknown         Code = iota // Unknown or unspecified error.
	CodeInvalidFormat               // Error code for invalid format.
	CodeInvalidInput                // Error code for invalid input.
	CodeNotFound                    // Error code for resource not found.
	CodeConflict                    // Error code for conflict situations (e.g., duplicate entries).
	CodeUnauthorized                // Error code for unauthorized access.
	CodeForbidden                   // Error code for forbidden actions.
	CodeContentTooLarge             // Error code for content too large.
	CodeTimeout                     // Error code for operation timeout.
	CodeInternal                    // Error code for internal server errors.
)

// Error is a custom error type that wraps an error with additional context, including type and code.
type Error struct {
	err     error  // The underlying error, if any.
	msg     string // A custom error message providing additional context.
	errType Type   // The type of the error (validation, business, server).
	code    Code   // The specific code of the error.
}

// Error implements the error interface and returns the error message as a string.
func (e *Error) Error() string {
	if e.err != nil {
		return e.err.Error()
	}

	if e.msg != "" {
		return e.msg
	}

	switch e.errType {
	case TypeBusiness:
		return "business error"
	case TypeServer:
		return "server error"
	case TypeValidation:
		return "validation error"
	default:
		return "unknown error"
	}
}

// Unwrap returns the underlying error, if any, allowing for error unwrapping.
func (e *Error) Unwrap() error {
	return e.err
}

// StatusCode maps the error type and code to an appropriate HTTP status code.
func (e *Error) StatusCode() int {
	switch e.code {
	case CodeInvalidFormat:
		return http.StatusBadRequest
	case CodeInvalidInput:
		return http.StatusUnprocessableEntity
	case CodeNotFound:
		return http.StatusNotFound
	case CodeUnauthorized:
		return http.StatusUnauthorized
	case CodeForbidden:
		return http.StatusForbidden
	case CodeTimeout:
		return http.StatusRequestTimeout
	case CodeConflict:
		return http.StatusConflict
	case CodeContentTooLarge:
		return http.StatusRequestEntityTooLarge
	case CodeInternal, CodeUnknown:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// NewBadFormat creates a new validation error for bad format with a custom message and an underlying error.
func NewBadFormat(msg string, err error) error {
	return &Error{
		err:     err,
		msg:     msg,
		errType: TypeValidation,
		code:    CodeInvalidFormat,
	}
}

// NewValidation creates a new validation error for invalid input with a custom message and an underlying error.
func NewValidation(msg string, err error) error {
	return &Error{
		err:     err,
		msg:     msg,
		errType: TypeValidation,
		code:    CodeInvalidInput,
	}
}

// NewServer creates a new server error with a custom message and an underlying error.
func NewServer(msg string, err error) error {
	return &Error{
		err:     err,
		msg:     msg,
		errType: TypeServer,
		code:    CodeInternal,
	}
}

// NewBusiness creates a new business logic error with a custom message and a specific error code.
func NewBusiness(msg string, code Code) error {
	return &Error{
		msg:     msg,
		errType: TypeBusiness,
		code:    code,
	}
}

// New creates a new custom error with the provided underlying error, message, type, and code.
//
// The function wraps an existing error (`err`) with additional context (`msg`), categorizing the error
// by type (`et`) and assigning a specific code (`code`). This allows for structured error handling
// that can be easily categorized and processed based on the error's type and code.
//
// Parameters:
//   - err: The underlying error that triggered this error. This can be `nil` if there is no underlying error.
//   - msg: A custom message providing additional context about the error. This can be an empty string.
//   - et: The type of the error, used to categorize the error (e.g., validation, business, server).
//   - code: The specific code of the error, used to identify the exact nature of the error.
//
// Returns:
//   - error: A pointer to an `Error` struct containing all the provided information.
//
// Example:
//
//	err := New(nil, "input validation failed", TypeValidation, CodeInvalidInput)
//	fmt.Println(err.Error()) // Output: input validation failed
func New(err error, msg string, et Type, code Code) error {
	return &Error{
		err:     err,
		msg:     msg,
		errType: et,
		code:    code,
	}
}
