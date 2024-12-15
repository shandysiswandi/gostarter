package domain

import "github.com/shopspring/decimal"

type Bill struct {
	ID            uint64
	TransactionID uint64
	ReferenceID   string
	Type          BillType
	Amount        decimal.Decimal
}

func (b *Bill) ScanColumn() []any {
	return []any{&b.ID, &b.TransactionID, &b.ReferenceID, &b.Type, &b.Amount}
}
