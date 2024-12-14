package domain

import "errors"

var ErrAccountNotCreated = errors.New("account not created")

type Account struct {
	ID     uint64
	UserID uint64
}
