package middleware

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	appErrors "github.com/overm-app/api-auth/internal/domain/errors"
	"github.com/overm-app/api-auth/internal/interface/api/handlers"
	"github.com/overm-app/api-auth/internal/infrastructure/service"
)

type AuthMiddlewareInterface interface {
	JWTAuthMiddleware(jwtService *service.JWTService, sugar *zap.SugaredLogger) gin.HandlerFunc
	GenerateRefreshToken() (string, string, error)
}

type authMiddleware struct{
}

func NewAuthMiddleware() *authMiddleware {
	return &authMiddleware{}
}

func (m *authMiddleware) JWTAuthMiddleware(jwtService *service.JWTService, sugar *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetString("request_id")
		var tokenString string

		clientType := c.GetHeader("X-Client-Type")

		if clientType == "web" {
			cookie, err := c.Cookie("overm_access_token")
			if err != nil {
				handlers.HandleError(c, sugar, appErrors.Unauthorized(appErrors.ErrUnauthorized, "Missing session cookie"), requestID)
				return
			}
			tokenString = cookie
		} else {
			bearer := c.GetHeader("Authorization")
			if len(bearer) < 8 || bearer[:7] != "Bearer " {
				handlers.HandleError(c, sugar, appErrors.Unauthorized(appErrors.ErrUnauthorized, "Missing or malformed Authorization header"), requestID)
				return
			}
			tokenString = bearer[7:]
		}

		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			handlers.HandleError(c, sugar, err, requestID)
			return
		}

		c.Set("user_id", claims.Subject)
		c.Set("user_email", claims.Email)
		c.Set("user_name", claims.Name)
		c.Next()
	}
}

func (m *authMiddleware) GenerateRefreshToken() (string, string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", "", err
	}
	token := base64.StdEncoding.EncodeToString(b)
	hash := sha256.Sum256([]byte(token))
	id := hex.EncodeToString(hash[:])
	return token, id, nil
}
