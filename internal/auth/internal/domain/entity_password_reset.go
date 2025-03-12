package domain

import (
	"errors"
	"time"
)

var ErrPasswordResetNotCreated = errors.New("password reset not created")

type PasswordReset struct {
	ID        uint64    `db:"id"`
	UserID    uint64    `db:"user_id"`
	Token     string    `db:"token"`
	ExpiresAt time.Time `db:"expires_at"`
}

func (PasswordReset) Table() string {
	return "password_resets"
}
