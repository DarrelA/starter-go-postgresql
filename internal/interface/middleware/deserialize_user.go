package middlewares

import (
	"context"
	"strings"

	"github.com/DarrelA/starter-go-postgresql/configs"
	redisDb "github.com/DarrelA/starter-go-postgresql/db/redis"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/factory"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/service"
	dto "github.com/DarrelA/starter-go-postgresql/internal/interface/transport/dto"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

var jwtCfg = configs.JWTSettings

func Deserializer(token service.TokenService, uf factory.UserFactory) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var access_token string
		authorization := c.Get("Authorization")

		if strings.HasPrefix(authorization, "Bearer ") {
			access_token = strings.TrimPrefix(authorization, "Bearer ")
		} else if c.Cookies("access_token") != "" {
			access_token = c.Cookies("access_token")
		}

		message := "please login again"

		if access_token == "" {
			return c.Status(fiber.StatusUnauthorized).
				JSON(fiber.Map{"status": "fail", "message": message})
		}

		tokenClaims, err := token.ValidateToken(access_token, jwtCfg.AccessTokenPublicKey)
		if err != nil {
			return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "error": err})
		}

		ctx := context.TODO()

		user_uuid, errGetTokenUUID := redisDb.RedisClient.Get(ctx, tokenClaims.TokenUUID).Result()
		if errGetTokenUUID == redis.Nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "message": message})
		}

		u, err := uf.GetUserByUUID(user_uuid)
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
