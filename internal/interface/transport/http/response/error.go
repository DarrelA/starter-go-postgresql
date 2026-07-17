package response

import (
	"errors"
	"net/http"

	"github.com/DarrelA/starter-go-postgresql/internal/domain/apperror"
	resterr "github.com/DarrelA/starter-go-postgresql/internal/error"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// Error maps an application error to a safe HTTP error response.
func Error(c *fiber.Ctx, err error) error {
	status, message := http.StatusInternalServerError, resterr.ErrMsgSomethingWentWrong

	switch {
	case errors.Is(err, apperror.ErrEmailConflict):
		status, message = http.StatusConflict, resterr.ErrMsgEmailIsAlreadyTaken
	case errors.Is(err, apperror.ErrOAuthIdentityConflict):
		status, message = http.StatusConflict, resterr.ErrMsgOAuthIdentityConflict
	case errors.Is(err, apperror.ErrInvalidCredentials):
		status, message = http.StatusUnauthorized, resterr.ErrMsgInvalidCredentials
	case errors.Is(err, apperror.ErrInvalidToken),
		errors.Is(err, apperror.ErrOAuthProfileInvalid),
		errors.Is(err, apperror.ErrSessionNotFound),
		errors.Is(err, apperror.ErrSessionRevoked),
		errors.Is(err, apperror.ErrRefreshTokenReused),
		errors.Is(err, apperror.ErrUserNotFound):
		status, message = http.StatusUnauthorized, resterr.ErrMsgPleaseLoginAgain
	case errors.Is(err, apperror.ErrInvalidCSRFToken):
		status, message = http.StatusForbidden, resterr.ErrMsgInvalidCSRFToken
	}

	if status >= http.StatusInternalServerError {
		log.Error().Err(err).Msg("request failed")
	}

	return c.Status(status).JSON(fiber.Map{
		"status": "fail",
		"error":  &resterr.RestErr{Status: status, Message: message},
	})
}
