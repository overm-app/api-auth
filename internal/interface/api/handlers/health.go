package handlers

import (
	"context"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	appErrors "github.com/overm-app/api-auth/internal/domain/errors"
	"github.com/overm-app/api-auth/internal/interface/api/response"
)

type Pinger interface {
	PingContext(ctx context.Context) error
}

type HealthHandler struct {
	db Pinger
	sugar *zap.SugaredLogger
	shuttingDown atomic.Bool
}

func NewHealthHandler(db Pinger, sugar *zap.SugaredLogger) *HealthHandler {
	return &HealthHandler{
		db: db,
		sugar: sugar,
	}
}

func (h *HealthHandler) SetShuttingDown() { h.shuttingDown.Store(true) }

func (h *HealthHandler) Live(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "api-auth",
	})
}

func (h *HealthHandler) Ready(c *gin.Context) {
	if h.shuttingDown.Load() {
		response.HandleError(c, h.sugar, appErrors.ServiceUnavailable("Service is shutting down", nil), c.GetString("request_id"))
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	if err := h.db.PingContext(ctx); err != nil {
		response.HandleError(c, h.sugar, appErrors.ServiceUnavailable("Database connection could not be established", err), c.GetString("request_id"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "api-auth",
	})
}