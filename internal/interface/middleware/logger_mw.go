package middleware

import (
	"net"
	"os"
	"time"

	errConst "github.com/DarrelA/starter-go-postgresql/internal/domain/error"
	restInterfaceErr "github.com/DarrelA/starter-go-postgresql/internal/interface/transport/http/error"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

const errMsgInvalidHostname = "failed to get hostname"

func LoggerMW(c *fiber.Ctx) error {
	start := time.Now()
	err := c.Next()
	duration := time.Since(start)
	ip := c.IP()

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

	hostname, hostnameErr := os.Hostname()
	if hostnameErr != nil {
		log.Error().Err(hostnameErr).Msg(errMsgInvalidHostname) // coverage:ignore
	}

	currentEnv, ok := c.Locals("env").(string)
	if !ok {
		err := restInterfaceErr.NewBadRequestError(errConst.ErrTypeError + ": currentEnv")
		log.Error().Err(err).Msg("")
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
			Str("correlationID", correlationID).
			Str("requestID", requestID).
			Msgf("request is completed in [%s] env", currentEnv)

	case "dev":
		log.Info().
			Str("method", c.Method()).
			Str("url", c.OriginalURL()).
			Int("status", c.Response().StatusCode()).
			Bytes("request_body", c.Request().Body()).
			Bytes("response_body", c.Response().Body()).
			Msgf("request is completed in [%s] env", currentEnv)

	default:
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
			Str("correlationID", correlationID).
			Str("requestID", requestID).
			Msgf("request is completed in [%s] env", currentEnv)
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
