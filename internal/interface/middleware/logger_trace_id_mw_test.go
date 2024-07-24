package middleware

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	errConst "github.com/DarrelA/starter-go-postgresql/internal/domain/error"
	restInterfaceErr "github.com/DarrelA/starter-go-postgresql/internal/interface/transport/http/error"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func TestCorrelationAndRequestID(t *testing.T) {
	app := fiber.New()
	app.Get("/id", CorrelationAndRequestID, func(c *fiber.Ctx) error {
		request_id, ok := c.Locals("request_id").(string)
		if !ok {
			err := restInterfaceErr.NewBadRequestError(errConst.ErrTypeError + ": request_id")
			log.Error().Err(err).Msg("")
		}

		correlation_id, ok := c.Locals("correlation_id").(string)
		if !ok {
			err := restInterfaceErr.NewBadRequestError(errConst.ErrTypeError + ": correlation_id")
			log.Error().Err(err).Msg("")
		}

		resp := map[string]string{
			"requestID":     request_id,
			"correlationID": correlation_id,
		}

		return c.JSON(resp)
	})

	req := httptest.NewRequest("GET", "/id", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("CorrelationAndRequestID middleware test failed: %v", err)
	}

	defer resp.Body.Close()

	var responseBody map[string]interface{}
	decodeErr := json.NewDecoder(resp.Body).Decode(&responseBody)
	if decodeErr != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	requestID, ok := responseBody["requestID"].(string)
	if !ok {
		t.Fatalf("Failed to get requestID from response body")
	}
	if requestID == "" {
		t.Fatal("Expected requestID but got an empty string")
	}

	correlationID, ok := responseBody["correlationID"].(string)
	if !ok {
		t.Fatalf("Failed to get correlationID from response body")
	}
	if correlationID == "" {
		t.Fatal("Expected correlationID but got an empty string")
	}
}
