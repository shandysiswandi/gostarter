package uid

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewUUIDString(t *testing.T) {
	tests := []struct {
		name string
		want *UUIDString
	}{
		{
			name: "CreateNewUUIDString",
			want: &UUIDString{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewUUIDString()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUUIDString_Generate(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "GenerateValidUUID",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			u := NewUUIDString()
			got := u.Generate()

			// Validate UUID format and length
			parsedUUID, err := uuid.Parse(got)
			assert.NoError(t, err)
			assert.Equal(t, 36, len(got)) // Length of UUID string with hyphens
			assert.True(t, parsedUUID.Version() == 4, "UUID version should be 4")
		})
	}
}
