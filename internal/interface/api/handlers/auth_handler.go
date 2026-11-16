package handlers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/overm-app/api-auth/internal/domain/models"
	"github.com/overm-app/api-auth/internal/usecase"
	"github.com/overm-app/api-auth/internal/interface/api/response"
)

type AuthHandler struct {
	loginUseCase *usecase.LoginUseCase
	cookieManager *CookieManager
	sugar        *zap.SugaredLogger
}

func NewAuthHandler(loginUseCase *usecase.LoginUseCase, cookieManager *CookieManager, sugar *zap.SugaredLogger) *AuthHandler {
	return &AuthHandler{
		loginUseCase: loginUseCase,
		cookieManager: cookieManager,
		sugar:        sugar,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if !response.BindJSON(c, &req) {
		return
	}

	resp, err := (*h.loginUseCase).Execute(c.Request.Context(), &req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	clientType := c.GetHeader("X-Client-Type")
	if clientType == "web" {
		h.cookieManager.SetAuth(c, resp)
		c.JSON(200, models.WebAuthResponse{User: resp.User})
		return
	}

	c.JSON(200, resp)
}
