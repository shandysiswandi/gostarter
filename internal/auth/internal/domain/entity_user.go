package domain

import "errors"

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrUserNotCreated = errors.New("user not created")
	ErrUserNotUpdated = errors.New("user not updated")
	ErrUserNotDeleted = errors.New("user not deleted")
)

type User struct {
	ID       uint64
	Email    string
	Password string
}

func (u *User) ScanColumn() []any {
	return []any{&u.ID, &u.Email, &u.Password}
}
