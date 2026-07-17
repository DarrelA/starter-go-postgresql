package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Response struct {
	Identity RequestIdentity `json:"identity"`
}

func TestAuthenticate(t *testing.T) {
	mockUUIDs := mockUUIDs{}
	mockUUIDs.initializeMockUUIDEntities()

	tokenService := &mockTokenService{mid: mockUUIDs}

	for _, test := range deserializerTests {
		t.Run(test.name, func(t *testing.T) {
			redisUserRepo := &mockRedisUserRepository{mid: mockUUIDs}
			app := fiber.New()
			app.Use(Authenticate(redisUserRepo, tokenService))
			app.Get("/", func(c *fiber.Ctx) error {
				identity, ok := RequestIdentityFrom(c)
				if !ok {
					return c.SendStatus(fiber.StatusInternalServerError)
				}
				return c.JSON(Response{Identity: identity})
			})

			req := httptest.NewRequest("GET", "/", nil)

			// Create a request with the Authorization header
			if test.header != "" {
				req.Header.Set("Authorization", test.header)
			}

			// Create a request with an attached cookie
			if test.cookieName != "" && test.cookieValue != "" {
				req.AddCookie(&http.Cookie{
					Name:  test.cookieName,
					Value: test.cookieValue,
				})
			}

			resp, err := app.Test(req)
			if err != nil {
				t.Errorf("authentication middleware test failed: %v", err)
			}

			defer resp.Body.Close()

			if !test.hasError {
				// Ensure response status code is 200 OK
				if resp.StatusCode != fiber.StatusOK {
					t.Errorf("Expected status '%d' but got '%d'", fiber.StatusOK, resp.StatusCode)
				}

				var respBody Response
				decodeErr := json.NewDecoder(resp.Body).Decode(&respBody)
				if decodeErr != nil {
					t.Errorf("Failed to decode response body: %v", decodeErr)
				}

				if respBody.Identity.UserUUID != *mockUUIDs.mockUserUUID {
					t.Fatalf("expected identity %s, got %s", mockUUIDs.mockUserUUID, respBody.Identity.UserUUID)
				}
			}

			if test.hasError {
				// Ensure response status code is 401 Unauthorized
				if resp.StatusCode != fiber.StatusUnauthorized {
					t.Errorf("Expected status '%d' but got '%d'", fiber.StatusUnauthorized, resp.StatusCode)
				}

				var respBody map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&respBody)
				if err != nil {
					t.Errorf("Failed to decode response body: %v", err)
				}

				errorMsg, ok := respBody["error"].(map[string]interface{})["message"].(string)
				if !ok || !strings.Contains(errorMsg, test.expectedErrMsg) {
					t.Errorf("Expected error message to contain %q, got %q", test.expectedErrMsg, errorMsg)
				}
			}
		})
	}
}

func TestAuthenticateRejectsMismatchedSessionIdentity(t *testing.T) {
	claims := mockUUIDs{}
	claims.initializeMockUUIDEntities()
	stored := claims
	otherUserUUID, err := uuid.NewV7()
	if err != nil {
		t.Fatalf("create mismatched user UUID: %v", err)
	}
	stored.mockUserUUID = &otherUserUUID

	repository := &mockRedisUserRepository{mid: stored}
	app := fiber.New()
	app.Use(Authenticate(repository, &mockTokenService{mid: claims}))
	called := false
	app.Get("/", func(c *fiber.Ctx) error {
		called = true
		return c.SendStatus(fiber.StatusOK)
	})
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.AddCookie(&http.Cookie{Name: "access_token", Value: "mockAccessToken"})
	response, err := app.Test(request)
	if err != nil {
		t.Fatalf("execute request: %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, response.StatusCode)
	}
	if called {
		t.Fatal("mismatched session identity reached the protected handler")
	}
}
