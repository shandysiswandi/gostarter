package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransfer_ScanColumn(t *testing.T) {
	tests := []struct {
		name   string
		entity *Transfer
		want   func(a *Transfer) []any
	}{
		{
			name:   "Success",
			entity: &Transfer{},
			want: func(entity *Transfer) []any {
				return []any{
					&entity.ID,
					&entity.TransactionID,
					&entity.SenderID,
					&entity.RecipientID,
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
