package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/genproto/googleapis/type/date"
)

func TestNewProtoValidator(t *testing.T) {
	tests := []struct {
		name    string
		wantErr error
	}{
		{
			name:    "Success",
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := NewProtoValidator()
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				assert.NotNil(t, got)
			}
		})
	}
}

func TestProtoValidator_Validate(t *testing.T) {
	tests := []struct {
		name    string
		data    any
		wantErr error
	}{
		{
			name:    "ErrorNotProtoMessage",
			data:    string("error"),
			wantErr: ErrNotProtoMessage,
		},
		{
			name:    "Success",
			data:    &date.Date{},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			v, _ := NewProtoValidator()
			err := v.Validate(tt.data)

			if tt.wantErr != nil {
				assert.Equal(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
