package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type Transaction struct {
	ID       uint64
	UserID   uint64
	Amount   decimal.Decimal
	Type     TransactionType
	Status   TransactionStatus
	Remark   string
	CreateAt time.Time
}

func (t *Transaction) ScanColumn() []any {
	return []any{&t.ID, &t.UserID, &t.Amount, &t.Type, &t.Status, &t.Remark, &t.CreateAt}
}
