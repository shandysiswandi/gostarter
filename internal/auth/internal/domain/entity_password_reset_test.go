package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasswordReset_ScanColumn(t *testing.T) {
	tests := []struct {
		name string
		pr   *PasswordReset
		want func(pr *PasswordReset) []any
	}{
		{
			name: "Success",
			pr:   &PasswordReset{},
			want: func(pr *PasswordReset) []any {
				return []any{
					&pr.ID,
					&pr.UserID,
					&pr.Token,
					&pr.ExpiresAt,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.pr.ScanColumn()
			assert.Equal(t, tt.want(tt.pr), got)
		})
	}
}
