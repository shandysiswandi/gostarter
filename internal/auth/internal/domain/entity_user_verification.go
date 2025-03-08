package domain

import (
	"database/sql"
	"errors"
	"time"
)

var ErrUserVerificationNotCreated = errors.New("user verification not created")

type UserVerification struct {
	ID         uint64              `db:"id"`
	UserID     uint64              `db:"user_id"`
	Code       string              `db:"code"`
	ExpiresAt  sql.Null[time.Time] `db:"expires_at"`
	VerifiedAt sql.Null[time.Time] `db:"verified_at"`
}

func (UserVerification) Table() string {
	return "user_verifications"
}
