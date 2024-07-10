package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

const correlation_id = "Correlation-ID"
const request_id = "Request-ID"

func CorrelationAndRequestID(c *fiber.Ctx) error {
	correlationID := c.Get(correlation_id)
	requestID := c.Get(request_id)

	if correlationID == "" {
		id, err := uuid.NewV7()
		if err != nil {
			log.Error().Err(err).Msg("correlationID_uuid_error: ")
			return fiber.ErrUnprocessableEntity
		}

		correlationID = id.String()
		c.Set(correlation_id, correlationID)
	}

	if requestID == "" {
		id, err := uuid.NewV7()
		if err != nil {
			log.Error().Err(err).Msg("requestID_uuid_error")
			return fiber.ErrUnprocessableEntity
		}

		requestID = id.String()
		c.Set(request_id, requestID)
	}

	c.Locals("correlation_id", correlationID)
	c.Locals("request_id", requestID)

	return c.Next()
}
