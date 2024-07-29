package http

import (
	"time"

	"github.com/DarrelA/starter-go-postgresql/internal/application/factory"
	dto "github.com/DarrelA/starter-go-postgresql/internal/application/repository/dto"
	appSvc "github.com/DarrelA/starter-go-postgresql/internal/application/service"
	errConst "github.com/DarrelA/starter-go-postgresql/internal/domain/error"
	r "github.com/DarrelA/starter-go-postgresql/internal/domain/repository/redis"
	domainSvc "github.com/DarrelA/starter-go-postgresql/internal/domain/service"
	restInterfaceErr "github.com/DarrelA/starter-go-postgresql/internal/interface/transport/http/error"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

const (
	errMsgRegisterPayload = "register_payload is not of type users.RegisterInput"
	errMsgLoginPayload    = "login_payload is not of type users.RegisterInput"
	errMsgAccessTokenUUID = "accessTokenUUID is not a string or not set"
)

type AuthService struct {
	r  r.RedisUserRepository
	uf factory.UserFactory
	ts domainSvc.TokenService
}

func NewAuthService(
	r r.RedisUserRepository,
	uf factory.UserFactory,
	ts domainSvc.TokenService,
) appSvc.AuthService {
	return &AuthService{r, uf, ts}
}

func (ah *AuthService) Register(c *fiber.Ctx) error {
	payload, ok := c.Locals("register_payload").(dto.RegisterInput)
	if !ok {
		err := restInterfaceErr.NewBadRequestError(errMsgRegisterPayload)
		log.Error().Err(err).Msg(errConst.ErrTypeError)
	}

	result, err := ah.uf.CreateUser(payload)
	if err != nil {
		return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "error": err})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "user": result})
}

func (ah *AuthService) Login(c *fiber.Ctx) error {
	payload, ok := c.Locals("login_payload").(dto.LoginInput)
	if !ok {
		err := restInterfaceErr.NewBadRequestError(errMsgLoginPayload)
		log.Error().Err(err).Msg(errConst.ErrTypeError)
	}

	user, err := ah.uf.GetUserByEmail(payload)
	if err != nil {
		return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "error": err})
	}

	jwtConfig := ah.uf.GetJWTConfig()
	accessTokenDetails, err := ah.ts.CreateToken(
		user.UUID.String(),
		jwtConfig.AccessTokenExpiredIn,
		jwtConfig.AccessTokenPrivateKey,
	)
	if err != nil {
		return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "error": err})
	}

	refreshTokenDetails, err := ah.ts.CreateToken(
		user.UUID.String(),
		jwtConfig.RefreshTokenExpiredIn,
		jwtConfig.RefreshTokenPrivateKey,
	)
	if err != nil {
		return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "error": err})
	}

	errAccess := ah.r.SetUserUUID(
		accessTokenDetails.TokenUUID,
		user.UUID.String(),
		*accessTokenDetails.ExpiresIn,
	)

	if errAccess != nil {
		return c.Status(errAccess.Status).JSON(fiber.Map{"status": "fail", "error": errAccess})
	}

	errRefresh := ah.r.SetUserUUID(
		refreshTokenDetails.TokenUUID,
		user.UUID.String(),
		*refreshTokenDetails.ExpiresIn,
	)

	if errRefresh != nil {
		return c.Status(errRefresh.Status).JSON(fiber.Map{"status": "fail", "error": errRefresh})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    *accessTokenDetails.Token,
		Path:     "/",
		Domain:   jwtConfig.Domain,
		MaxAge:   jwtConfig.AccessTokenMaxAge * 60,
		Secure:   jwtConfig.Secure,
		HTTPOnly: jwtConfig.HttpOnly,
		SameSite: "strict",
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    *refreshTokenDetails.Token,
		Path:     "/",
		Domain:   jwtConfig.Domain,
		MaxAge:   jwtConfig.RefreshTokenMaxAge * 60,
		Secure:   jwtConfig.Secure,
		HTTPOnly: jwtConfig.HttpOnly,
		SameSite: "strict",
	})

	return c.Status(fiber.StatusOK).
		JSON(fiber.Map{"status": "success", "access_token": accessTokenDetails.Token})
}

func (ah *AuthService) RefreshAccessToken(c *fiber.Ctx) error {
	refresh_token := c.Cookies("refresh_token")

	if refresh_token == "" {
		clientErr := restInterfaceErr.NewBadRequestError(errConst.ErrMsgPleaseLoginAgain)
		return c.Status(clientErr.Status).JSON(fiber.Map{"status": "fail", "error": clientErr})
	}

	jwtConfig := ah.uf.GetJWTConfig()
	tokenClaims, err := ah.ts.ValidateToken(refresh_token, jwtConfig.RefreshTokenPublicKey)
	if err != nil {
		return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "error": err})
	}

	userUuid, errGetTokenUUID := ah.r.GetUserUUID(tokenClaims.TokenUUID)
	if errGetTokenUUID != nil {
		return c.Status(errGetTokenUUID.Status).JSON(fiber.Map{"status": "fail", "error": errGetTokenUUID})
	}

	user, err := ah.uf.GetUserByUUID(userUuid)
	if err != nil {
		return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "error": err})
	}

	accessTokenDetails, err := ah.ts.CreateToken(
		user.UUID.String(),
		jwtConfig.AccessTokenExpiredIn,
		jwtConfig.AccessTokenPrivateKey,
	)
	if err != nil {
		return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "error": err})
	}

	errAccess := ah.r.SetUserUUID(
		accessTokenDetails.TokenUUID,
		user.UUID.String(),
		*accessTokenDetails.ExpiresIn,
	)

	if errAccess != nil {
		return c.Status(errAccess.Status).JSON(fiber.Map{"status": "fail", "error": errAccess})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    *accessTokenDetails.Token,
		Path:     "/",
		Domain:   jwtConfig.Domain,
		MaxAge:   jwtConfig.AccessTokenMaxAge * 60,
		Secure:   jwtConfig.Secure,
		HTTPOnly: jwtConfig.HttpOnly,
		SameSite: "strict",
	})

	return c.Status(fiber.StatusOK).
		JSON(fiber.Map{"status": "success", "access_token": accessTokenDetails.Token})
}

func (ah *AuthService) Logout(c *fiber.Ctx) error {
	refresh_token := c.Cookies("refresh_token")

	if refresh_token == "" {
		clientErr := restInterfaceErr.NewBadRequestError(errConst.ErrMsgPleaseLoginAgain)
		return c.Status(clientErr.Status).JSON(fiber.Map{"status": "fail", "error": clientErr})
	}

	jwtConfig := ah.uf.GetJWTConfig()
	tokenClaims, err := ah.ts.ValidateToken(refresh_token, jwtConfig.RefreshTokenPublicKey)
	if err != nil {
		return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "error": err})
	}

	accessTokenUUID, ok := c.Locals("accessTokenUUID").(string) // type assertion
	if !ok {
		internalErr := restInterfaceErr.NewBadRequestError(errMsgAccessTokenUUID)
		log.Error().Err(internalErr).Msg(errConst.ErrTypeError)
		clientErr := restInterfaceErr.NewBadRequestError(errConst.ErrMsgPleaseLoginAgain)
		return c.Status(clientErr.Status).JSON(fiber.Map{"status": "fail", "error": clientErr})
	}

	ah.r.DelUserUUID(tokenClaims.TokenUUID, accessTokenUUID)
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
