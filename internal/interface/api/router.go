package api

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/overm-app/api-auth/internal/domain/ports"
	"github.com/overm-app/api-auth/internal/interface/api/handlers"
)

type Router struct {
	authHandler *handlers.AuthHandler
	userHandler *handlers.UserHandler
	jwtService  ports.JWTService
	sugar       *zap.SugaredLogger
}

func NewRouter(authHandler *handlers.AuthHandler, userHandler *handlers.UserHandler, jwtService ports.JWTService, sugar *zap.SugaredLogger) *Router {
	return &Router{
		authHandler: authHandler,
		userHandler: userHandler,
		jwtService:  jwtService,
		sugar:       sugar,
	}
}

func (r *Router) SetupRouter(logger *zap.SugaredLogger) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	ginEngine := gin.New()
	ginEngine.Use(ginZapLogger(logger))
	ginEngine.Use(gin.Recovery())

	ginEngine.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"http://localhost:3000"}, // your frontend dev URL
    AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Client-Type", "X-CSRF-Token"},
    ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
    AllowCredentials: true,
    MaxAge:           12 * time.Hour,
	}))

	r.setupRoutes(ginEngine)

	return ginEngine
}

func (r *Router) setupRoutes(ginEngine *gin.Engine) {
	ginEngine.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "api-auth",
		})
	})

	auth := ginEngine.Group("/auth/v1")
	{
		auth.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
		auth.POST("/login", r.authHandler.Login)
		auth.POST("/register", r.userHandler.Register)
	}
}

func ginZapLogger(sugar *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		fields := []any{
			"status", c.Writer.Status(),
			"method", c.Request.Method,
			"path", path,
			"latency_ms", time.Since(start).Milliseconds(),
			"ip", c.ClientIP(),
		}

		if errs := c.Errors.ByType(gin.ErrorTypePrivate).String(); errs != "" {
			fields = append(fields, "error", errs)
			sugar.Errorw("HTTP request with error", fields...)
			return
		}

		sugar.Infow("HTTP request", fields...)
	}
}
