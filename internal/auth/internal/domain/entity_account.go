package domain

import "errors"

var ErrAccountNotCreated = errors.New("account not created")

type Account struct {
	ID     uint64 `db:"id"`
	UserID uint64 `db:"user_id"`
}

func (a Account) Table() string {
	return "accounts"
}
