package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	dto "github.com/DarrelA/starter-go-postgresql/internal/application/dto"
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	UserRecord      *dto.UserRecord `json:"userRecord"`
	AccessTokenUUID string          `json:"accessTokenUUID"`
}

func TestDeserializer(t *testing.T) {
	mockUUIDs := mockUUIDs{}
	mockUUIDs.initializeMockUUIDEntities()

	redisUserRepo := &mockRedisUserRepository{mid: mockUUIDs}
	tokenService := &mockTokenService{mid: mockUUIDs}
	userService := &mockUserService{}

	app := fiber.New()
	app.Use(Deserializer(redisUserRepo, tokenService, userService))
	app.Get("/", func(c *fiber.Ctx) error {
		userRecord := c.Locals("userRecord").(*dto.UserRecord)
		accessTokenUUID := c.Locals("accessTokenUUID").(string)

		resp := Response{
			UserRecord:      userRecord,
			AccessTokenUUID: accessTokenUUID,
		}

		return c.JSON(resp)
	})

	for _, test := range deserializerTests {
		t.Run(test.name, func(t *testing.T) {
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
				t.Errorf("Deserializer middleware test failed: %v", err)
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

				if respBody.UserRecord == nil {
					t.Error("Expected userRecord but got nil")
				}

				if respBody.AccessTokenUUID == "" {
					t.Error("Expected accessTokenUUID but got an empty string")
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
