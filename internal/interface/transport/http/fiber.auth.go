// coverage:ignore file
// Testing with integration test
package http

import (
	"time"

	dto "github.com/DarrelA/starter-go-postgresql/internal/application/dto"
	appSvc "github.com/DarrelA/starter-go-postgresql/internal/application/service"
	"github.com/DarrelA/starter-go-postgresql/internal/application/usecase"
	r "github.com/DarrelA/starter-go-postgresql/internal/domain/repository/redis"
	domainSvc "github.com/DarrelA/starter-go-postgresql/internal/domain/service"
	restErr "github.com/DarrelA/starter-go-postgresql/internal/error"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

const (
	errMsgRegisterPayload = "register_payload is not of type users.RegisterInput"
	errMsgLoginPayload    = "login_payload is not of type users.RegisterInput"
	errMsgAccessTokenUUID = "accessTokenUUID is not a string or not set"
)

type AuthUseCase struct {
	r  r.RedisUserRepository
	us appSvc.UserService
	ts domainSvc.TokenService
}

func NewAuthUseCase(
	r r.RedisUserRepository,
	us appSvc.UserService,
	ts domainSvc.TokenService,
) usecase.AuthUseCase {
	return &AuthUseCase{r, us, ts}
}

func (auc *AuthUseCase) Register(c *fiber.Ctx) error {
	payload, ok := c.Locals("register_payload").(dto.RegisterInput)
	if !ok {
		err := restErr.NewBadRequestError(errMsgRegisterPayload)
		log.Error().Err(err).Msg(restErr.ErrTypeError)
	}

	result, err := auc.us.CreateUser(payload)
	if err != nil {
		return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "error": err})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "user": result})
}

func (auc *AuthUseCase) Login(c *fiber.Ctx) error {
	payload, ok := c.Locals("login_payload").(dto.LoginInput)
	if !ok {
		err := restErr.NewBadRequestError(errMsgLoginPayload)
		log.Error().Err(err).Msg(restErr.ErrTypeError)
	}

	user, err := auc.us.GetUserByEmail(payload)
	if err != nil {
		return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "error": err})
	}

	jwtConfig := auc.us.GetJWTConfig()
	accessTokenDetails, err := auc.ts.CreateToken(
		user.UUID.String(),
		jwtConfig.AccessTokenExpiredIn,
		jwtConfig.AccessTokenPrivateKey,
	)
	if err != nil {
		return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "error": err})
	}

	refreshTokenDetails, err := auc.ts.CreateToken(
		user.UUID.String(),
		jwtConfig.RefreshTokenExpiredIn,
		jwtConfig.RefreshTokenPrivateKey,
	)
	if err != nil {
		return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "error": err})
	}

	errAccess := auc.r.SetUserUUID(
		accessTokenDetails.TokenUUID,
		user.UUID.String(),
		*accessTokenDetails.ExpiresIn,
	)

	if errAccess != nil {
		return c.Status(errAccess.Status).JSON(fiber.Map{"status": "fail", "error": errAccess})
	}

	errRefresh := auc.r.SetUserUUID(
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

func (auc *AuthUseCase) RefreshAccessToken(c *fiber.Ctx) error {
	refresh_token := c.Cookies("refresh_token")

	if refresh_token == "" {
		clientErr := restErr.NewBadRequestError(restErr.ErrMsgPleaseLoginAgain)
		return c.Status(clientErr.Status).JSON(fiber.Map{"status": "fail", "error": clientErr})
	}

	jwtConfig := auc.us.GetJWTConfig()
	tokenClaims, err := auc.ts.ValidateToken(refresh_token, jwtConfig.RefreshTokenPublicKey)
	if err != nil {
		return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "error": err})
	}

	userUuid, errGetTokenUUID := auc.r.GetUserUUID(tokenClaims.TokenUUID)
	if errGetTokenUUID != nil {
		return c.Status(errGetTokenUUID.Status).JSON(fiber.Map{"status": "fail", "error": errGetTokenUUID})
	}

	user, err := auc.us.GetUserByUUID(userUuid)
	if err != nil {
		return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "error": err})
	}

	accessTokenDetails, err := auc.ts.CreateToken(
		user.UUID.String(),
		jwtConfig.AccessTokenExpiredIn,
		jwtConfig.AccessTokenPrivateKey,
	)
	if err != nil {
		return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "error": err})
	}

	errAccess := auc.r.SetUserUUID(
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

func (auc *AuthUseCase) Logout(c *fiber.Ctx) error {
	refresh_token := c.Cookies("refresh_token")

	if refresh_token == "" {
		clientErr := restErr.NewBadRequestError(restErr.ErrMsgPleaseLoginAgain)
		return c.Status(clientErr.Status).JSON(fiber.Map{"status": "fail", "error": clientErr})
	}

	jwtConfig := auc.us.GetJWTConfig()
	tokenClaims, err := auc.ts.ValidateToken(refresh_token, jwtConfig.RefreshTokenPublicKey)
	if err != nil {
		return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "error": err})
	}

	accessTokenUUID, ok := c.Locals("accessTokenUUID").(string) // type assertion
	if !ok {
		internalErr := restErr.NewBadRequestError(errMsgAccessTokenUUID)
		log.Error().Err(internalErr).Msg(restErr.ErrTypeError)
		clientErr := restErr.NewBadRequestError(restErr.ErrMsgPleaseLoginAgain)
		return c.Status(clientErr.Status).JSON(fiber.Map{"status": "fail", "error": clientErr})
	}

	auc.r.DelUserUUID(tokenClaims.TokenUUID, accessTokenUUID)
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
