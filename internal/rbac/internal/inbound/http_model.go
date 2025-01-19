package inbound

type Role struct {
	ID          uint64 `json:"id,string"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Permission struct {
	ID          uint64 `json:"id,string"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type (
	CreateRoleRequest struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	CreateRoleResponse struct {
		ID uint64 `json:"id,string"`
	}
)

type (
	CreatePermissionRequest struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	CreatePermissionResponse struct {
		ID uint64 `json:"id,string"`
	}
)
