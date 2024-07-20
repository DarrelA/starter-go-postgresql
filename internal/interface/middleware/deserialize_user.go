package middlewares

import (
	"strings"

	"github.com/DarrelA/starter-go-postgresql/internal/application/factory"
	dto "github.com/DarrelA/starter-go-postgresql/internal/application/repository/dto"
	errConst "github.com/DarrelA/starter-go-postgresql/internal/domain/error"
	r "github.com/DarrelA/starter-go-postgresql/internal/domain/repository/redis"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/service"
	"github.com/gofiber/fiber/v2"
)

func Deserializer(
	r r.RedisUserRepository,
	ts service.TokenService,
	uf factory.UserFactory,
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
			return c.Status(fiber.StatusUnauthorized).
				JSON(fiber.Map{"status": "fail", "message": errConst.ErrMsgPleaseLoginAgain})
		}

		tokenClaims, err := ts.ValidateToken(access_token, uf.GetJWTConfig().AccessTokenPublicKey)
		if err != nil {
			return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "error": err})
		}

		userUuid, errGetTokenUUID := r.GetUserUUID(tokenClaims.TokenUUID)
		if errGetTokenUUID != nil {
			return c.Status(fiber.StatusForbidden).
				JSON(fiber.Map{"status": "fail", "message": errConst.ErrMsgPleaseLoginAgain})
		}

		u, err := uf.GetUserByUUID(userUuid)
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

		c.Locals("user_record", userRecord)
		c.Locals("access_token_uuid", tokenClaims.TokenUUID)

		return c.Next()
	}
}
