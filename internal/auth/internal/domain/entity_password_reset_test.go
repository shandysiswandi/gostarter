package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasswordReset_Table(t *testing.T) {
	tests := []struct {
		name string
		pr   PasswordReset
		want string
	}{
		{
			name: "Success",
			pr:   PasswordReset{},
			want: "password_resets",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.pr.Table()
			assert.Equal(t, tt.want, got)
		})
	}
}
