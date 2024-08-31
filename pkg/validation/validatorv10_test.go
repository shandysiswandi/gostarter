package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	Name  string `validate:"required"`
	Email string `validate:"required,email"`
	Age   int    `validate:"gte=0,lte=130"`
}

func TestNewV10Validator(t *testing.T) {
	v := NewV10Validator()
	assert.NotNil(t, v)
	assert.NotNil(t, v.validate)
	assert.IsType(t, &V10Validator{}, v)
}

func TestV10Validator_Validate(t *testing.T) {
	tests := []struct {
		name    string
		data    any
		wantErr bool
	}{
		{
			name: "ValidData",
			data: testStruct{
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   30,
			},
			wantErr: false,
		},
		{
			name: "MissingRequiredField",
			data: testStruct{
				Email: "john@example.com",
				Age:   30,
			},
			wantErr: true,
		},
		{
			name: "InvalidEmailFormat",
			data: testStruct{
				Name:  "John Doe",
				Email: "invalid-email",
				Age:   30,
			},
			wantErr: true,
		},
		{
			name: "InvalidAgeNegative",
			data: testStruct{
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   -5,
			},
			wantErr: true,
		},
		{
			name: "InvalidAgeTooHigh",
			data: testStruct{
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   150,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			v := NewV10Validator()
			err := v.Validate(tt.data)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
