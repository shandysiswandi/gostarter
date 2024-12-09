package inbound

type (
	ProfileResponse struct {
		ID    uint64 `json:"id"`
		Email string `json:"email"`
	}
)

type (
	LogoutResponse struct {
		Message string `json:"message"`
	}
)
