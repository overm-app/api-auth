package api

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SetupRouter(logger *zap.SugaredLogger) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
    r := gin.New()
    r.Use(ginZapLogger(logger))
    r.Use(gin.Recovery())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	SetupRoutes(r)

	return r
}

func SetupRoutes(r *gin.Engine) {
	r.GET("/healthz", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "status":  "ok",
            "service": "api-auth",
        })
    })

	auth := r.Group("/auth/v1/")
	{
		auth.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
	}
}

func ginZapLogger(sugar *zap.SugaredLogger) gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path

        c.Next()

		fields := []any{
			"status",     c.Writer.Status(),
			"method",     c.Request.Method,
			"path",       path,
			"latency_ms", time.Since(start).Milliseconds(),
			"ip",         c.ClientIP(),
		}

		if errs := c.Errors.ByType(gin.ErrorTypePrivate).String(); errs != "" {
			fields = append(fields, "error", errs)
			sugar.Errorw("HTTP request with error", fields...)
			return
		}

		sugar.Infow("HTTP request", fields...)
    }
}