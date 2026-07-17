package http

import (
	"context"
	"encoding/json"
	stdhttp "net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DarrelA/starter-go-postgresql/internal/application/dto"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/apperror"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	testAccessToken  = "signed-access-token"
	testRefreshToken = "signed-refresh-token"
)

func TestAuthCookiesAreNotReturnedInResponseBodies(t *testing.T) {
	handler, repository := newAuthTestHandler()

	t.Run("login", func(t *testing.T) {
		app := fiber.New()
		app.Use(func(c *fiber.Ctx) error {
			c.Locals("login_payload", dto.LoginInput{Email: "user@example.com", Password: "Password1!"})
			return c.Next()
		})
		app.Post("/login", handler.Login)

		resp, err := app.Test(httptest.NewRequest(stdhttp.MethodPost, "/login", nil))
		if err != nil {
			t.Fatalf("execute login: %v", err)
		}
		defer resp.Body.Close()

		assertSuccessBodyWithoutTokens(t, resp)
		cookies := cookiesByName(resp.Cookies())
		assertIssuedCookie(t, cookies[accessTokenCookieName], testAccessToken, handler.jwt.AccessTokenMaxAge*60)
		assertIssuedCookie(t, cookies[refreshTokenCookieName], testRefreshToken, handler.jwt.RefreshTokenMaxAge*60)
		assertCSRFCookie(t, cookies[csrfTokenCookieName], "csrf-refresh-token-id")
		if len(repository.created) != 1 {
			t.Fatalf("expected one token family, got %d", len(repository.created))
		}
	})

	t.Run("refresh", func(t *testing.T) {
		app := fiber.New()
		app.Post("/refresh", handler.RefreshAccessToken)
		req := httptest.NewRequest(stdhttp.MethodPost, "/refresh", nil)
		req.AddCookie(&stdhttp.Cookie{Name: refreshTokenCookieName, Value: testRefreshToken})
		req.AddCookie(&stdhttp.Cookie{Name: csrfTokenCookieName, Value: "csrf-refresh-token-id"})
		req.Header.Set(csrfTokenHeaderName, "csrf-refresh-token-id")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("execute refresh: %v", err)
		}
		defer resp.Body.Close()

		assertSuccessBodyWithoutTokens(t, resp)
		cookies := cookiesByName(resp.Cookies())
		assertIssuedCookie(t, cookies[accessTokenCookieName], testAccessToken, handler.jwt.AccessTokenMaxAge*60)
		assertIssuedCookie(t, cookies[refreshTokenCookieName], testRefreshToken, handler.jwt.RefreshTokenMaxAge*60)
		assertCSRFCookie(t, cookies[csrfTokenCookieName], "csrf-refresh-token-id")
		if len(repository.rotated) != 1 {
			t.Fatalf("expected one token rotation, got %d", len(repository.rotated))
		}
	})
}

func TestRefreshRejectsMissingOrMismatchedCSRFToken(t *testing.T) {
	handler, _ := newAuthTestHandler()
	for _, test := range []struct {
		name, cookie, header string
	}{
		{name: "missing", cookie: "", header: ""},
		{name: "mismatch", cookie: "csrf-refresh-token-id", header: "different"},
		{name: "invalid signature", cookie: "invalid", header: "invalid"},
	} {
		t.Run(test.name, func(t *testing.T) {
			app := fiber.New()
			app.Post("/refresh", handler.RefreshAccessToken)
			req := httptest.NewRequest(stdhttp.MethodPost, "/refresh", nil)
			req.AddCookie(&stdhttp.Cookie{Name: refreshTokenCookieName, Value: testRefreshToken})
			if test.cookie != "" {
				req.AddCookie(&stdhttp.Cookie{Name: csrfTokenCookieName, Value: test.cookie})
			}
			req.Header.Set(csrfTokenHeaderName, test.header)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("execute refresh: %v", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != stdhttp.StatusForbidden {
				t.Fatalf("expected status %d, got %d", stdhttp.StatusForbidden, resp.StatusCode)
			}
		})
	}
}

func TestLogoutClearsCookiesUsingIssuanceAttributes(t *testing.T) {
	handler, repository := newAuthTestHandler()
	app := fiber.New()
	app.Post("/logout", handler.Logout)
	req := httptest.NewRequest(stdhttp.MethodPost, "/logout", nil)
	req.AddCookie(&stdhttp.Cookie{Name: refreshTokenCookieName, Value: testRefreshToken})
	req.AddCookie(&stdhttp.Cookie{Name: csrfTokenCookieName, Value: "csrf-refresh-token-id"})
	req.Header.Set(csrfTokenHeaderName, "csrf-refresh-token-id")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("execute logout: %v", err)
	}
	defer resp.Body.Close()

	assertSuccessBodyWithoutTokens(t, resp)
	if len(repository.revoked) != 1 || repository.revoked[0].refreshTokenUUID != "refresh-token-id" {
		t.Fatalf("unexpected revoked sessions: %#v", repository.revoked)
	}

	cookies := cookiesByName(resp.Cookies())
	for _, name := range []string{accessTokenCookieName, refreshTokenCookieName, csrfTokenCookieName} {
		cookie := cookies[name]
		if cookie == nil {
			t.Fatalf("missing expired %s cookie", name)
		}
		if cookie.Value != "" || cookie.MaxAge >= 0 || !cookie.Expires.Before(time.Now()) {
			t.Fatalf("cookie %s was not expired: %#v", name, cookie)
		}
		if name == csrfTokenCookieName {
			assertCSRFCookieAttributes(t, cookie)
		} else {
			assertCookieAttributes(t, cookie)
		}
	}
}

