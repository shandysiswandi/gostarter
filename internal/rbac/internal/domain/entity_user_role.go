package domain

import "errors"

var ErrUserRoleNotCreated = errors.New("user role not created")

type UserRole struct {
	UserID uint64
	RoleID uint64
}

func (ur *UserRole) ScanColumn() []any {
	return []any{&ur.RoleID, &ur.UserID}
}
