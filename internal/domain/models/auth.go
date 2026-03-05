package models

type AuthResponse struct {
	AccessToken 	string `json:"access_token"`
	RefreshToken 	string `json:"refresh_token"`
	CSRFToken 		string `json:"-"`
	User 			User   `json:"user"`
}

type RegisterRequest struct {
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name  string `json:"name" binding:"required"`
}

type LoginRequest struct {
	Email string `json:"email" binding:"required,email"`					
	Password string `json:"password" binding:"required,min=8"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type WebAuthResponse struct {
	User User `json:"user"`
}

const (
	AuthProviderLocal = "local"
	AuthProviderGoogle = "google"
)

