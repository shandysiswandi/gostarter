package domain

import (
	"errors"
	"time"
)

var (
	ErrTokenNotFound       = errors.New("token not found")
	ErrTokenNoRowsAffected = errors.New("token not created or update")
)

type Token struct {
	ID               uint64
	UserID           uint64
	AccessToken      string
	RefreshToken     string
	AccessExpiredAt  time.Time
	RefreshExpiredAt time.Time
}

func (t *Token) ScanColumn() []any {
	return []any{&t.ID, &t.UserID, &t.AccessToken, &t.RefreshToken, &t.AccessExpiredAt, &t.RefreshExpiredAt}
}
