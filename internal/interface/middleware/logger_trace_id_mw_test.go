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
		requestID, ok := c.Locals("requestID").(string)
		if !ok {
			err := restInterfaceErr.NewBadRequestError(errConst.ErrTypeError + ": requestID")
			log.Error().Err(err).Msg("")
		}

		correlationID, ok := c.Locals("correlationID").(string)
		if !ok {
			err := restInterfaceErr.NewBadRequestError(errConst.ErrTypeError + ": correlationID")
			log.Error().Err(err).Msg("")
		}

		resp := map[string]string{
			"requestID":     requestID,
			"correlationID": correlationID,
		}

		return c.JSON(resp)
	})

	req := httptest.NewRequest("GET", "/id", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("CorrelationAndRequestID middleware test failed: %v", err)
	}

	defer resp.Body.Close()

	var respBody map[string]interface{}
	decodeErr := json.NewDecoder(resp.Body).Decode(&respBody)
	if decodeErr != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	requestID, ok := respBody["requestID"].(string)
	if !ok {
		t.Fatalf("Failed to get requestID from response body")
	}
	if requestID == "" {
		t.Fatal("Expected requestID but got an empty string")
	}

	correlationID, ok := respBody["correlationID"].(string)
	if !ok {
		t.Fatalf("Failed to get correlationID from response body")
	}
	if correlationID == "" {
		t.Fatal("Expected correlationID but got an empty string")
	}
}
