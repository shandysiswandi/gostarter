package domain

import (
	"testing"

	"github.com/shandysiswandi/gostarter/pkg/enum"
	"github.com/stretchr/testify/assert"
)

func TestBillType_Values(t *testing.T) {
	tests := []struct {
		name string
		want map[enum.Enumerate]string
	}{
		{
			name: "Success",
			want: map[enum.Enumerate]string{
				BillTypeUnknown:  "UNKNOWN",
				BillTypePulsa:    "PULSA",
				BillTypeListrik:  "LISTRIK",
				BillTypeInternet: "INTERNET",
				BillTypeDonasi:   "DONASI",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, BillType(0).Values())
		})
	}
}

func TestBill_ScanColumn(t *testing.T) {
	tests := []struct {
		name   string
		entity *Bill
		want   func(a *Bill) []any
	}{
		{
			name:   "Success",
			entity: &Bill{},
			want: func(entity *Bill) []any {
				return []any{
					&entity.ID,
					&entity.TransactionID,
					&entity.ReferenceID,
					&entity.Type,
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
