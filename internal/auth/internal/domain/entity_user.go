package domain

import (
	"database/sql"
	"errors"
	"time"
)

var ErrUserNotCreated = errors.New("user not created")

type User struct {
	ID         uint64              `db:"id"`
	Name       string              `db:"name"`
	Email      string              `db:"email"`
	Password   string              `db:"password"`
	VerifiedAt sql.Null[time.Time] `db:"verified_at"`
}

func (u User) Table() string {
	return "users"
}
