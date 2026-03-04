package service

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"

	"github.com/overm-app/api-auth/internal/domain/models"  
    appErrors "github.com/overm-app/api-auth/internal/domain/errors"
)

type JWTServiceInter interface {
	GenerateToken(user *models.User) (string, error)
	ValidateToken(tokenString string) (*models.JWTClaims, error)
}

type JWTService struct {
	secret     []byte
	expiration time.Duration
	sugar      *zap.SugaredLogger
}

func NewJWTService(secret []byte, expiration time.Duration, sugar *zap.SugaredLogger) JWTServiceInter {
	return &JWTService{
		secret:     secret,
		expiration: expiration,
		sugar:      sugar,
	}
}

func (j *JWTService) GenerateToken(user *models.User) (string, error) {
	claims := &models.JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.Id,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "overm-api-auth",
		},
		Email: user.Email,
		Name:  user.Name,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(j.secret)
	if err != nil {
		return "", appErrors.Internal("Failed to sign token", err)
	}

	return signed, nil
}

func (j *JWTService) ValidateToken(tokenString string) (*models.JWTClaims, error) {
	parsed, err := jwt.ParseWithClaims(tokenString, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, appErrors.Unauthorized(appErrors.ErrTokenInvalid, "Unexpected signing method")
		}
		return j.secret, nil
	})

	if err != nil {
		return nil, appErrors.Unauthorized(appErrors.ErrTokenInvalid, "Invalid or expired token")
	}

	claims, ok := parsed.Claims.(*models.JWTClaims)
	if !ok || !parsed.Valid {
		return nil, appErrors.Unauthorized(appErrors.ErrTokenInvalid, "Invalid token claims")
	}

	return claims, nil
}
