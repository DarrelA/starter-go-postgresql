package http

import (
	"crypto/subtle"
	"errors"
	"time"

	dto "github.com/DarrelA/starter-go-postgresql/internal/application/dto"
	appSvc "github.com/DarrelA/starter-go-postgresql/internal/application/service"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/apperror"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/repository"
	domainSvc "github.com/DarrelA/starter-go-postgresql/internal/domain/service"
	restErr "github.com/DarrelA/starter-go-postgresql/internal/error"
	"github.com/DarrelA/starter-go-postgresql/internal/interface/transport/http/authcookie"
	"github.com/DarrelA/starter-go-postgresql/internal/interface/transport/http/response"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

const (
	errMsgRegisterPayload  = "register_payload is not of type users.RegisterInput"
	errMsgLoginPayload     = "login_payload is not of type users.RegisterInput"
	accessTokenCookieName  = authcookie.AccessTokenName
	refreshTokenCookieName = authcookie.RefreshTokenName
	csrfTokenCookieName    = authcookie.CSRFTokenName
	csrfTokenHeaderName    = authcookie.CSRFHeaderName
	authCookiePath         = authcookie.Path
	authCookieSameSite     = authcookie.SameSite
)

// AuthHandler translates authentication HTTP requests into application operations.
type AuthHandler struct {
	r       repository.TokenRepository
	us      appSvc.UserService
	ss      appSvc.SessionService
	ts      domainSvc.TokenService
	csrf    domainSvc.CSRFService
	cookies *authcookie.Manager
	jwt     *entity.JWTConfig
}

// NewAuthHandler constructs an authentication HTTP handler.
func NewAuthHandler(
	r repository.TokenRepository,
	us appSvc.UserService,
	ss appSvc.SessionService,
	ts domainSvc.TokenService,
	csrfService domainSvc.CSRFService,
	jwtConfig *entity.JWTConfig,
) *AuthHandler {
	return &AuthHandler{
		r: r, us: us, ss: ss, ts: ts, csrf: csrfService,
		cookies: authcookie.NewManager(jwtConfig), jwt: jwtConfig,
	}
}

// Register creates a password user from a preprocessed registration payload.
func (auc *AuthHandler) Register(c *fiber.Ctx) error {
	payload, ok := c.Locals("register_payload").(dto.RegisterInput)
	if !ok {
		err := restErr.NewBadRequestError(errMsgRegisterPayload)
		log.Error().Err(err).Msg(restErr.ErrTypeError)
		return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "error": err})
	}

	result, err := auc.us.CreateUser(c.UserContext(), payload)
	if err != nil {
		return response.Error(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "user": result})
}

// Login authenticates credentials and issues access and refresh token sessions.
func (auc *AuthHandler) Login(c *fiber.Ctx) error {
	payload, ok := c.Locals("login_payload").(dto.LoginInput)
	if !ok {
		err := restErr.NewBadRequestError(errMsgLoginPayload)
		log.Error().Err(err).Msg(restErr.ErrTypeError)
		return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "error": err})
	}

	user, err := auc.us.Authenticate(c.UserContext(), payload)
	if err != nil {
		return response.Error(c, err)
	}

	session, err := auc.ss.Issue(c.UserContext(), user.UUID.String())
	if err != nil {
		return response.Error(c, err)
	}
	auc.cookies.Set(c, session)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}

