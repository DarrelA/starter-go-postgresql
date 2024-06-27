package middlewares

import (
	"context"
	"strings"

	"github.com/DarrelA/starter-go-postgresql/configs"
	redisDb "github.com/DarrelA/starter-go-postgresql/db/redis"
	"github.com/DarrelA/starter-go-postgresql/internal/domains/users"
	"github.com/DarrelA/starter-go-postgresql/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

var jwtCfg = configs.JWTSettings

func Deserializer(c *fiber.Ctx) error {

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

	tokenClaims, err := services.ValidateToken(access_token, jwtCfg.AccessTokenPublicKey)
	if err != nil {
		return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "error": err})
	}

	ctx := context.TODO()

	user_uuid, errGetTokenUUID := redisDb.RedisClient.Get(ctx, tokenClaims.TokenUUID).Result()
	if errGetTokenUUID == redis.Nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "message": message})
	}

	user, err := services.GetUserByUUID(user_uuid)
	if err != nil {
		return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "error": err})
	}

	userRecord := &users.UserRecord{
		UUID:      user.UUID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	c.Locals("user_record", userRecord)
	c.Locals("access_token_uuid", tokenClaims.TokenUUID)

	return c.Next()
}
