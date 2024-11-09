package domain

import "errors"

var ErrUserNotCreated = errors.New("user not created")

type User struct {
	ID       uint64
	Email    string
	Password string
}
