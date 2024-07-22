package middleware

import (
	errConst "github.com/DarrelA/starter-go-postgresql/internal/domain/error"
	restInterfaceErr "github.com/DarrelA/starter-go-postgresql/internal/interface/transport/http/error"
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
			log.Error().Err(err).Msg(errConst.ErrUUIDError)
			return restInterfaceErr.NewUnprocessableEntityError(errConst.ErrUUIDError)
		}

		correlationID = id.String()
		c.Set(correlation_id, correlationID)
	}

	if requestID == "" {
		id, err := uuid.NewV7()
		if err != nil {
			log.Error().Err(err).Msg(errConst.ErrUUIDError)
			return restInterfaceErr.NewUnprocessableEntityError(errConst.ErrUUIDError)
		}

		requestID = id.String()
		c.Set(request_id, requestID)
	}

	c.Locals("correlation_id", correlationID)
	c.Locals("request_id", requestID)

	return c.Next()
}
