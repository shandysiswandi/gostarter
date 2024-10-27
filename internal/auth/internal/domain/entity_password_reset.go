package domain

import (
	"errors"
	"time"
)

var (
	ErrPasswordResetNotFound   = errors.New("password reset not found")
	ErrPasswordResetNotCreated = errors.New("password reset not created")
)

type PasswordReset struct {
	ID        uint64
	UserID    uint64
	Token     string
	ExpiresAt time.Time
}

func (p *PasswordReset) ScanColumn() []any {
	return []any{&p.ID, &p.UserID, &p.Token, &p.ExpiresAt}
}
