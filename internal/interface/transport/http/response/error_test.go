package response

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DarrelA/starter-go-postgresql/internal/domain/apperror"
	"github.com/gofiber/fiber/v2"
)

func TestErrorMapsDomainErrorsAtHTTPBoundary(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		status int
	}{
		{name: "email conflict", err: apperror.ErrEmailConflict, status: http.StatusConflict},
		{name: "wrapped invalid token", err: errors.Join(apperror.ErrInvalidToken, errors.New("signature failed")), status: http.StatusUnauthorized},
		{name: "refresh token reused", err: apperror.ErrRefreshTokenReused, status: http.StatusUnauthorized},
		{name: "invalid CSRF token", err: apperror.ErrInvalidCSRFToken, status: http.StatusForbidden},
		{name: "invalid OAuth profile", err: apperror.ErrOAuthProfileInvalid, status: http.StatusUnauthorized},
		{name: "OAuth identity conflict", err: apperror.ErrOAuthIdentityConflict, status: http.StatusConflict},
		{name: "unknown infrastructure error", err: errors.New("database unavailable"), status: http.StatusInternalServerError},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			app := fiber.New()
			app.Get("/", func(c *fiber.Ctx) error { return Error(c, test.err) })

			response, err := app.Test(httptest.NewRequest(http.MethodGet, "/", nil))
			if err != nil {
				t.Fatalf("execute request: %v", err)
			}
			defer response.Body.Close()
			if response.StatusCode != test.status {
				t.Fatalf("expected status %d, got %d", test.status, response.StatusCode)
			}
		})
	}
}
