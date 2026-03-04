package handlers

import (
	"os"

	"github.com/gin-gonic/gin"

	"github.com/overm-app/api-auth/internal/domain/models"
)

const (
    accessTokenCookie  = "overm_access_token"
    refreshTokenCookie = "overm_refresh_token"
    csrfTokenCookie    = "overm_csrf_token"

    accessTokenMaxAge  = 60 * 60 * 24        // 24h
    refreshTokenMaxAge = 60 * 60 * 24 * 30   // 30d
    csrfTokenMaxAge    = 60 * 60 * 24        // 24h
)

type CookieConfig struct {
	Secure bool
	Domain string
}

func NewCookieConfig() CookieConfig {
	return CookieConfig{
		Secure: os.Getenv("ENV") == "production",
		Domain: os.Getenv("COOKIE_DOMAIN"),
	}
}

func (cfg CookieConfig) setAuthCookies(c *gin.Context, resp *models.AuthResponse) {
    c.SetCookie(accessTokenCookie, resp.AccessToken, accessTokenMaxAge, "/", cfg.Domain, cfg.Secure, true)
    c.SetCookie(refreshTokenCookie, resp.RefreshToken, refreshTokenMaxAge, "/auth/v1/refresh", cfg.Domain, cfg.Secure, true)
    c.SetCookie(csrfTokenCookie, resp.CSRFToken, csrfTokenMaxAge, "/", cfg.Domain, cfg.Secure, false)
}

// func (cfg CookieConfig) clearAuthCookies(c *gin.Context) {
//     c.SetCookie(accessTokenCookie, "", -1, "/", "", cfg.Secure, true)
//     c.SetCookie(refreshTokenCookie, "", -1, "/auth/v1/refresh", "", cfg.Secure, true)
//     c.SetCookie(csrfTokenCookie, "", -1, "/", "", cfg.Secure, false)
// }