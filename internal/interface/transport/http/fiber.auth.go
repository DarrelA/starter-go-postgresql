package http

import (
	"context"
	"time"

	"github.com/DarrelA/starter-go-postgresql/configs"
	redisDb "github.com/DarrelA/starter-go-postgresql/db/redis"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/factory"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/service"
	dto "github.com/DarrelA/starter-go-postgresql/internal/interface/transport/dto"
	"github.com/DarrelA/starter-go-postgresql/internal/utils/err_rest"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type AuthHandler struct {
	uf factory.UserFactory
	ts service.TokenService
}

func NewAuthService(uf factory.UserFactory, ts service.TokenService) *AuthHandler {
	return &AuthHandler{uf, ts}
}

var jwtCfg = configs.JWTSettings

func (ah *AuthHandler) Register(c *fiber.Ctx) error {
	payload, ok := c.Locals("register_payload").(dto.RegisterInput)
	if !ok {
		err := err_rest.NewBadRequestError("register_payload is not of type users.RegisterInput")
		log.Error().Err(err).Msg("type_error")
	}

	result, err := ah.uf.CreateUser(payload)
	if err != nil {
		return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "error": err})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "user": result})
}

func (ah *AuthHandler) Login(c *fiber.Ctx) error {
	payload, ok := c.Locals("login_payload").(dto.LoginInput)
	if !ok {
		err := err_rest.NewBadRequestError("login_payload is not of type users.RegisterInput")
		log.Error().Err(err).Msg("type_error")
	}

	user, err := ah.uf.GetUser(payload)
	if err != nil {
		return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "error": err})
	}

	accessTokenDetails, err := ah.ts.CreateToken(
		user.UUID.String(),
		jwtCfg.AccessTokenExpiredIn,
		jwtCfg.AccessTokenPrivateKey,
	)
	if err != nil {
		return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "error": err})
	}

	refreshTokenDetails, err := ah.ts.CreateToken(
		user.UUID.String(),
		jwtCfg.RefreshTokenExpiredIn,
		jwtCfg.RefreshTokenPrivateKey,
	)
	if err != nil {
		return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "error": err})
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
		log.Error().Err(errAccess).Msg("redis_error")
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"status": "fail", "message": "something went wrong"})
	}

	errRefresh := redisDb.RedisClient.Set(
		ctx,
		refreshTokenDetails.TokenUUID,
		user.UUID.String(),
		time.Unix(*refreshTokenDetails.ExpiresIn, 0).Sub(timeNow),
	).Err()

	if errRefresh != nil {
		log.Error().Err(errRefresh).Msg("redis_error")
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"status": "fail", "message": "something went wrong"})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    *accessTokenDetails.Token,
		Path:     "/",
		Domain:   jwtCfg.Domain,
		MaxAge:   jwtCfg.AccessTokenMaxAge * 60,
		Secure:   jwtCfg.Secure,
		HTTPOnly: jwtCfg.HttpOnly,
		SameSite: "strict",
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    *refreshTokenDetails.Token,
		Path:     "/",
		Domain:   jwtCfg.Domain,
		MaxAge:   jwtCfg.RefreshTokenMaxAge * 60,
		Secure:   jwtCfg.Secure,
		HTTPOnly: jwtCfg.HttpOnly,
		SameSite: "strict",
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "access_token": accessTokenDetails.Token})
}

func (ah *AuthHandler) RefreshAccessToken(c *fiber.Ctx) error {
	message := "please login again"
	refresh_token := c.Cookies("refresh_token")

	if refresh_token == "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "message": message})
	}

	ctx := context.TODO()

	tokenClaims, err := ah.ts.ValidateToken(refresh_token, jwtCfg.RefreshTokenPublicKey)
	if err != nil {
		return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "error": err})
	}

	user_uuid, errGetTokenUUID := redisDb.RedisClient.Get(ctx, tokenClaims.TokenUUID).Result()
	if errGetTokenUUID == redis.Nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "message": message})
	}

	user, err := ah.uf.GetUserByUUID(user_uuid)
	if err != nil {
		return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "error": err})
	}

	accessTokenDetails, err := ah.ts.CreateToken(
		user.UUID.String(),
		jwtCfg.AccessTokenExpiredIn,
		jwtCfg.AccessTokenPrivateKey,
	)
	if err != nil {
		return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "error": err})
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
		Domain:   jwtCfg.Domain,
		MaxAge:   jwtCfg.AccessTokenMaxAge * 60,
		Secure:   jwtCfg.Secure,
		HTTPOnly: jwtCfg.HttpOnly,
		SameSite: "strict",
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "access_token": accessTokenDetails.Token})
}

func (ah *AuthHandler) Logout(c *fiber.Ctx) error {
	message := "please login again"
	refresh_token := c.Cookies("refresh_token")

	if refresh_token == "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "message": message})
	}

	ctx := context.TODO()

	tokenClaims, err := ah.ts.ValidateToken(refresh_token, jwtCfg.RefreshTokenPublicKey)
	if err != nil {
		return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "error": err})
	}

	accessTokenUUID, ok := c.Locals("access_token_uuid").(string) // type assertion
	if !ok {
		err := err_rest.NewBadRequestError("access_token is not a string or not set")
		log.Error().Err(err).Msg("")
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

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}
