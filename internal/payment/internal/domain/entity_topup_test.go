package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTopup_ScanColumn(t *testing.T) {
	tests := []struct {
		name   string
		entity *Topup
		want   func(a *Topup) []any
	}{
		{
			name:   "Success",
			entity: &Topup{},
			want: func(entity *Topup) []any {
				return []any{
					&entity.ID,
					&entity.TransactionID,
					&entity.ReferenceID,
					&entity.Amount,
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
