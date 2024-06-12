package users

import (
	"context"
	"time"

	"github.com/DarrelA/starter-go-postgresql/configs"
	redis "github.com/DarrelA/starter-go-postgresql/db/redis"
	"github.com/DarrelA/starter-go-postgresql/internal/domains/users"
	"github.com/DarrelA/starter-go-postgresql/internal/services"
	"github.com/DarrelA/starter-go-postgresql/internal/utils"
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

	accessTokenDetails, err := utils.CreateToken(
		user.UUID.String(),
		jwtCfg.AccessTokenExpiredIn,
		jwtCfg.AccessTokenPrivateKey,
	)
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	refreshTokenDetails, err := utils.CreateToken(
		user.UUID.String(),
		jwtCfg.RefreshTokenExpiredIn,
		jwtCfg.RefreshTokenPrivateKey,
	)
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	ctx := context.TODO()
	timeNow := time.Now()

	errAccess := redis.RedisClient.Set(
		ctx,
		accessTokenDetails.TokenUUID,
		user.UUID.String(),
		time.Unix(*accessTokenDetails.ExpiresIn, 0).Sub(timeNow),
	).Err()

	if errAccess != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"message": errAccess.Error()})
	}

	errRefresh := redis.RedisClient.Set(
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
		Secure:   jwtCfg.Secure,
		HTTPOnly: jwtCfg.HttpOnly,
		Domain:   jwtCfg.Domain,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "access_token": accessTokenDetails.Token})
}
