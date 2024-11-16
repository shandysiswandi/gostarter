// Package validation provides interfaces and implementations for data validation.
package validation

import (
	"errors"
	"fmt"

	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/protobuf/proto"
)

var ErrNotProtoMessage = errors.New("provided data does not implement proto.Message")

// ProtoValidator is a concrete implementation of the Validator interface.
type ProtoValidator struct {
	validate *protovalidate.Validator
}

// NewProtoValidator creates a new instance of ProtoValidator.
func NewProtoValidator() (*ProtoValidator, error) {
	p, err := protovalidate.New()
	if err != nil {
		return nil, err
	}

	return &ProtoValidator{validate: p}, nil
}

// Validate checks the given data against validation rules defined in its
// Protocol Buffers struct tags. The data must implement the proto.Message
// interface, which is the interface for Protocol Buffers messages.
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
