package inbound

type User struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type (
	UpdateRequest struct {
		Name string `json:"name"`
	}

	UpdateResponse struct {
		Name string `json:"name"`
	}
)

type (
	LogoutResponse struct {
		Message string `json:"message"`
	}
)
