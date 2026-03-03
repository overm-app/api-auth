package ports

import (
	"context"

	"github.com/overm-app/api-auth/internal/domain/models"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
}
