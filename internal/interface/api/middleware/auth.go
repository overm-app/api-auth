package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	appErrors "github.com/overm-app/api-auth/internal/domain/errors"
	"github.com/overm-app/api-auth/internal/domain/ports"
	"github.com/overm-app/api-auth/internal/interface/api/response"
)

type authMiddleware struct {
}

func NewAuthMiddleware() *authMiddleware {
	return &authMiddleware{}
}

func (m *authMiddleware) JWTAuthMiddleware(jwtService ports.JWTService, sugar *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		clientType := c.GetHeader("X-Client-Type")

		if clientType == "web" {
			cookie, err := c.Cookie("overm_access_token")
			if err != nil {
				response.HandleError(c, appErrors.Unauthorized(appErrors.ErrUnauthorized, "Missing session cookie"))
				return
			}
			tokenString = cookie
		} else {
			bearer := c.GetHeader("Authorization")
			if len(bearer) < 8 || bearer[:7] != "Bearer " {
				response.HandleError(c, appErrors.Unauthorized(appErrors.ErrUnauthorized, "Missing or malformed Authorization header"))
				return
			}
			tokenString = bearer[7:]
		}

		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			response.HandleError(c, appErrors.Unauthorized(appErrors.ErrUnauthorized, "Invalid token"))
			return
		}

		c.Set("user_id", claims.Subject)
		c.Set("user_email", claims.Email)
		c.Set("user_name", claims.Name)
		c.Next()
	}
}
