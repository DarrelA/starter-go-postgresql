package service

import (
	"time"

	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	"github.com/DarrelA/starter-go-postgresql/internal/utils/err_rest"
	"github.com/gofiber/fiber/v2"
)

type AuthService interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	RefreshAccessToken(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
}

type TokenService interface {
	CreateToken(user_uuid string, ttl time.Duration, privateKey string) (*entity.Token, *err_rest.RestErr)
	ValidateToken(token string, publicKey string) (*entity.Token, *err_rest.RestErr)
}
