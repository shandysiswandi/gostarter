package domain

import (
	"errors"
	"time"
)

var ErrTokenNoRowsAffected = errors.New("token not created or update")

type Token struct {
	ID               uint64
	UserID           uint64
	AccessToken      string
	RefreshToken     string
	AccessExpiredAt  time.Time
	RefreshExpiredAt time.Time
}
