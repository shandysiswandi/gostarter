package domain

import (
	"errors"
)

var ErrPermissionNotCreated = errors.New("permission not created")

type Permission struct {
	ID          uint64
	Name        string
	Description string
}

func (p *Permission) ScanColumn() []any {
	return []any{&p.ID, &p.Name, &p.Description}
}
