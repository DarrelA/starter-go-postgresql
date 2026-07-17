package http

import (
	stdhttp "net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestAuthHandlersRejectMissingValidatedPayload(t *testing.T) {
	app := fiber.New()
	auth := &AuthHandler{}
	app.Post("/register", auth.Register)
	app.Post("/login", auth.Login)

	for _, path := range []string{"/register", "/login"} {
		t.Run(path, func(t *testing.T) {
			resp, err := app.Test(httptest.NewRequest(stdhttp.MethodPost, path, nil))
			if err != nil {
				t.Fatalf("execute request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != stdhttp.StatusBadRequest {
				t.Fatalf("expected status %d, got %d", stdhttp.StatusBadRequest, resp.StatusCode)
			}
		})
	}
}
