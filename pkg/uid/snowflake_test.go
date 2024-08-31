package uid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSnowflakeNumber(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "CreateNewSnowflakeUIDNumber",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := NewSnowflakeNumber()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				// Optionally, you might want to verify that the node ID is within expected range
				assert.True(t, got.node != nil)
			}
		})
	}
}

func TestSnowflakeNumber_Generate(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "GenerateUniqueID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s, err := NewSnowflakeNumber()
			assert.NoError(t, err)
			got := s.Generate()
			assert.True(t, got > 0, "Generated ID should be greater than 0")
		})
	}
}
