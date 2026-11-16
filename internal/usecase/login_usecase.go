package usecase

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"

	appErrors "github.com/overm-app/api-auth/internal/domain/errors"
	"github.com/overm-app/api-auth/internal/domain/models"
	"github.com/overm-app/api-auth/internal/domain/ports"
	"github.com/overm-app/api-auth/internal/infrastructure/service"
)

type LoginUseCase struct {
	userRepo   ports.UserRepository
	tokenRepo  ports.TokenRepository
	jwtService ports.JWTService
	refreshTTL time.Duration
}

func NewLoginUseCase(userRepo ports.UserRepository, tokenRepo ports.TokenRepository, jwtService ports.JWTService, refreshTTL time.Duration) *LoginUseCase {
	return &LoginUseCase{
		userRepo:   userRepo,
		tokenRepo:  tokenRepo,
		jwtService: jwtService,
		refreshTTL: refreshTTL,
	}
}

func (uc *LoginUseCase) Execute(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error) {
	user, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, appErrors.Unauthorized(appErrors.ErrInvalidCredentials, "Email or password is incorrect")
	}

	if err := verifyPassword(req.Password, user.PasswordHash); err != nil {
		return nil, err
	}

	accesstoken, err := uc.jwtService.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	rawToken, tokenID, err := service.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	if err := uc.tokenRepo.DeleteByUserID(ctx, user.ID); err != nil {
		return nil, err
	}

	refreshToken := &models.RefreshToken{
		ID:        tokenID,
		UserID:    user.ID,
		Token:     rawToken,
		ExpiresAt: time.Now().Add(uc.refreshTTL),
	}

	if err := uc.tokenRepo.Save(ctx, refreshToken); err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		AccessToken:  accesstoken,
		RefreshToken: rawToken,
		User:         *user,
	}, nil
}

func verifyPassword(providedPassword, storedHash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(providedPassword))
	if err != nil {
		return appErrors.Unauthorized(appErrors.ErrInvalidCredentials, "Email or password is incorrect")
	}
	return nil
}
