package middlewares

import (
	"net"
	"os"
	"time"

	"github.com/DarrelA/starter-go-postgresql/internal/utils/err_rest"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func LoggerMW(c *fiber.Ctx) error {
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

	currentEnv, ok := c.Locals("env").(string)
	if !ok {
		log.Panic().Msg("invalid_currentEnv")
	}

	switch currentEnv {
	case "prod":
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
			Str("ip_version", getIPVersion(ip)).
			Str("user_agent", c.Get("User-Agent")).
			Str("correlation_id", correlation_id).
			Str("request_id", request_id).
			Msg("request is completed")

	case "dev":
		log.Info().
			Str("method", c.Method()).
			Str("url", c.OriginalURL()).
			Int("status", c.Response().StatusCode()).
			Bytes("request_body", c.Request().Body()).
			Bytes("response_body", c.Response().Body()).
			Msgf("request is completed in [%s] env", currentEnv)

	case "test":
		log.Info().
			Str("hostname", hostname).
			Str("method", c.Method()).
			Str("referer", c.Get("Referer")).
			Str("url", c.OriginalURL()).
			Int("status", c.Response().StatusCode()).
			Bytes("request_body", c.Request().Body()).
			Bytes("response_body", c.Response().Body()).
			Str("user_host", c.Get("Host")).
			Dur("response_time", duration).
			Int64("latency_ms", duration.Milliseconds()).
			Str("ip_address", ip).
			Str("ip_version", getIPVersion(ip)).
			Str("user_agent", c.Get("User-Agent")).
			Str("correlation_id", correlation_id).
			Str("request_id", request_id).
			Msgf("request is completed in [%s] env", currentEnv)

	default:
		log.Info().Msgf("expecting prod, dev or test env but, current the env is [%s]", currentEnv)
	}

	return err
}

func getIPVersion(ip string) string {
	parsedIP := net.ParseIP(ip)

	if parsedIP == nil {
		return "unknown"
	}
	if parsedIP.To4() != nil {
		return "IPv4"
	}

	return "IPv6"
}
