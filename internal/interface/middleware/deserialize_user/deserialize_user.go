package middleware

import (
	"strings"

	dto "github.com/DarrelA/starter-go-postgresql/internal/application/dto"
	appSvc "github.com/DarrelA/starter-go-postgresql/internal/application/service"
	r "github.com/DarrelA/starter-go-postgresql/internal/domain/repository/redis"
	domainSvc "github.com/DarrelA/starter-go-postgresql/internal/domain/service"
	restErr "github.com/DarrelA/starter-go-postgresql/internal/error"
	"github.com/gofiber/fiber/v2"
)

func Deserializer(
	r r.RedisUserRepository,
	ts domainSvc.TokenService,
	us appSvc.UserService,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var access_token string
		authorization := c.Get("Authorization")

		if strings.HasPrefix(authorization, "Bearer ") {
			access_token = strings.TrimPrefix(authorization, "Bearer ")
		} else if c.Cookies("access_token") != "" {
			access_token = c.Cookies("access_token")
		}

		if access_token == "" {
			err := restErr.NewUnauthorizedError(restErr.ErrMsgPleaseLoginAgain)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "error": err})
		}

		tokenClaims, err := ts.ValidateToken(access_token, us.GetJWTConfig().AccessTokenPublicKey)
		if err != nil {
			return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "error": err})
		}

		userUuid, errGetTokenUUID := r.GetUserUUID(tokenClaims.TokenUUID)
		if errGetTokenUUID != nil {
			return c.Status(fiber.StatusUnauthorized).
				JSON(fiber.Map{"status": "fail", "message": errGetTokenUUID})
		}

		u, err := us.GetUserByUUID(userUuid)
		if err != nil {
			return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "error": err})
		}

		userRecord := &dto.UserRecord{
			UUID:      u.UUID,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Email:     u.Email,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		}

		c.Locals("userRecord", userRecord)
		c.Locals("accessTokenUUID", tokenClaims.TokenUUID)

		return c.Next()
	}
}
