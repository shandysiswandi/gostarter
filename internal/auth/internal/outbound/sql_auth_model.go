package outbound

import (
	"time"

	"github.com/shandysiswandi/gostarter/internal/auth/internal/domain"
)

type user struct {
	id       uint64
	email    string
	password string
}

func (u *user) ScanColumn() []any {
	return []any{&u.id, &u.email, &u.password}
}

func (u *user) toEntity() *domain.User {
	if u == nil {
		return nil
	}

	return &domain.User{
		ID:       u.id,
		Email:    u.email,
		Password: u.password,
	}
}

type token struct {
	id               uint64
	userID           uint64
	accessToken      string
	refreshToken     string
	accessExpiredAt  time.Time
	refreshExpiredAt time.Time
}

func (t *token) ScanColumn() []any {
	return []any{
		&t.id,
		&t.userID,
		&t.accessToken,
		&t.refreshToken,
		&t.accessExpiredAt,
		&t.refreshExpiredAt,
	}
}

func (t *token) toEntity() *domain.Token {
	if t == nil {
		return nil
	}

	return &domain.Token{
		ID:               t.id,
		UserID:           t.userID,
		AccessToken:      t.accessToken,
		RefreshToken:     t.refreshToken,
		AccessExpiredAt:  t.accessExpiredAt,
		RefreshExpiredAt: t.refreshExpiredAt,
	}
}

type passwordReset struct {
	id        uint64
	userID    uint64
	token     string
	expiresAt time.Time
}

func (p *passwordReset) ScanColumn() []any {
	return []any{&p.id, &p.userID, &p.token, &p.expiresAt}
}

func (p *passwordReset) toEntity() *domain.PasswordReset {
	if p == nil {
		return nil
	}

	return &domain.PasswordReset{
		ID:        p.id,
		UserID:    p.userID,
		Token:     p.token,
		ExpiresAt: p.expiresAt,
	}
}