type tokenRevocation struct {
	familyUUID, refreshTokenUUID string
}

type authTokenRepositoryStub struct {
	userUUID     string
	created      []entity.TokenSession
	rotated      []entity.TokenSession
	revoked      []tokenRevocation
	rotateErr    error
	getUserErr   error
	getUserCalls int
}

func (r *authTokenRepositoryStub) Create(_ context.Context, session entity.TokenSession) error {
	r.created = append(r.created, session)
	return nil
}

func (r *authTokenRepositoryStub) Rotate(_ context.Context, _ string, session entity.TokenSession) error {
	if r.rotateErr != nil {
		return r.rotateErr
	}
	r.rotated = append(r.rotated, session)
	return nil
}

func (r *authTokenRepositoryStub) GetUserUUID(context.Context, string, string) (string, error) {
	r.getUserCalls++
	return r.userUUID, r.getUserErr
}

func (r *authTokenRepositoryStub) Revoke(_ context.Context, familyUUID, refreshTokenUUID string) error {
	r.revoked = append(r.revoked, tokenRevocation{familyUUID: familyUUID, refreshTokenUUID: refreshTokenUUID})
	return nil
}

type authUserServiceStub struct{ userUUID uuid.UUID }

func (s authUserServiceStub) CreateUser(context.Context, dto.RegisterInput) (*dto.UserResponse, error) {
	return &dto.UserResponse{UUID: &s.userUUID}, nil
}

func (s authUserServiceStub) Authenticate(context.Context, dto.LoginInput) (*dto.UserResponse, error) {
	return &dto.UserResponse{UUID: &s.userUUID}, nil
}

func (s authUserServiceStub) GetUserByUUID(context.Context, string) (*entity.User, error) {
	return &entity.User{UUID: &s.userUUID}, nil
}

type authSessionServiceStub struct{ repository *authTokenRepositoryStub }

func (s authSessionServiceStub) Issue(_ context.Context, userUUID string) (*dto.AuthSession, error) {
	s.repository.created = append(s.repository.created, entity.TokenSession{
		UserUUID: userUUID, FamilyUUID: "family-id",
		AccessTokenUUID: "access-token-id", RefreshTokenUUID: "refresh-token-id",
	})
	return &dto.AuthSession{
		AccessToken: testAccessToken, RefreshToken: testRefreshToken, CSRFToken: "csrf-refresh-token-id",
	}, nil
}

type authTokenServiceStub struct {
	access  entity.Token
	refresh entity.Token
}

func (s *authTokenServiceStub) CreateAccessToken(userUUID, familyUUID string, _ time.Duration) (*entity.Token, error) {
	token := s.access
	token.UserUUID, token.FamilyUUID = userUUID, familyUUID
	return &token, nil
}

func (s *authTokenServiceStub) CreateRefreshToken(userUUID, familyUUID string, _ time.Duration) (*entity.Token, error) {
	token := s.refresh
	token.UserUUID, token.FamilyUUID = userUUID, familyUUID
	return &token, nil
}

