package http

import (
	"context"
	stdhttp "net/http"
	"net/http/httptest"
	"testing"

	"github.com/DarrelA/starter-go-postgresql/internal/domain/apperror"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	"github.com/DarrelA/starter-go-postgresql/internal/infrastructure/config"
	"github.com/gofiber/fiber/v2"
)

func TestRouterRegistersLogoutAsPost(t *testing.T) {
	handler, repository := newAuthTestHandler()
	envConfig := &config.EnvConfig{EnvConfig: entity.EnvConfig{
		Env:            "test",
		BaseURLsConfig: &entity.BaseURLsConfig{},
		CORSConfig:     &entity.CORSConfig{AllowedOrigins: "http://localhost"},
	}}
	app := NewRouter(envConfig, repository, handler.ts, noOpUserRoutes{}, handler, noOpOAuthRoutes{})

	getReq := httptest.NewRequest(stdhttp.MethodGet, "/auth/api/v1/users/logout", nil)
	getReq.AddCookie(&stdhttp.Cookie{Name: accessTokenCookieName, Value: testAccessToken})
	getReq.AddCookie(&stdhttp.Cookie{Name: refreshTokenCookieName, Value: testRefreshToken})
	getResp, err := app.Test(getReq)
	if err != nil {
		t.Fatalf("execute GET logout: %v", err)
	}
	defer getResp.Body.Close()
	if getResp.StatusCode != stdhttp.StatusNotFound {
		t.Fatalf("expected GET logout status %d, got %d", stdhttp.StatusNotFound, getResp.StatusCode)
	}

	postReq := httptest.NewRequest(stdhttp.MethodPost, "/auth/api/v1/users/logout", nil)
	postReq.AddCookie(&stdhttp.Cookie{Name: accessTokenCookieName, Value: testAccessToken})
	postReq.AddCookie(&stdhttp.Cookie{Name: refreshTokenCookieName, Value: testRefreshToken})
	postReq.AddCookie(&stdhttp.Cookie{Name: csrfTokenCookieName, Value: "csrf-refresh-token-id"})
	postReq.Header.Set(csrfTokenHeaderName, "csrf-refresh-token-id")
	postResp, err := app.Test(postReq)
	if err != nil {
		t.Fatalf("execute POST logout: %v", err)
	}
	defer postResp.Body.Close()
	if postResp.StatusCode != stdhttp.StatusOK {
		t.Fatalf("expected POST logout status %d, got %d", stdhttp.StatusOK, postResp.StatusCode)
	}
}

func TestRouterRegistersRefreshAsPost(t *testing.T) {
	handler, repository := newAuthTestHandler()
	envConfig := &config.EnvConfig{EnvConfig: entity.EnvConfig{
		Env:            "test",
		BaseURLsConfig: &entity.BaseURLsConfig{},
		CORSConfig:     &entity.CORSConfig{AllowedOrigins: "http://localhost"},
	}}
	app := NewRouter(envConfig, repository, handler.ts, noOpUserRoutes{}, handler, noOpOAuthRoutes{})

	response, err := app.Test(httptest.NewRequest(stdhttp.MethodGet, "/auth/api/v1/users/refresh", nil))
	if err != nil {
		t.Fatalf("execute GET refresh: %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode != stdhttp.StatusNotFound {
		t.Fatalf("expected GET refresh status %d, got %d", stdhttp.StatusNotFound, response.StatusCode)
	}

	request := httptest.NewRequest(stdhttp.MethodPost, "/auth/api/v1/users/refresh", nil)
	request.AddCookie(&stdhttp.Cookie{Name: refreshTokenCookieName, Value: testRefreshToken})
	request.AddCookie(&stdhttp.Cookie{Name: csrfTokenCookieName, Value: "csrf-refresh-token-id"})
	request.Header.Set(csrfTokenHeaderName, "csrf-refresh-token-id")
	response, err = app.Test(request)
	if err != nil {
		t.Fatalf("execute POST refresh: %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode != stdhttp.StatusOK {
		t.Fatalf("expected POST refresh status %d without an access token, got %d", stdhttp.StatusOK, response.StatusCode)
	}
}

func TestProtectedRouteRejectsInvalidSessionWithoutLoadingUser(t *testing.T) {
	handler, repository := newAuthTestHandler()
	repository.getUserErr = apperror.ErrSessionNotFound
	users := &countingUserReader{}
	envConfig := &config.EnvConfig{EnvConfig: entity.EnvConfig{
		Env:            "test",
		BaseURLsConfig: &entity.BaseURLsConfig{},
		CORSConfig:     &entity.CORSConfig{AllowedOrigins: "http://localhost"},
	}}
	app := NewRouter(
		envConfig, repository, handler.ts, NewUserHandler(users), handler, noOpOAuthRoutes{},
	)
	request := httptest.NewRequest(stdhttp.MethodGet, "/auth/api/v1/users/me", nil)
	request.AddCookie(&stdhttp.Cookie{Name: accessTokenCookieName, Value: testAccessToken})
	response, err := app.Test(request)
	if err != nil {
		t.Fatalf("execute protected request: %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode != stdhttp.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", stdhttp.StatusUnauthorized, response.StatusCode)
	}
	if users.calls != 0 {
		t.Fatalf("invalid session caused %d PostgreSQL user lookup(s)", users.calls)
	}
}

type countingUserReader struct{ calls int }

func (r *countingUserReader) GetUserByUUID(context.Context, string) (*entity.User, error) {
	r.calls++
	return nil, apperror.ErrUserNotFound
}

type noOpUserRoutes struct{}

func (noOpUserRoutes) GetUserRecord(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) }

type noOpOAuthRoutes struct{}

func (noOpOAuthRoutes) Login(c *fiber.Ctx) error    { return c.SendStatus(fiber.StatusOK) }
func (noOpOAuthRoutes) Callback(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) }
