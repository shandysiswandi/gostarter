package inbound

type (
	CreateRoleRequest struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	CreateRoleResponse struct {
		ID uint64 `json:"id,string"`
	}
)
