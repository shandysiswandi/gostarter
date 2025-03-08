package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserVerification_Table(t *testing.T) {
	tests := []struct {
		name string
		uv   UserVerification
		want string
	}{
		{
			name: "Success",
			uv:   UserVerification{},
			want: "user_verifications",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.uv.Table()
			assert.Equal(t, tt.want, got)
		})
	}
}
