package domain

import (
	"errors"
	"time"
)

var ErrPasswordResetNotCreated = errors.New("password reset not created")

type PasswordReset struct {
	ID        uint64
	UserID    uint64
	Token     string
	ExpiresAt time.Time
}
