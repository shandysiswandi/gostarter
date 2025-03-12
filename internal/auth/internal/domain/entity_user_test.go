package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser_Table(t *testing.T) {
	tests := []struct {
		name string
		u    User
		want string
	}{
		{
			name: "Success",
			u:    User{},
			want: "users",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.u.Table()
			assert.Equal(t, tt.want, got)
		})
	}
}
