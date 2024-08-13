// Package validation provides interfaces and implementations for data validation.
package validation

import (
	"errors"
	"fmt"

	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/protobuf/proto"
)

var (
	ErrNotProtoMessage    = errors.New("provided data does not implement proto.Message")
	ErrInitProtoValidator = errors.New("failed to create protovalidate validator")
)

// ProtoValidator is a concrete implementation of the Validator interface
// that uses the protovalidate library to validate Protocol Buffers messages.
type ProtoValidator struct {
	validate *protovalidate.Validator
}

// NewProtoValidator creates a new instance of ProtoValidator with default
// validation settings provided by the protovalidate library.
//
// Returns:
// - *ProtoValidator: A new instance of ProtoValidator, ready for use.
// - error: An error if initialization fails, otherwise nil.
func NewProtoValidator() (*ProtoValidator, error) {
	p, err := protovalidate.New()
	if err != nil {
		return nil, ErrInitProtoValidator
	}

	return &ProtoValidator{
		validate: p,
	}, nil
}

// Validate checks the given data against validation rules defined in its
// Protocol Buffers struct tags. The data must implement the proto.Message
// interface, which is the interface for Protocol Buffers messages.
//
// Parameters:
//   - data: An interface{} type which is expected to be a Protocol Buffers message
//     (i.e., a type that implements proto.Message).
//
// Returns:
//   - error: An error if validation fails, or nil if the data is valid.
//     If the provided data does not implement proto.Message, an error is returned
//     indicating that the type assertion failed.
func (v *ProtoValidator) Validate(data any) error {
	msg, ok := data.(proto.Message)
	if !ok {
		return ErrNotProtoMessage
	}

	if err := v.validate.Validate(msg); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return nil
}
