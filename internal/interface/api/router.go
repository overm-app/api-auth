package api

import (
	"errors"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/overm-app/api-auth/internal/domain/ports"
	"github.com/overm-app/api-auth/internal/interface/api/handlers"
	appErrors "github.com/overm-app/api-auth/internal/domain/errors"
)

type Router struct {
	authHandler   *handlers.AuthHandler
	userHandler   *handlers.UserHandler
	healthHandler *handlers.HealthHandler
	jwtService    ports.JWTService
	sugar         *zap.SugaredLogger
}

func NewRouter(authHandler *handlers.AuthHandler, userHandler *handlers.UserHandler, healthHandler *handlers.HealthHandler, jwtService ports.JWTService, sugar *zap.SugaredLogger) *Router {
	return &Router{
		authHandler:   authHandler,
		userHandler:   userHandler,
		healthHandler: healthHandler,
		jwtService:    jwtService,
		sugar:         sugar,
	}
}

func (r *Router) SetupRouter(logger *zap.SugaredLogger) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	ginEngine := gin.New()
	ginEngine.Use(requestIDMiddleware())
	ginEngine.Use(ginZapLogger(logger))
	ginEngine.Use(gin.Recovery())

	ginEngine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // your frontend dev URL
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
	ginEngine.GET("/health/livez", r.healthHandler.Live)
	ginEngine.GET("/health/readyz", r.healthHandler.Ready)

	auth := ginEngine.Group("/auth/v1")
	{
		auth.POST("/login", r.authHandler.Login)
		auth.POST("/register", r.userHandler.Register)
	}
}

func ginZapLogger(sugar *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		switch c.Request.URL.Path {
		case "/health/livez", "/health/readyz":
			c.Next()
			return
		}
		start := time.Now()
		path := c.Request.URL.Path
		clientType := c.GetHeader("X-Client-Type")

		c.Next()

		status := c.Writer.Status()
		fields := []any{
			"request_id", c.GetString("request_id"),
			"status", status,
			"method", c.Request.Method,
			"client_type", clientType,
			"path", path,
			"latency_ms", time.Since(start).Milliseconds(),
			"ip", c.ClientIP(),
		}

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			if appErr, ok := errors.AsType[*appErrors.AppError](err); ok {
				fields = append(fields, "code", appErr.Code)
				if appErr.Err != nil {
					fields = append(fields, "cause", appErr.Err.Error())
				}
			} else {
				fields = append(fields, "error", err.Error())
			}
		}

		switch {
		case status >= 500:
			sugar.Errorw("HTTP request", fields...)
		case status >= 400:
			sugar.Infow("HTTP request", fields...)
		default:
			sugar.Infow("HTTP request", fields...)
		}
	}
}

func requestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = c.GetHeader("X-Amzn-Trace-Id")
		}
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}
