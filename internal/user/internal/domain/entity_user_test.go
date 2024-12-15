package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser_ScanColumn(t *testing.T) {
	tests := []struct {
		name string
		u    *User
		want func(u *User) []any
	}{
		{
			name: "Success",
			u:    &User{},
			want: func(u *User) []any {
				return []any{
					&u.ID,
					&u.Name,
					&u.Email,
					&u.Password,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.u.ScanColumn()
			assert.Equal(t, tt.want(tt.u), got)
		})
	}
}
