package domain

import (
	"testing"

	"github.com/shandysiswandi/gostarter/pkg/enum"
	"github.com/stretchr/testify/assert"
)

func TestTransactionStatus_Values(t *testing.T) {
	tests := []struct {
		name string
		want map[enum.Enumerate]string
	}{
		{
			name: "Success",
			want: map[enum.Enumerate]string{
				TransactionStatusUnknown: "UNKNOWN",
				TransactionStatusPending: "PENDING",
				TransactionStatusFailed:  "FAILED",
				TransactionStatusSuccess: "SUCCESS",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, TransactionStatus(0).Values())
		})
	}
}

func TestTransactionType_Values(t *testing.T) {
	tests := []struct {
		name string
		want map[enum.Enumerate]string
	}{
		{
			name: "Success",
			want: map[enum.Enumerate]string{
				TransactionStatusUnknown: "UNKNOWN",
				TransactionTypeDebit:     "DEBIT",
				TransactionTypeCredit:    "CREDIT",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, TransactionType(0).Values())
		})
	}
}

func TestTransaction_ScanColumn(t *testing.T) {
	tests := []struct {
		name   string
		entity *Transaction
		want   func(a *Transaction) []any
	}{
		{
			name:   "Success",
			entity: &Transaction{},
			want: func(entity *Transaction) []any {
				return []any{
					&entity.ID,
					&entity.UserID,
					&entity.Amount,
					&entity.Type,
					&entity.Status,
					&entity.Remark,
					&entity.CreateAt,
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
