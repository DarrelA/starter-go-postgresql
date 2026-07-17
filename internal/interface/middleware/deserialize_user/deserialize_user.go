package middleware

import (
	"strings"

	"github.com/DarrelA/starter-go-postgresql/internal/domain/apperror"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/repository"
	domainSvc "github.com/DarrelA/starter-go-postgresql/internal/domain/service"
	restErr "github.com/DarrelA/starter-go-postgresql/internal/error"
	"github.com/DarrelA/starter-go-postgresql/internal/interface/transport/http/response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Authenticate verifies the access JWT and its Redis-backed session.
func Authenticate(
	r repository.TokenRepository,
	ts domainSvc.TokenService,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var accessToken string
		authorization := c.Get("Authorization")

		if strings.HasPrefix(authorization, "Bearer ") {
			accessToken = strings.TrimPrefix(authorization, "Bearer ")
		} else if c.Cookies("access_token") != "" {
			accessToken = c.Cookies("access_token")
		}

		if accessToken == "" {
			err := restErr.NewUnauthorizedError(restErr.ErrMsgPleaseLoginAgain)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "error": err})
		}

		tokenClaims, err := ts.ValidateAccessToken(accessToken)
		if err != nil {
			return response.Error(c, err)
		}

		storedUserUUID, err := r.GetUserUUID(
			c.UserContext(), tokenClaims.FamilyUUID, tokenClaims.TokenUUID,
		)
		if err != nil {
			return response.Error(c, err)
		}

		userUUID, err := uuid.Parse(storedUserUUID)
		claimsUserUUID, claimsErr := uuid.Parse(tokenClaims.UserUUID)
		if err != nil || claimsErr != nil || userUUID != claimsUserUUID {
			return response.Error(c, apperror.ErrInvalidToken)
		}

		SetRequestIdentity(c, RequestIdentity{UserUUID: userUUID})

		return c.Next()
	}
}
