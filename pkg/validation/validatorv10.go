// Package validation provides interfaces and implementations for data validation.
package validation

import "github.com/go-playground/validator/v10"

// V10Validator is a concrete implementation of the Validator interface
// using the go-playground/validator library.
type V10Validator struct {
	validate *validator.Validate
}

// NewV10Validator creates a new instance of V10Validator with the default
// validation settings provided by the go-playground/validator library.
func NewV10Validator() *V10Validator {
	return &V10Validator{
		validate: validator.New(),
	}
}

// Validate checks the given data against validation rules defined in its struct tags.
// It returns an error if validation fails or nil if the data is valid.
// The input data should typically be a struct type.
func (v *V10Validator) Validate(data any) error {
	return v.validate.Struct(data)
}
