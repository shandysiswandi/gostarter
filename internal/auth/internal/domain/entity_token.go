package domain

import (
	"errors"
	"time"
)

var ErrTokenNoRowsAffected = errors.New("token not created or update")

type Token struct {
	ID               uint64    `db:"id"`
	UserID           uint64    `db:"user_id"`
	AccessToken      string    `db:"access_token"`
	RefreshToken     string    `db:"refresh_token"`
	AccessExpiresAt  time.Time `db:"access_expires_at"`
	RefreshExpiresAt time.Time `db:"refresh_expires_at"`
}

func (Token) Table() string {
	return "tokens"
}
