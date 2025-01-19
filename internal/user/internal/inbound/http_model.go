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
)

type (
	UpdatePasswordRequest struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
)

type (
	LogoutResponse struct {
		Message string `json:"message"`
	}
)
