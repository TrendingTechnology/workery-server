package models

// The struct used to represent the user's `register` POST request data.
type RegisterRequest struct {
	TenantId uint64     `json:"tenant_id,string"`
	FirstName string     `json:"first_name"`
	LastName string     `json:"last_name"`
	Email string        `json:"email"`
	Password string     `json:"password"`
}

// The struct used to represent the system's response when the `register` POST request was a success.
type RegisterResponse struct {
	Message string `json:"message"`
}

// The struct used to represent the user's `login` POST request data.
type LoginRequest struct {
	Email string    `json:"email"`
	Password string `json:"password"`
}

// The struct used to represent the system's response when the `login` POST request was a success.
type LoginResponse struct {
	AccessToken string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
}

// The struct used to represent the user's `refresh token` POST request data.
type RefreshTokenRequest struct {
	Value string     `json:"value"`
}

// The struct used to represent the system's response when the `refresh token` POST request was a success.
type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
}
