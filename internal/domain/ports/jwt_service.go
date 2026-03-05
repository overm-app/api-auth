package ports

import "github.com/overm-app/api-auth/internal/domain/models"

type JWTService interface {
    GenerateToken(user *models.User) (string, error)
    ValidateToken(tokenString string) (*models.JWTClaims, error)
}