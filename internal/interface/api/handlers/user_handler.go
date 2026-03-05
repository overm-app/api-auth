package handlers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	appErrors "github.com/overm-app/api-auth/internal/domain/errors"
	"github.com/overm-app/api-auth/internal/domain/models"
	"github.com/overm-app/api-auth/internal/usecase"
	"github.com/overm-app/api-auth/internal/interface/api/response"
)

type UserHandler struct {
	registerUseCase *usecase.RegisterUseCase
	cookieConfig   CookieConfig
	sugar           *zap.SugaredLogger
}

func NewUserHandler(registerUseCase *usecase.RegisterUseCase, cookieConfig CookieConfig, sugar *zap.SugaredLogger) *UserHandler {
	return &UserHandler{
		registerUseCase: registerUseCase,
		cookieConfig:   cookieConfig,
		sugar:           sugar,
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	requestID := c.GetString("request_id")

	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, h.sugar, appErrors.Validation(appErrors.ErrValidation, err.Error()), requestID)
		return
	}

	resp, err := (*h.registerUseCase).Execute(c.Request.Context(), &req)
	if err != nil {
		response.HandleError(c, h.sugar, err, requestID)
		return
	}

	clientType := c.GetHeader("X-Client-Type")
	if clientType == "web" {
		h.cookieConfig.setAuthCookies(c, resp)
		c.JSON(200, models.WebAuthResponse{User: resp.User})
		return
	}

	c.JSON(200, resp)
}