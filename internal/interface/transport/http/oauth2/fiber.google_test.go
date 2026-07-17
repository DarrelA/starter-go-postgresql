package http

import (
	"context"
	"encoding/json"
	"io"
	stdhttp "net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/DarrelA/starter-go-postgresql/internal/application/dto"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/apperror"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

func TestLoginCreatesBoundOAuthState(t *testing.T) {
	app := fiber.New()
	oauth := NewGoogleOAuth2(&entity.OAuth2Config{
		GoogleRedirectURL: "http://localhost/auth/google_callback",
		GoogleClientID:    "test-client",
		Scopes:            []string{"openid"},
	}, &entity.JWTConfig{}, nil, nil)
	app.Get("/login", oauth.Login)

	firstState := requestOAuthState(t, app)
	secondState := requestOAuthState(t, app)
	if firstState == secondState {
		t.Fatal("expected a unique OAuth state for each login")
	}
}

func TestCallbackRejectsMismatchedOAuthState(t *testing.T) {
	app := fiber.New()
	oauth := NewGoogleOAuth2(&entity.OAuth2Config{}, &entity.JWTConfig{}, nil, nil)
	app.Get("/callback", oauth.Callback)

	req := httptest.NewRequest(stdhttp.MethodGet, "/callback?state=unexpected&code=unused", nil)
	req.AddCookie(&stdhttp.Cookie{Name: oauthStateCookie, Value: "expected"})
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("execute callback request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != stdhttp.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", stdhttp.StatusBadRequest, resp.StatusCode)
	}
}

func TestCallbackEstablishesLocalApplicationSession(t *testing.T) {
	providerClient := &stdhttp.Client{Transport: roundTripFunc(func(r *stdhttp.Request) (*stdhttp.Response, error) {
		switch r.URL.Path {
		case "/token":
			return providerResponse(`{"access_token":"google-token","token_type":"Bearer"}`), nil
		case "/userinfo":
			if r.Header.Get("Authorization") != "Bearer google-token" {
				t.Errorf("unexpected authorization header: %q", r.Header.Get("Authorization"))
			}
			return providerResponse(`{
				"id":"google-subject","email":"user@example.com","verified_email":true,
				"given_name":"Test","family_name":"User"
			}`), nil
		default:
			return &stdhttp.Response{
				StatusCode: stdhttp.StatusNotFound, Body: io.NopCloser(strings.NewReader("")),
				Header: make(stdhttp.Header),
			}, nil
		}
	})}

	userUUID := uuid.MustParse("018f0000-0000-7000-8000-000000000020")
	oauthService := &oauthServiceStub{result: &dto.UserResponse{UUID: &userUUID, Email: "user@example.com"}}
	sessionService := &oauthSessionServiceStub{result: &dto.AuthSession{
		AccessToken: "access", RefreshToken: "refresh", CSRFToken: "csrf",
	}}
	oauth := NewGoogleOAuth2(
		&entity.OAuth2Config{GoogleClientID: "client", GoogleClientSecret: "secret"},
		&entity.JWTConfig{AccessTokenMaxAge: 1, RefreshTokenMaxAge: 2},
		oauthService,
		sessionService,
	)
	oauth.GoogleLoginConfig.Endpoint = oauth2.Endpoint{TokenURL: "https://provider.test/token"}
	oauth.userInfoEndpoint = "https://provider.test/userinfo"
	oauth.httpClient = providerClient

	app := fiber.New()
	app.Get("/callback", oauth.Callback)
	req := httptest.NewRequest(stdhttp.MethodGet, "/callback?state=expected&code=code", nil)
	req.AddCookie(&stdhttp.Cookie{Name: oauthStateCookie, Value: "expected"})
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("execute callback: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != stdhttp.StatusOK {
		t.Fatalf("expected status %d, got %d", stdhttp.StatusOK, resp.StatusCode)
	}
	if oauthService.profile.Provider != entity.OAuthProviderGoogle ||
		oauthService.profile.Subject != "google-subject" || !oauthService.profile.EmailVerified {
		t.Fatalf("unexpected typed profile: %#v", oauthService.profile)
	}
	if sessionService.userUUID != userUUID.String() {
		t.Fatalf("session issued for %q", sessionService.userUUID)
	}
	cookies := resp.Cookies()
	for _, name := range []string{"access_token", "refresh_token", "csrf_token"} {
		if !hasCookie(cookies, name) {
			t.Fatalf("missing %s cookie", name)
		}
	}
	var body map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode callback response: %v", err)
	}
	for _, name := range []string{"access_token", "refresh_token", "csrf_token"} {
		if _, exposed := body[name]; exposed {
			t.Fatalf("callback response exposed %s", name)
		}
	}
}

