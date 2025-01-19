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

type Pagination struct {
	NextCursor string `json:"next_cursor"`
	HasMore    bool   `json:"has_more"`
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
	UpdateRoleRequest struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
)

type (
	FetchRoleRequest struct {
		ID          string `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
	}

	FetchRoleResponse struct {
		Roles      []Role     `json:"roles"`
		Pagination Pagination `json:"pagination"`
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

type (
	UpdatePermissionRequest struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
)

type (
	FetchPermissionRequest struct {
		ID          string `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
	}

	FetchPermissionResponse struct {
		Permissions []Permission `json:"permissions"`
		Pagination  Pagination   `json:"pagination"`
	}
)
