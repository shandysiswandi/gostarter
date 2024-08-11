// Package errs provides a structured way to define, create, and handle errors in a Go application.
// It supports categorizing errors by type (e.g., validation, business, server) and code (e.g., specific error codes).
package errs

import (
	"errors"
	"fmt"
)

// Type represents the type of an error, categorizing it into different domains.
type Type int

// Different types of errors that can occur.
const (
	TypeValidation Type = iota // Validation errors (e.g., input validation failures)
	TypeBusiness               // Business logic errors (e.g., domain rules violations)
	TypeServer                 // Server-side errors (e.g., database or network issues)
)

// Code represents a specific error code that can be used to identify the exact error.
type Code int

// Common error codes. Extend these based on application needs.
const (
	CodeUnknown      Code = iota // Unknown or unspecified error
	CodeInvalidInput             // Specific code for invalid input
	CodeNotFound                 // Specific code for resource not found
	CodeConflict                 // Specific code for conflict situations (e.g., duplicate entries)
	CodeUnauthorized             // Specific code for unauthorized access
	CodeForbidden                // Specific code for forbidden actions
	CodeTimeout                  // Specific code for operation timeout
	CodeInternal                 // Specific code for internal server errors
)

// Error struct wraps an error with additional context, including type and code.
type Error struct {
	err     error  // The underlying error
	msg     string // Custom error message
	errType Type   // The type of the error (validation, business, server)
	code    Code   // The specific code of the error
}

// Error returns the error message as a string, implementing the error interface.
func (e *Error) Error() string {
	if e.err != nil {
		return e.err.Error()
	}

	if e.msg != "" {
		return e.msg
	}

	return ""
}

// Unwrap returns the underlying error, if any, allowing error unwrapping.
func (e *Error) Unwrap() error {
	return e.err
}

// Code returns the specific error code, allowing further categorization.
func (e *Error) Code() Code {
	return e.code
}

// Type returns the type of the error, indicating its category.
func (e *Error) Type() Type {
	return e.errType
}

// IsValidationError checks if an error is of TypeValidation.
func IsValidationError(err error) bool {
	var e *Error
	if errors.As(err, &e) {
		return e.errType == TypeValidation
	}

	return false
}

// IsBusinessError checks if an error is of TypeBusiness.
func IsBusinessError(err error) bool {
	var e *Error
	if errors.As(err, &e) {
		return e.errType == TypeBusiness
	}

	return false
}

// IsServerError checks if an error is of TypeServer.
func IsServerError(err error) bool {
	var e *Error
	if errors.As(err, &e) {
		return e.errType == TypeServer
	}

	return false
}

// NewValidation creates a new validation error with a custom message.
func NewValidation(msg string) error {
	return &Error{
		msg:     msg,
		errType: TypeValidation,
		code:    CodeInvalidInput,
	}
}

// NewBusiness creates a new business logic error with a custom message.
func NewBusiness(msg string) error {
	return &Error{
		msg:     msg,
		errType: TypeBusiness,
		code:    CodeNotFound,
	}
}

// NewServer creates a new server error with a custom message.
func NewServer(msg string) error {
	return &Error{
		msg:     msg,
		errType: TypeServer,
		code:    CodeInternal,
	}
}

// NewValidationFrom wraps an existing error as a validation error.
func NewValidationFrom(err error) error {
	return &Error{
		err:     err,
		errType: TypeValidation,
		code:    CodeInvalidInput,
	}
}

// NewBusinessFrom wraps an existing error as a business logic error.
func NewBusinessFrom(err error) error {
	return &Error{
		err:     err,
		errType: TypeBusiness,
		code:    CodeNotFound,
	}
}

// NewServerFrom wraps an existing error as a server error.
func NewServerFrom(err error) error {
	return &Error{
		err:     err,
		errType: TypeServer,
		code:    CodeInternal,
	}
}

// Wrap adds context to an existing error, preserving its type and code.
func Wrap(err error, msg string) error {
	var e *Error
	if errors.As(err, &e) {
		return &Error{
			err:     fmt.Errorf("%s: %w", msg, err),
			msg:     msg,
			errType: e.errType,
			code:    e.code,
		}
	}

	return &Error{
		err:     fmt.Errorf("%s: %w", msg, err),
		msg:     msg,
		errType: TypeServer,
		code:    CodeUnknown,
	}
}

// New creates a new custom error with the provided underlying error, message, type, and code.
//
// The function wraps an existing error (`err`) with additional context (`msg`), while also
// categorizing the error using a specific type (`et`) and assigning a specific code (`code`).
//
// This is useful for generating structured errors that can be easily categorized and handled
// based on their type and code.
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
