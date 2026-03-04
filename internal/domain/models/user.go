package models

import "time"

type User struct {
    ID           int64     `json:"-"`          
    PublicID     string    `json:"id"`
    Email        string    `json:"email"`
    Name         string    `json:"name"`
    PasswordHash string    `json:"-"`
    AuthProvider string    `json:"auth_provider"`
    AvatarURL    *string   `json:"avatar_url"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}