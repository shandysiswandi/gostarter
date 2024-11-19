package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToken_ScanColumn(t *testing.T) {
	tests := []struct {
		name string
		to   *Token
		want func(to *Token) []any
	}{
		{
			name: "Success",
			to:   &Token{},
			want: func(to *Token) []any {
				return []any{
					&to.ID,
					&to.UserID,
					&to.AccessToken,
					&to.RefreshToken,
					&to.AccessExpiredAt,
					&to.RefreshExpiredAt,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.to.ScanColumn()
			assert.Equal(t, tt.want(tt.to), got)
		})
	}
}
