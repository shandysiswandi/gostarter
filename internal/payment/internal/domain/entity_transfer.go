package domain

import "github.com/shopspring/decimal"

type Transfer struct {
	ID            uint64
	TransactionID uint64
	SenderID      uint64
	RecipientID   uint64
	Amount        decimal.Decimal
}

func (t *Transfer) ScanColumn() []any {
	return []any{&t.ID, &t.TransactionID, &t.SenderID, &t.RecipientID, &t.Amount}
}
