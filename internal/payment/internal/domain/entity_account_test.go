package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccount_ScanColumn(t *testing.T) {
	tests := []struct {
		name   string
		entity *Account
		want   func(a *Account) []any
	}{
		{
			name:   "Success",
			entity: &Account{},
			want: func(entity *Account) []any {
				return []any{
					&entity.ID,
					&entity.UserID,
					&entity.Balanace,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.entity.ScanColumn()
			assert.Equal(t, tt.want(tt.entity), got)
		})
	}
}
