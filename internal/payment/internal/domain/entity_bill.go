package domain

import (
	"github.com/shandysiswandi/gostarter/pkg/enum"
	"github.com/shopspring/decimal"
)

type BillType int

const (
	BillTypeUnknown BillType = iota
	BillTypePulsa
	BillTypeListrik
	BillTypeInternet
	BillTypeDonasi
)

func (bt BillType) Values() map[enum.Enumerate]string {
	return map[enum.Enumerate]string{
		BillTypeUnknown:  "UNKNOWN",
		BillTypePulsa:    "PULSA",
		BillTypeListrik:  "LISTRIK",
		BillTypeInternet: "INTERNET",
		BillTypeDonasi:   "DONASI",
	}
}

type Bill struct {
	ID            uint64
	TransactionID uint64
	ReferenceID   string
	Type          BillType
	Amount        decimal.Decimal
}

func (b *Bill) ScanColumn() []any {
	return []any{
		&b.ID,
		&b.TransactionID,
		&b.ReferenceID,
		&b.Type,
		&b.Amount,
	}
}
