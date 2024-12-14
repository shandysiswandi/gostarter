package inbound

type (
	LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	LoginResponse struct {
		AccessToken      string `json:"access_token"`
		RefreshToken     string `json:"refresh_token"`
		AccessExpiresIn  int64  `json:"access_expires_in"`  // in seconds
		RefreshExpiresIn int64  `json:"refresh_expires_in"` // in seconds
	}
)

type (
	RegisterRequest struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	RegisterResponse struct {
		Email string `json:"email"`
	}
)

type (
	RefreshTokenRequest struct {
		RefreshToken string `json:"refresh_token"`
	}

	RefreshTokenResponse struct {
		AccessToken      string `json:"access_token"`
		RefreshToken     string `json:"refresh_token"`
		AccessExpiresIn  int64  `json:"access_expires_in"`  // in seconds
		RefreshExpiresIn int64  `json:"refresh_expires_in"` // in seconds
	}
)

type (
	ForgotPasswordRequest struct {
		Email string `json:"email"`
	}

	ForgotPasswordResponse struct {
		Email   string `json:"email"`
		Message string `json:"message"`
	}
)

type (
	ResetPasswordRequest struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}

	ResetPasswordResponse struct {
		Message string `json:"message"`
	}
)
