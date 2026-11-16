package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/overm-app/api-auth/internal/domain/models"
)

const (
	accessTokenCookie  = "overm_access_token"
	refreshTokenCookie = "overm_refresh_token"
	csrfTokenCookie    = "overm_csrf_token"

	refreshCookiePath = "/auth/v1"
)

type CookieConfig struct {
	Secure     bool
	Domain     string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

type CookieManager struct {
	cfg CookieConfig
}

func NewCookieManager(cfg CookieConfig) *CookieManager {
	return &CookieManager{cfg: cfg}
}

type cookieSpec struct {
	name     string
	path     string
	httpOnly bool
	sameSite http.SameSite
}

var (
	accessSpec  = cookieSpec{accessTokenCookie, "/", true, http.SameSiteLaxMode}
	refreshSpec = cookieSpec{refreshTokenCookie, refreshCookiePath, true, http.SameSiteStrictMode}
	csrfSpec    = cookieSpec{csrfTokenCookie, "/", false, http.SameSiteLaxMode}
)

func (m *CookieManager) write(c *gin.Context, s cookieSpec, value string, maxAge int) {
	c.SetSameSite(s.sameSite)
	c.SetCookie(s.name, value, maxAge, s.path, m.cfg.Domain, m.cfg.Secure, s.httpOnly)
}

func (m *CookieManager) SetAuth(c *gin.Context, resp *models.AuthResponse) {
	m.write(c, accessSpec, resp.AccessToken, seconds(m.cfg.AccessTTL))
	m.write(c, refreshSpec, resp.RefreshToken, seconds(m.cfg.RefreshTTL))
	m.write(c, csrfSpec, resp.CSRFToken, seconds(m.cfg.RefreshTTL))
}

func (m *CookieManager) ClearAuth(c *gin.Context) {
	for _, s := range []cookieSpec{accessSpec, refreshSpec, csrfSpec} {
		m.write(c, s, "", -1)
	}
}

func seconds(d time.Duration) int { return int(d.Seconds()) }
