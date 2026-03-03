package models

type AuthResponse struct {
	AccessToken 	string `json:"access_token"`
	RefreshToken 	string `json:"refresh_token"`
	User 			User   `json:"user"`
}

type RegisterRequest struct {
	Email string `json:"email"`
	Password string `json:"password,omitempty"`
	Name  string `json:"name"`
}

type LoginRequest struct {
	Email string `json:"email"`
	Password string `json:"password,omitempty"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}