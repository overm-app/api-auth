package models

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTClaims struct {
	jwt.RegisteredClaims
	Email string `json:"email"`
	Name  string `json:"name"`
}

type RefreshToken struct {
	ID        string
	UserID    int64
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
}
