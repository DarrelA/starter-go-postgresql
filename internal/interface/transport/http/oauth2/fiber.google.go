package http

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/DarrelA/starter-go-postgresql/internal/application/dto"
	appservice "github.com/DarrelA/starter-go-postgresql/internal/application/service"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	restErr "github.com/DarrelA/starter-go-postgresql/internal/error"
	"github.com/DarrelA/starter-go-postgresql/internal/interface/transport/http/authcookie"
	"github.com/DarrelA/starter-go-postgresql/internal/interface/transport/http/response"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	googleUserInfoEndpoint = "https://www.googleapis.com/oauth2/v2/userinfo"
	oauthStateCookie       = "oauth_state"
	oauthStateTTL          = 10 * time.Minute
	maxUserInfoSize        = 1 << 20
)

// GoogleOAuth2 handles the Google OAuth authorization-code flow.
type GoogleOAuth2 struct {
	GoogleLoginConfig oauth2.Config
	httpClient        *http.Client
	userInfoEndpoint  string
	oauth             appservice.OAuthService
	sessions          appservice.SessionService
	cookies           *authcookie.Manager
}

// NewGoogleOAuth2 constructs Google OAuth handlers from application configuration.
func NewGoogleOAuth2(
	OAuth2Config *entity.OAuth2Config,
	jwtConfig *entity.JWTConfig,
	oauthService appservice.OAuthService,
	sessionService appservice.SessionService,
) *GoogleOAuth2 {
	googleLoginConfig := oauth2.Config{
		RedirectURL:  OAuth2Config.GoogleRedirectURL,
		ClientID:     OAuth2Config.GoogleClientID,
		ClientSecret: OAuth2Config.GoogleClientSecret,
		Scopes:       OAuth2Config.Scopes,
		Endpoint:     google.Endpoint,
	}

	return &GoogleOAuth2{
		GoogleLoginConfig: googleLoginConfig,
		httpClient:        &http.Client{Timeout: 10 * time.Second},
		userInfoEndpoint:  googleUserInfoEndpoint,
		oauth:             oauthService,
		sessions:          sessionService,
		cookies:           authcookie.NewManager(jwtConfig),
	}
}

// Login creates OAuth state and redirects the client to Google.
func (oa GoogleOAuth2) Login(c *fiber.Ctx) error {
	state, err := newOAuthState()
	if err != nil {
		log.Error().Err(err).Msg(restErr.ErrMsgGoogleOAuth2Error)
		serverErr := restErr.NewInternalServerError(restErr.ErrMsgSomethingWentWrong)
		return c.Status(serverErr.Status).JSON(fiber.Map{"status": "fail", "error": serverErr})
	}

	c.Cookie(&fiber.Cookie{
		Name:     oauthStateCookie,
		Value:    state,
		Path:     "/auth/google_callback",
		MaxAge:   int(oauthStateTTL.Seconds()),
		Secure:   c.Protocol() == "https",
		HTTPOnly: true,
		SameSite: "lax",
	})

	url := oa.GoogleLoginConfig.AuthCodeURL(state)
	return c.Redirect(url, fiber.StatusSeeOther)
}

// Callback validates OAuth state, resolves a local user, and establishes a session.
func (oa GoogleOAuth2) Callback(c *fiber.Ctx) error {
	state := c.Query("state")
	expectedState := c.Cookies(oauthStateCookie)
	expireOAuthStateCookie(c)

	if state == "" || expectedState == "" || len(state) != len(expectedState) ||
		subtle.ConstantTimeCompare([]byte(state), []byte(expectedState)) != 1 {
		err := restErr.NewBadRequestError(restErr.ErrMsgPleaseLoginAgain)
		return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "error": err})
	}

	code := c.Query("code")
	if code == "" {
		err := restErr.NewBadRequestError(restErr.ErrMsgPleaseLoginAgain)
		return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "error": err})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 10*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, oauth2.HTTPClient, oa.httpClient)

	token, err := oa.GoogleLoginConfig.Exchange(ctx, code)
	if err != nil {
		err := restErr.NewBadRequestError(restErr.ErrMsgPleaseLoginAgain)
		return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "error": err})
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, oa.userInfoEndpoint, nil)
	if err != nil {
		log.Error().Err(err).Msg(restErr.ErrMsgGoogleOAuth2Error)
		serverErr := restErr.NewInternalServerError(restErr.ErrMsgSomethingWentWrong)
		return c.Status(serverErr.Status).JSON(fiber.Map{"status": "fail", "error": serverErr})
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	resp, err := oa.httpClient.Do(req)
	if err != nil {
		log.Error().Err(err).Msg(restErr.ErrMsgGoogleOAuth2Error)
		serverErr := restErr.NewBadGatewayError(restErr.ErrMsgSomethingWentWrong)
		return c.Status(serverErr.Status).JSON(fiber.Map{"status": "fail", "error": serverErr})
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		log.Error().Int("status", resp.StatusCode).Msg(restErr.ErrMsgGoogleOAuth2Error)
		serverErr := restErr.NewBadGatewayError(restErr.ErrMsgSomethingWentWrong)
		return c.Status(serverErr.Status).JSON(fiber.Map{"status": "fail", "error": serverErr})
	}

	var googleProfile dto.GoogleProfile
	decoder := json.NewDecoder(io.LimitReader(resp.Body, maxUserInfoSize))
	if err := decoder.Decode(&googleProfile); err != nil {
		log.Error().Err(err).Msg(restErr.ErrMsgGoogleOAuth2Error)
		err := restErr.NewInternalServerError(restErr.ErrMsgSomethingWentWrong)
		return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "error": err})
	}

	user, err := oa.oauth.Authenticate(ctx, dto.OAuthProfile{
		Provider: entity.OAuthProviderGoogle, Subject: googleProfile.Subject,
		Email: googleProfile.Email, EmailVerified: googleProfile.EmailVerified,
		FirstName: googleProfile.FirstName, LastName: googleProfile.LastName,
	})
	if err != nil {
		return response.Error(c, err)
	}
	session, err := oa.sessions.Issue(ctx, user.UUID.String())
	if err != nil {
		return response.Error(c, err)
	}
	oa.cookies.Set(c, session)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "user": user})
}

func newOAuthState() (string, error) {
	state := make([]byte, 32)
	if _, err := rand.Read(state); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(state), nil
}

func expireOAuthStateCookie(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     oauthStateCookie,
		Value:    "",
		Path:     "/auth/google_callback",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		Secure:   c.Protocol() == "https",
		HTTPOnly: true,
		SameSite: "lax",
	})
}
