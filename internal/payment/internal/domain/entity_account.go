package domain

import (
	"errors"

	"github.com/shopspring/decimal"
)

var ErrAccountNoRowsAffected = errors.New("account not created or update")

type Account struct {
	ID       uint64
	UserID   uint64
	Balanace decimal.Decimal
}

func (a *Account) ScanColumn() []any {
	return []any{
		&a.ID,
		&a.UserID,
		&a.Balanace,
	}
}
