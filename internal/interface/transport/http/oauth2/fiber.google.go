package http

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/DarrelA/starter-go-postgresql/internal/application/usecase"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	errConst "github.com/DarrelA/starter-go-postgresql/internal/domain/error"
	restInterfaceErr "github.com/DarrelA/starter-go-postgresql/internal/interface/transport/http/error"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const googleUserInfoEndpoint = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

type GoogleOAuth2 struct {
	GoogleLoginConfig oauth2.Config
}

func NewGoogleOAuth2(OAuth2Config *entity.OAuth2Config) usecase.OAuth2UseCase {
	googleLoginConfig := oauth2.Config{
		RedirectURL:  OAuth2Config.GoogleRedirectURL,
		ClientID:     OAuth2Config.GoogleClientID,
		ClientSecret: OAuth2Config.GoogleClientSecret,
		Scopes:       OAuth2Config.Scopes,
		Endpoint:     google.Endpoint,
	}

	return &GoogleOAuth2{GoogleLoginConfig: googleLoginConfig}
}

func (oa GoogleOAuth2) Login(c *fiber.Ctx) error {
	url := oa.GoogleLoginConfig.AuthCodeURL("randomstate")
	return c.Redirect(url, fiber.StatusSeeOther)
}

func (oa GoogleOAuth2) Callback(c *fiber.Ctx) error {
	state := c.Query("state")
	if state != "randomstate" {
		err := restInterfaceErr.NewBadRequestError(errConst.ErrMsgPleaseLoginAgain)
		return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "error": err})
	}

	code := c.Query("code")
	token, err := oa.GoogleLoginConfig.Exchange(context.Background(), code)
	if err != nil {
		err := restInterfaceErr.NewBadRequestError(errConst.ErrMsgPleaseLoginAgain)
		return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "error": err})
	}

	resp, err := http.Get(googleUserInfoEndpoint + token.AccessToken)
	if err != nil {
		log.Error().Err(err).Msg(errConst.ErrMsgGoogleOAuth2Error)
		err := restInterfaceErr.NewInternalServerError(errConst.ErrMsgSomethingWentWrong)
		return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "error": err})
	}

	userInfo, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg(errConst.ErrMsgGoogleOAuth2Error)
		err := restInterfaceErr.NewInternalServerError(errConst.ErrMsgSomethingWentWrong)
		return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "error": err})
	}

	var user map[string]interface{}
	if err := json.Unmarshal(userInfo, &user); err != nil {
		log.Error().Err(err).Msg("")
		err := restInterfaceErr.NewInternalServerError(errConst.ErrMsgSomethingWentWrong)
		return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "error": err})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "user": user})
}
