package validation

import "github.com/go-playground/validator/v10"

// V10Validator is a concrete implementation of the Validator interface.
type V10Validator struct {
	validate *validator.Validate
}

// NewV10Validator creates a new instance of V10Validator.
func NewV10Validator() *V10Validator {
	return &V10Validator{
		validate: validator.New(),
	}
}

// Validate checks the given data against validation rules defined in its struct tags.
func (v *V10Validator) Validate(data any) error {
	return v.validate.Struct(data)
}
