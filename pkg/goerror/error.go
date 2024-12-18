// Package goerror provides a structured way to define, create, and handle errors in Go applications.
// It integrates with both HTTP status codes and gRPC status codes, categorizing errors by type (validation,
// business, server) and code (specific error codes). This allows for easier error handling and debugging.
package goerror

import (
	"fmt"
	"net/http"

	"github.com/shandysiswandi/gostarter/pkg/goerror/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GoError is a custom error type that wraps an error with additional context, including type and code.
type GoError struct {
	err     error  // The underlying error, if any.
	msg     string // A custom error message providing additional context.
	errType Type   // The type of the error (validation, business, server).
	code    Code   // The specific code of the error.
}

// Error implements the error interface and returns the error message as a string.
func (e *GoError) Error() string {
	if e.err != nil {
		return e.err.Error()
	}

	if e.msg != "" {
		return e.msg
	}

	if e.errType == TypeValidation {
		return "Validation violation"
	}

	if e.errType == TypeBusiness {
		return "Logical business not meet with requirement"
	}

	if e.errType == TypeServer {
		return "Internal error"
	}

	return "Unknown error"
}

// String returns a formatted string with details about the error.
func (e *GoError) String() string {
	return fmt.Sprintf(
		"Error Type: %s, Code: %s, Message: %s, Underlying Error: %v",
		e.errType.String(),
		e.code.String(),
		e.msg,
		e.err,
	)
}

// Msg returns the custom error message.
func (e *GoError) Msg() string {
	return e.msg
}

// Type returns the error type (validation, business, or server).
func (e *GoError) Type() Type {
	return e.errType
}

// Code returns the specific code of the error.
func (e *GoError) Code() Code {
	return e.code
}

// Unwrap returns the underlying error, if any, allowing error chaining.
func (e *GoError) Unwrap() error {
	return e.err
}

// GRPCStatus maps the error to a gRPC status code.
//
// It returns a gRPC status object based on the error's code and type, facilitating integration with
// gRPC-based systems.
func (e *GoError) GRPCStatus() *status.Status {
	var sts *status.Status

	switch e.code {
	case CodeInvalidFormat:
		sts = status.New(codes.InvalidArgument, e.msg)
	case CodeInvalidInput:
		sts = status.New(codes.FailedPrecondition, e.msg)
	case CodeNotFound:
		sts = status.New(codes.NotFound, e.msg)
	case CodeUnauthorized:
		sts = status.New(codes.Unauthenticated, e.msg)
	case CodeForbidden:
		sts = status.New(codes.PermissionDenied, e.msg)
	case CodeTimeout:
		sts = status.New(codes.DeadlineExceeded, e.msg)
	case CodeConflict:
		sts = status.New(codes.AlreadyExists, e.msg)
	case CodeInternal:
		sts = status.New(codes.Internal, e.msg)
	case CodeUnknown:
		sts = status.New(codes.Unknown, e.msg)
	default:
		sts = status.New(codes.Unknown, e.msg)
	}

	ge := &pb.Error{
		Code:       e.code.String(),
		Type:       e.errType.String(),
		Message:    e.msg,
		Attributes: nil, // next-mr: will be added later
	}

	if swd, err := sts.WithDetails(ge); err == nil {
		return swd
	}

	return sts
}

// StatusCode maps the error to an appropriate HTTP status code.
func (e *GoError) StatusCode() int {
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
	case CodeInternal, CodeUnknown:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// New creates a new custom error with the provided error, message, type, and code.
func New(err error, msg string, et Type, code Code) error {
	return &GoError{
		err:     err,
		msg:     msg,
		errType: et,
		code:    code,
	}
}

// NewServerInternal creates a server-type error with the provided error.
func NewServerInternal(err error) error {
	return New(err, "internal server error", TypeServer, CodeInternal)
}

// NewBusiness creates a business-type error with the specified message and code.
func NewBusiness(msg string, code Code) error {
	return New(nil, msg, TypeBusiness, code)
}

// NewInvalidInput creates a validation error for invalid input with a message and underlying error.
func NewInvalidInput(msg string, err error) error {
	return New(err, msg, TypeValidation, CodeInvalidInput)
}

// NewInvalidFormat creates a validation error for invalid format with a message and underlying error.
func NewInvalidFormat(msg string) error {
	return New(nil, msg, TypeValidation, CodeInvalidFormat)
}
