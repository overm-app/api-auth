package usecase

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	appErrors "github.com/overm-app/api-auth/internal/domain/errors"
	"github.com/overm-app/api-auth/internal/domain/models"
	"github.com/overm-app/api-auth/internal/domain/ports"
)

type LoginUseCase interface {
	Execute(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error)
}

type loginUseCase struct {
	userRepo ports.UserRepository
}

func NewLoginUseCase(userRepo ports.UserRepository) LoginUseCase {
	return &loginUseCase{
		userRepo: userRepo,
	}
}

func (uc *loginUseCase) Execute(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error) {
	user, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, appErrors.Internal("Failed to find user", err)
	}

	if user == nil {
		return nil, appErrors.Unauthorized(appErrors.ErrInvalidCredentials, "Email or password is incorrect")
	}

	if err := verifyPassword(req.Password, user.PasswordHash); err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		AccessToken: "fake-access-token",
	}, nil
}

func verifyPassword(providedPassword, storedHash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(providedPassword))
	if err != nil {
		return appErrors.Unauthorized(appErrors.ErrInvalidCredentials, "Email or password is incorrect")
	}
	return nil
}
