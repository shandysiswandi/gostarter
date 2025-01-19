package domain

import "errors"

var ErrRolePermissionNotCreated = errors.New("role permission not created")

type RolePermission struct {
	RoleID       uint64
	PermissionID uint64
}

func (rp *RolePermission) ScanColumn() []any {
	return []any{&rp.RoleID, &rp.PermissionID}
}
