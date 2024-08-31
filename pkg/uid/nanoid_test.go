package uid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewNanoIDString(t *testing.T) {
	tests := []struct {
		name    string
		length  int
		wantErr bool
	}{
		{
			name:    "ValidLength",
			length:  21,
			wantErr: false,
		},
		{
			name:    "ValidLengthEdge",
			length:  1,
			wantErr: false,
		},
		{
			name:    "ValidLengthMax",
			length:  255,
			wantErr: false,
		},
		{
			name:    "InvalidLengthTooLow",
			length:  0,
			wantErr: true,
		},
		{
			name:    "InvalidLengthTooHigh",
			length:  256,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := NewNanoIDString(tt.length)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.length, got.length)
			}
		})
	}
}

func TestNanoIDString_Generate(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{
			name:   "GenerateLength21",
			length: 21,
		},
		{
			name:   "GenerateLength1",
			length: 1,
		},
		{
			name:   "GenerateLength255",
			length: 255,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			n, err := NewNanoIDString(tt.length)
			assert.NoError(t, err)
			got := n.Generate()
			assert.Equal(t, tt.length, len(got))
		})
	}
}
