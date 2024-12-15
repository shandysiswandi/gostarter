package domain

import "github.com/shopspring/decimal"

type Account struct {
	ID       uint64
	UserID   uint64
	Balanace decimal.Decimal
}

func (a *Account) ScanColumn() []any {
	return []any{&a.ID, &a.UserID, &a.Balanace}
}
