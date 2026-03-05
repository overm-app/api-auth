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

const cost = bcrypt.DefaultCost

type RegisterUseCase struct {
	userRepo   ports.UserRepository
	tokenRepo  ports.TokenRepository
	jwtService ports.JWTService
	

}

func NewRegisterUseCase(userRepo ports.UserRepository, tokenRepo ports.TokenRepository, jwtService ports.JWTService) *RegisterUseCase {
	return &RegisterUseCase{
		userRepo:   userRepo,
		tokenRepo:  tokenRepo,
		jwtService: jwtService,
	}
}

func (uc *RegisterUseCase) Execute(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error) {
	existing, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, appErrors.Internal("Failed to check existing user", err)
	}

	if existing != nil {
		return nil, appErrors.Conflict(appErrors.ErrEmailAlreadyExists, "Email is already registered")
	}

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return nil, appErrors.Internal("Failed to hash password", err)
	}

	user := &models.User{
		Email:        req.Email,
		Name:         req.Name,
		PasswordHash: hashedPassword,
		AuthProvider: models.AuthProviderLocal,
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, appErrors.Internal("failed to create user", err)
	}

	accesstoken, err := uc.jwtService.GenerateToken(user)
	if err != nil {
		return nil, appErrors.Internal("failed to generate access token", err)
	}

	rawToken, tokenID, err := service.GenerateRefreshToken()
	if err != nil {
		return nil, appErrors.Internal("failed to generate refresh token", err)
	}

	refreshToken := &models.RefreshToken{
		ID:        tokenID,
		UserID:    user.ID,
		Token:     rawToken,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}

	if err := uc.tokenRepo.Save(ctx, refreshToken); err != nil {
		return nil, appErrors.Internal("failed to save refresh token", err)
	}

	return &models.AuthResponse{
		AccessToken:  accesstoken,
		RefreshToken: rawToken,
		User:         *user,
	}, nil
}

func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}
