package middlewares

import (
	"os"
	"time"

	"github.com/DarrelA/starter-go-postgresql/internal/utils"
	"github.com/DarrelA/starter-go-postgresql/internal/utils/err_rest"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func LogRequest(c *fiber.Ctx) error {
	start := time.Now()
	err := c.Next()
	duration := time.Since(start)
	ip := c.IP()

	request_id, ok := c.Locals("request_id").(string)
	if !ok {
		err := err_rest.NewBadRequestError("request_id is not a string or not set")
		log.Error().Err(err).Msg("")
	}

	correlation_id, ok := c.Locals("correlation_id").(string)
	if !ok {
		err := err_rest.NewBadRequestError("correlation_id is not a string or not set")
		log.Error().Err(err).Msg("")
	}

	hostname, hostnameErr := os.Hostname()
	if hostnameErr != nil {
		log.Fatal().Err(hostnameErr).Msg("failed to get hostname")
	}

	log.Info().
		Str("hostname", hostname).
		Str("method", c.Method()).
		Str("referer", c.Get("Referer")).
		Str("url", c.OriginalURL()).
		Int("status", c.Response().StatusCode()).
		Str("user_host", c.Get("Host")).
		Dur("response_time", duration).
		Int64("latency_ms", duration.Milliseconds()).
		Str("ip_address", ip).
		Str("ip_version", utils.GetIPVersion(ip)).
		Str("user_agent", c.Get("User-Agent")).
		Str("correlation_id", correlation_id).
		Str("request_id", request_id).
		Msg("Request is completed")

	return err
}
