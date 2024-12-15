package domain

import "github.com/shopspring/decimal"

type Topup struct {
	ID            uint64
	TransactionID uint64
	ReferenceID   string
	Amount        decimal.Decimal
}

func (t *Topup) ScanColumn() []any {
	return []any{&t.ID, &t.TransactionID, &t.ReferenceID, &t.Amount}
}
