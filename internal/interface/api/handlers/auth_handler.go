package handlers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	appErrors "github.com/overm-app/api-auth/internal/domain/errors"
	"github.com/overm-app/api-auth/internal/domain/models"
	"github.com/overm-app/api-auth/internal/usecase"
	"github.com/overm-app/api-auth/internal/interface/api/response"
)

type AuthHandler struct {
	loginUseCase *usecase.LoginUseCase
	cookieConfig CookieConfig
	sugar        *zap.SugaredLogger
}

func NewAuthHandler(loginUseCase *usecase.LoginUseCase, cookieConfig CookieConfig, sugar *zap.SugaredLogger) *AuthHandler {
	return &AuthHandler{
		loginUseCase: loginUseCase,
		cookieConfig: cookieConfig,
		sugar:        sugar,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, appErrors.Validation(appErrors.ErrValidation, err.Error()))
		return
	}

	resp, err := (*h.loginUseCase).Execute(c.Request.Context(), &req)
	if err != nil {
		response.HandleError(c, err)
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
