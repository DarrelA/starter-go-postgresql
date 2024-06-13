package handlers

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

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
		return c.Status(err.Status).JSON(err)
	}

	result, err := services.CreateUser(payload)
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func Login(c *fiber.Ctx) error {
	var payload users.LoginInput

	if err := c.BodyParser(&payload); err != nil {
		err := errors.NewBadRequestError("invalid json body")
		return c.Status(err.Status).JSON(err)
	}

	user, err := services.GetUser(payload)
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	accessTokenDetails, err := services.CreateToken(
		user.UUID.String(),
		jwtCfg.AccessTokenExpiredIn,
		jwtCfg.AccessTokenPrivateKey,
	)
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	refreshTokenDetails, err := services.CreateToken(
		user.UUID.String(),
		jwtCfg.RefreshTokenExpiredIn,
		jwtCfg.RefreshTokenPrivateKey,
	)
	if err != nil {
		return c.Status(err.Status).JSON(err)
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
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"message": errAccess.Error()})
	}

	errRefresh := redisDb.RedisClient.Set(
		ctx,
		refreshTokenDetails.TokenUUID,
		user.UUID.String(),
		time.Unix(*refreshTokenDetails.ExpiresIn, 0).Sub(timeNow),
	).Err()

	if errRefresh != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"message": errRefresh.Error()})
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
	message := "could not refresh access token"
	refresh_token := c.Cookies("refresh_token")

	if refresh_token == "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "message": message})
	}

	ctx := context.TODO()

	tokenClaims, err := services.ValidateToken(refresh_token, jwtCfg.RefreshTokenPublicKey)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	user_uuid, err := redisDb.RedisClient.Get(ctx, tokenClaims.TokenUUID).Result()
	if err == redis.Nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "message": message})
	}

	// @TODO: Streamline custom error handling and Fiber error handling
	user, restErr := services.GetUserByUUID(user_uuid)
	if restErr != nil {
		return c.Status(restErr.Status).JSON(restErr)
	}

	accessTokenDetails, restErr := services.CreateToken(
		user.UUID.String(),
		jwtCfg.AccessTokenExpiredIn,
		jwtCfg.AccessTokenPrivateKey,
	)
	if restErr != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	timeNow := time.Now()

	errAccess := redisDb.RedisClient.Set(
		ctx,
		accessTokenDetails.TokenUUID,
		user.UUID.String(),
		time.Unix(*accessTokenDetails.ExpiresIn, 0).Sub(timeNow),
	).Err()

	if errAccess != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"message": errAccess.Error()})
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
