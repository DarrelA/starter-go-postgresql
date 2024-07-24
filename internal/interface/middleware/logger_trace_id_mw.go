package middleware

import (
	errConst "github.com/DarrelA/starter-go-postgresql/internal/domain/error"
	restInterfaceErr "github.com/DarrelA/starter-go-postgresql/internal/interface/transport/http/error"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

const correlationID = "Correlation-ID"
const requestID = "Request-ID"

func CorrelationAndRequestID(c *fiber.Ctx) error {
	correlationID := c.Get(correlationID)
	requestID := c.Get(requestID)

	if correlationID == "" {
		id, err := uuid.NewV7()
		if err != nil { // coverage:ignore
			log.Error().Err(err).Msg(errConst.ErrUUIDError)
			return restInterfaceErr.NewUnprocessableEntityError(errConst.ErrUUIDError)
		}

		correlationID = id.String()
		c.Set(correlationID, correlationID)
	}

	if requestID == "" {
		id, err := uuid.NewV7()
		if err != nil { // coverage:ignore
			log.Error().Err(err).Msg(errConst.ErrUUIDError)
			return restInterfaceErr.NewUnprocessableEntityError(errConst.ErrUUIDError)
		}

		requestID = id.String()
		c.Set(requestID, requestID)
	}

	c.Locals("correlationID", correlationID)
	c.Locals("requestID", requestID)

	return c.Next()
}