// RefreshAccessToken validates a refresh session and issues a new access token.
func (auc *AuthHandler) RefreshAccessToken(c *fiber.Ctx) error {
	refreshToken := c.Cookies(refreshTokenCookieName)

	if refreshToken == "" {
		clientErr := restErr.NewBadRequestError(restErr.ErrMsgPleaseLoginAgain)
		return c.Status(clientErr.Status).JSON(fiber.Map{"status": "fail", "error": clientErr})
	}

	jwtConfig := auc.jwt
	tokenClaims, err := auc.ts.ValidateRefreshToken(refreshToken)
	if err != nil {
		return response.Error(c, err)
	}

	if err := auc.validateCSRFToken(c, tokenClaims.TokenUUID); err != nil {
		return response.Error(c, err)
	}

	accessTokenDetails, err := auc.ts.CreateAccessToken(
		tokenClaims.UserUUID,
		tokenClaims.FamilyUUID,
		jwtConfig.AccessTokenExpiredIn,
	)
	if err != nil {
		return response.Error(c, err)
	}

	refreshTokenDetails, err := auc.ts.CreateRefreshToken(
		tokenClaims.UserUUID,
		tokenClaims.FamilyUUID,
		jwtConfig.RefreshTokenExpiredIn,
	)
	if err != nil {
		return response.Error(c, err)
	}
	csrfToken, err := auc.csrf.Create(refreshTokenDetails.TokenUUID)
	if err != nil {
		return response.Error(c, err)
	}

	if err := auc.r.Rotate(c.UserContext(), tokenClaims.TokenUUID, tokenSession(accessTokenDetails, refreshTokenDetails)); err != nil {
		if errors.Is(err, apperror.ErrRefreshTokenReused) || errors.Is(err, apperror.ErrSessionRevoked) {
			auc.clearSessionCookies(c)
		}
		return response.Error(c, err)
	}

	auc.cookies.Set(c, &dto.AuthSession{
		AccessToken: *accessTokenDetails.Token, RefreshToken: *refreshTokenDetails.Token, CSRFToken: csrfToken,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}

// Logout revokes the current access and refresh token sessions.
func (auc *AuthHandler) Logout(c *fiber.Ctx) error {
	refreshToken := c.Cookies(refreshTokenCookieName)

	if refreshToken == "" {
		clientErr := restErr.NewBadRequestError(restErr.ErrMsgPleaseLoginAgain)
		return c.Status(clientErr.Status).JSON(fiber.Map{"status": "fail", "error": clientErr})
	}

	tokenClaims, err := auc.ts.ValidateRefreshToken(refreshToken)
	if err != nil {
		auc.clearSessionCookies(c)
		return response.Error(c, err)
	}

	if err := auc.validateCSRFToken(c, tokenClaims.TokenUUID); err != nil {
		return response.Error(c, err)
	}

	if revokeErr := auc.r.Revoke(c.UserContext(), tokenClaims.FamilyUUID, tokenClaims.TokenUUID); revokeErr != nil &&
		!errors.Is(revokeErr, apperror.ErrSessionNotFound) {
		return response.Error(c, revokeErr)
	}
	auc.clearSessionCookies(c)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}

func tokenSession(accessToken, refreshToken *entity.Token) entity.TokenSession {
	return entity.TokenSession{
		UserUUID:         refreshToken.UserUUID,
		FamilyUUID:       refreshToken.FamilyUUID,
		AccessTokenUUID:  accessToken.TokenUUID,
		RefreshTokenUUID: refreshToken.TokenUUID,
		AccessExpiresAt:  time.Unix(*accessToken.ExpiresIn, 0),
		RefreshExpiresAt: time.Unix(*refreshToken.ExpiresIn, 0),
	}
}

func (auc *AuthHandler) validateCSRFToken(c *fiber.Ctx, sessionID string) error {
	cookieToken := c.Cookies(csrfTokenCookieName)
	headerToken := c.Get(csrfTokenHeaderName)
	if cookieToken == "" || headerToken == "" || len(cookieToken) != len(headerToken) ||
		subtle.ConstantTimeCompare([]byte(cookieToken), []byte(headerToken)) != 1 {
		return apperror.ErrInvalidCSRFToken
	}
	return auc.csrf.Validate(sessionID, headerToken)
}

func (auc *AuthHandler) clearSessionCookies(c *fiber.Ctx) {
	auc.cookies.Clear(c)
}
