package ports

import (
	"context"

	"github.com/overm-app/api-auth/internal/domain/models"
)

type TokenRepository interface {
	Save(ctx context.Context, token *models.RefreshToken) error
	GetByUserID(ctx context.Context, userID string) (*models.RefreshToken, error)
	GetByID(ctx context.Context, id string) (*models.RefreshToken, error)
	DeleteByUserID(ctx context.Context, userID string) error
	RotateToken(ctx context.Context, oldID string, newToken *models.RefreshToken) error
}