func (s *authTokenServiceStub) ValidateAccessToken(string) (*entity.Token, error) {
	token := s.access
	return &token, nil
}

type authCSRFServiceStub struct{}

func (authCSRFServiceStub) Create(sessionID string) (string, error) { return "csrf-" + sessionID, nil }
func (authCSRFServiceStub) Validate(sessionID, token string) error {
	if token != "csrf-"+sessionID {
		return apperror.ErrInvalidCSRFToken
	}
	return nil
}

func (s *authTokenServiceStub) ValidateRefreshToken(string) (*entity.Token, error) {
	token := s.refresh
	return &token, nil
}

func newAuthTestHandler() (*AuthHandler, *authTokenRepositoryStub) {
	userUUID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	accessExpires := time.Now().Add(time.Minute).Unix()
	refreshExpires := time.Now().Add(3 * time.Minute).Unix()
	accessToken := testAccessToken
	refreshToken := testRefreshToken
	repository := &authTokenRepositoryStub{userUUID: userUUID.String()}
	tokenService := &authTokenServiceStub{
		access: entity.Token{
			Token: &accessToken, TokenUUID: "access-token-id", UserUUID: userUUID.String(),
			FamilyUUID: "family-id", ExpiresIn: &accessExpires,
		},
		refresh: entity.Token{Token: &refreshToken, TokenUUID: "refresh-token-id", ExpiresIn: &refreshExpires},
	}
	handler := NewAuthHandler(
		repository,
		authUserServiceStub{userUUID: userUUID},
		authSessionServiceStub{repository: repository},
		tokenService,
		authCSRFServiceStub{},
		&entity.JWTConfig{
			Domain:                "example.com",
			Secure:                true,
			HttpOnly:              true,
			AccessTokenMaxAge:     1,
			RefreshTokenMaxAge:    3,
			AccessTokenExpiredIn:  time.Minute,
			RefreshTokenExpiredIn: 3 * time.Minute,
		},
	)
	return handler, repository
}

func assertSuccessBodyWithoutTokens(t *testing.T, resp *stdhttp.Response) {
	t.Helper()
	if resp.StatusCode != stdhttp.StatusOK {
		t.Fatalf("expected status %d, got %d", stdhttp.StatusOK, resp.StatusCode)
	}
	var body map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body["status"] != "success" {
		t.Fatalf("unexpected response body: %#v", body)
	}
	for _, field := range []string{"access_token", "refresh_token"} {
		if _, ok := body[field]; ok {
			t.Fatalf("response exposes %s: %#v", field, body)
		}
	}
}

func cookiesByName(cookies []*stdhttp.Cookie) map[string]*stdhttp.Cookie {
	result := make(map[string]*stdhttp.Cookie, len(cookies))
	for _, cookie := range cookies {
		result[cookie.Name] = cookie
	}
	return result
}

func assertIssuedCookie(t *testing.T, cookie *stdhttp.Cookie, value string, maxAge int) {
	t.Helper()
	if cookie == nil {
		t.Fatal("missing authentication cookie")
	}
	if cookie.Value != value || cookie.MaxAge != maxAge {
		t.Fatalf("unexpected issued cookie: %#v", cookie)
	}
	assertCookieAttributes(t, cookie)
}

func assertCSRFCookie(t *testing.T, cookie *stdhttp.Cookie, value string) {
	t.Helper()
	if cookie == nil || cookie.Value != value {
		t.Fatalf("unexpected CSRF cookie: %#v", cookie)
	}
	assertCSRFCookieAttributes(t, cookie)
}

func assertCSRFCookieAttributes(t *testing.T, cookie *stdhttp.Cookie) {
	t.Helper()
	if cookie.Path != authCookiePath || cookie.Domain != "example.com" || !cookie.Secure || cookie.HttpOnly || cookie.SameSite != stdhttp.SameSiteStrictMode {
		t.Fatalf("unexpected CSRF cookie attributes: %#v", cookie)
	}
}

func assertCookieAttributes(t *testing.T, cookie *stdhttp.Cookie) {
	t.Helper()
	if cookie.Path != authCookiePath || cookie.Domain != "example.com" || !cookie.Secure || !cookie.HttpOnly || cookie.SameSite != stdhttp.SameSiteStrictMode {
		t.Fatalf("unexpected authentication cookie attributes: %#v", cookie)
	}
}
