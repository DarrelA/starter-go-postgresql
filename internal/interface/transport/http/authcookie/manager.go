// Package authcookie owns the browser-cookie representation of application sessions.
package authcookie

import (
	"time"

	"github.com/DarrelA/starter-go-postgresql/internal/application/dto"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	"github.com/gofiber/fiber/v2"
)

const (
	AccessTokenName  = "access_token"
	RefreshTokenName = "refresh_token"
	CSRFTokenName    = "csrf_token"
	CSRFHeaderName   = "X-CSRF-Token"
	Path             = "/"
	SameSite         = "Strict"
)

type Manager struct{ config *entity.JWTConfig }

func NewManager(config *entity.JWTConfig) *Manager { return &Manager{config: config} }

func (m *Manager) Set(c *fiber.Ctx, session *dto.AuthSession) {
	c.Cookie(m.authCookie(
		AccessTokenName, session.AccessToken, m.config.AccessTokenMaxAge*60,
	))
	c.Cookie(m.authCookie(
		RefreshTokenName, session.RefreshToken, m.config.RefreshTokenMaxAge*60,
	))
	c.Cookie(&fiber.Cookie{
		Name: CSRFTokenName, Value: session.CSRFToken, Path: Path,
		Domain: m.config.Domain, MaxAge: m.config.RefreshTokenMaxAge * 60,
		Secure: m.config.Secure, HTTPOnly: false, SameSite: SameSite,
	})
}

func (m *Manager) Clear(c *fiber.Ctx) {
	for _, name := range []string{AccessTokenName, RefreshTokenName} {
		cookie := m.authCookie(name, "", -1)
		cookie.Expires = time.Unix(0, 0).UTC()
		c.Cookie(cookie)
	}
	c.Cookie(&fiber.Cookie{
		Name: CSRFTokenName, Value: "", Path: Path, Domain: m.config.Domain,
		Expires: time.Unix(0, 0).UTC(), MaxAge: -1, Secure: m.config.Secure,
		HTTPOnly: false, SameSite: SameSite,
	})
}

func (m *Manager) authCookie(name, value string, maxAgeSeconds int) *fiber.Cookie {
	return &fiber.Cookie{
		Name: name, Value: value, Path: Path, Domain: m.config.Domain,
		MaxAge: maxAgeSeconds, Secure: m.config.Secure, HTTPOnly: true, SameSite: SameSite,
	}
}
