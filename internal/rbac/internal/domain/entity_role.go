package domain

import "errors"

var ErrRoleNotCreated = errors.New("role not created")

type Role struct {
	ID          uint64
	Name        string
	Description string
}

func (r *Role) ScanColumn() []any {
	return []any{&r.ID, &r.Name, &r.Description}
}