func TestCallbackReturnsIdentityConflict(t *testing.T) {
	providerClient := &stdhttp.Client{Transport: roundTripFunc(func(r *stdhttp.Request) (*stdhttp.Response, error) {
		if r.URL.Path == "/token" {
			return providerResponse(`{"access_token":"google-token","token_type":"Bearer"}`), nil
		}
		return providerResponse(`{
			"id":"different-subject","email":"user@example.com","verified_email":true
		}`), nil
	})}
	oauth := NewGoogleOAuth2(
		&entity.OAuth2Config{GoogleClientID: "client", GoogleClientSecret: "secret"},
		&entity.JWTConfig{},
		&oauthServiceStub{err: apperror.ErrOAuthIdentityConflict},
		&oauthSessionServiceStub{},
	)
	oauth.GoogleLoginConfig.Endpoint = oauth2.Endpoint{TokenURL: "https://provider.test/token"}
	oauth.userInfoEndpoint, oauth.httpClient = "https://provider.test/userinfo", providerClient
	app := fiber.New()
	app.Get("/callback", oauth.Callback)
	req := httptest.NewRequest(stdhttp.MethodGet, "/callback?state=expected&code=code", nil)
	req.AddCookie(&stdhttp.Cookie{Name: oauthStateCookie, Value: "expected"})
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("execute callback: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != stdhttp.StatusConflict {
		t.Fatalf("expected status %d, got %d", stdhttp.StatusConflict, resp.StatusCode)
	}
}

type oauthServiceStub struct {
	profile dto.OAuthProfile
	result  *dto.UserResponse
	err     error
}

func (s *oauthServiceStub) Authenticate(_ context.Context, profile dto.OAuthProfile) (*dto.UserResponse, error) {
	s.profile = profile
	return s.result, s.err
}

type oauthSessionServiceStub struct {
	userUUID string
	result   *dto.AuthSession
	err      error
}

func (s *oauthSessionServiceStub) Issue(_ context.Context, userUUID string) (*dto.AuthSession, error) {
	s.userUUID = userUUID
	return s.result, s.err
}

func hasCookie(cookies []*stdhttp.Cookie, name string) bool {
	for _, cookie := range cookies {
		if cookie.Name == name {
			return true
		}
	}
	return false
}

type roundTripFunc func(*stdhttp.Request) (*stdhttp.Response, error)

func (f roundTripFunc) RoundTrip(request *stdhttp.Request) (*stdhttp.Response, error) {
	return f(request)
}

func providerResponse(body string) *stdhttp.Response {
	return &stdhttp.Response{
		StatusCode: stdhttp.StatusOK,
		Header:     stdhttp.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func requestOAuthState(t *testing.T, app *fiber.App) string {
	t.Helper()

	resp, err := app.Test(httptest.NewRequest(stdhttp.MethodGet, "/login", nil))
	if err != nil {
		t.Fatalf("execute login request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != stdhttp.StatusSeeOther {
		t.Fatalf("expected status %d, got %d", stdhttp.StatusSeeOther, resp.StatusCode)
	}

	redirectURL, err := url.Parse(resp.Header.Get("Location"))
	if err != nil {
		t.Fatalf("parse OAuth redirect: %v", err)
	}
	state := redirectURL.Query().Get("state")
	if state == "" {
		t.Fatal("expected OAuth redirect to contain state")
	}

	for _, cookie := range resp.Cookies() {
		if cookie.Name != oauthStateCookie {
			continue
		}
		if cookie.Value != state {
			t.Fatal("expected OAuth cookie state to match redirect state")
		}
		if !cookie.HttpOnly {
			t.Fatal("expected OAuth state cookie to be HTTP-only")
		}
		return state
	}

	t.Fatal("expected OAuth state cookie")
	return ""
}
