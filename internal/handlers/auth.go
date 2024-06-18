package handlers

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"

	"github.com/DarrelA/starter-go-postgresql/configs"
	redisDb "github.com/DarrelA/starter-go-postgresql/db/redis"
	"github.com/DarrelA/starter-go-postgresql/internal/domains/users"
	"github.com/DarrelA/starter-go-postgresql/internal/services"
	"github.com/DarrelA/starter-go-postgresql/internal/utils/errors"
	"github.com/gofiber/fiber/v2"
)

var jwtCfg = configs.JWTSettings

func Register(c *fiber.Ctx) error {
	var payload users.RegisterInput

	if err := c.BodyParser(&payload); err != nil {
		err := errors.NewBadRequestError("invalid json body")
		return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "message": err.Message})
	}

	result, err := services.CreateUser(payload)
	if err != nil {
		return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "message": err.Message})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "user": result})
}

func Login(c *fiber.Ctx) error {
	var payload users.LoginInput

	if err := c.BodyParser(&payload); err != nil {
		err := errors.NewBadRequestError("invalid json body")
		return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "message": err.Message})
	}

	user, err := services.GetUser(payload)
	if err != nil {
		return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "message": err.Message})
	}

	accessTokenDetails, err := services.CreateToken(
		user.UUID.String(),
		jwtCfg.AccessTokenExpiredIn,
		jwtCfg.AccessTokenPrivateKey,
	)
	if err != nil {
		return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "message": err.Message})
	}

	refreshTokenDetails, err := services.CreateToken(
		user.UUID.String(),
		jwtCfg.RefreshTokenExpiredIn,
		jwtCfg.RefreshTokenPrivateKey,
	)
	if err != nil {
		return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "message": err.Message})
	}

	ctx := context.TODO()
	timeNow := time.Now()

	errAccess := redisDb.RedisClient.Set(
		ctx,
		accessTokenDetails.TokenUUID,
		user.UUID.String(),
		time.Unix(*accessTokenDetails.ExpiresIn, 0).Sub(timeNow),
	).Err()

	if errAccess != nil {
		log.Error().Msg("redis_error: " + errAccess.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"status": "fail", "message": "something went wrong"})
	}

	errRefresh := redisDb.RedisClient.Set(
		ctx,
		refreshTokenDetails.TokenUUID,
		user.UUID.String(),
		time.Unix(*refreshTokenDetails.ExpiresIn, 0).Sub(timeNow),
	).Err()

	if errRefresh != nil {
		log.Error().Msg("redis_error: " + errRefresh.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"status": "fail", "message": "something went wrong"})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    *accessTokenDetails.Token,
		Path:     "/",
		MaxAge:   jwtCfg.AccessTokenMaxAge * 60,
		Secure:   jwtCfg.Secure,
		HTTPOnly: jwtCfg.HttpOnly,
		Domain:   jwtCfg.Domain,
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    *refreshTokenDetails.Token,
		Path:     "/",
		MaxAge:   jwtCfg.RefreshTokenMaxAge * 60,
		Secure:   jwtCfg.Secure,
		HTTPOnly: jwtCfg.HttpOnly,
		Domain:   jwtCfg.Domain,
	})

	c.Cookie(&fiber.Cookie{
		Name:     "logged_in",
		Value:    "true",
		Path:     "/",
		MaxAge:   jwtCfg.AccessTokenMaxAge * 60,
		Secure:   false,
		HTTPOnly: false,
		Domain:   "localhost",
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "access_token": accessTokenDetails.Token})
}

func RefreshAccessToken(c *fiber.Ctx) error {
	message := "please login again"
	refresh_token := c.Cookies("refresh_token")

	if refresh_token == "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "message": message})
	}

	ctx := context.TODO()

	tokenClaims, err := services.ValidateToken(refresh_token, jwtCfg.RefreshTokenPublicKey)
	if err != nil {
		return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "message": err.Message})
	}

	user_uuid, errGetTokenUUID := redisDb.RedisClient.Get(ctx, tokenClaims.TokenUUID).Result()
	if errGetTokenUUID == redis.Nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "message": message})
	}

	user, err := services.GetUserByUUID(user_uuid)
	if err != nil {
		return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "message": err.Message})
	}

	accessTokenDetails, err := services.CreateToken(
		user.UUID.String(),
		jwtCfg.AccessTokenExpiredIn,
		jwtCfg.AccessTokenPrivateKey,
	)
	if err != nil {
		return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "message": err.Message})
	}

	timeNow := time.Now()

	errAccess := redisDb.RedisClient.Set(
		ctx,
		accessTokenDetails.TokenUUID,
		user.UUID.String(),
		time.Unix(*accessTokenDetails.ExpiresIn, 0).Sub(timeNow),
	).Err()

	if errAccess != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"status": "fail", "message": "something went wrong"})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    *accessTokenDetails.Token,
		Path:     "/",
		MaxAge:   jwtCfg.AccessTokenMaxAge * 60,
		Secure:   jwtCfg.Secure,
		HTTPOnly: jwtCfg.HttpOnly,
		Domain:   jwtCfg.Domain,
	})

	c.Cookie(&fiber.Cookie{
		Name:     "logged_in",
		Value:    "true",
		Path:     "/",
		MaxAge:   jwtCfg.AccessTokenMaxAge * 60,
		Secure:   false,
		HTTPOnly: false,
		Domain:   "localhost",
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "access_token": accessTokenDetails.Token})
}

func Logout(c *fiber.Ctx) error {
	message := "please login again"
	refresh_token := c.Cookies("refresh_token")

	if refresh_token == "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "message": message})
	}

	ctx := context.TODO()

	tokenClaims, err := services.ValidateToken(refresh_token, jwtCfg.RefreshTokenPublicKey)
	if err != nil {
		return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "message": err.Message})
	}

	accessTokenUUID, ok := c.Locals("access_token_uuid").(string) // type assertion
	if !ok {
		log.Error().Msg("access_token is not a string or not set")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": message})
	}

	_, errDelTokenUUID := redisDb.RedisClient.Del(ctx, tokenClaims.TokenUUID, accessTokenUUID).Result()
	if errDelTokenUUID != nil {
		return c.Status(fiber.StatusBadGateway).
			JSON(fiber.Map{"status": "fail", "message": errDelTokenUUID.Error()})
	}

	expired := time.Now().Add(-time.Hour * 24)

	c.Cookie(&fiber.Cookie{
		Name:    "access_token",
		Value:   "",
		Expires: expired,
	})

	c.Cookie(&fiber.Cookie{
		Name:    "refresh_token",
		Value:   "",
		Expires: expired,
	})

	c.Cookie(&fiber.Cookie{
		Name:    "logged_in",
		Value:   "",
		Expires: expired,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}
