package http

import (
	"context"

	dto "github.com/DarrelA/starter-go-postgresql/internal/application/dto"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	restErr "github.com/DarrelA/starter-go-postgresql/internal/error"
	authmw "github.com/DarrelA/starter-go-postgresql/internal/interface/middleware/deserialize_user"
	"github.com/DarrelA/starter-go-postgresql/internal/interface/transport/http/response"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type userReader interface {
	GetUserByUUID(ctx context.Context, userUUID string) (*entity.User, error)
}

// UserHandler serves authenticated user resources.
type UserHandler struct{ users userReader }

// NewUserHandler constructs a user HTTP handler.
func NewUserHandler(users userReader) *UserHandler {
	return &UserHandler{users: users}
}

// GetUserRecord loads and returns the complete record for the verified identity.
func (uuc *UserHandler) GetUserRecord(c *fiber.Ctx) error {
	identity, ok := authmw.RequestIdentityFrom(c)
	if !ok {
		err := restErr.NewInternalServerError(restErr.ErrMsgSomethingWentWrong)
		log.Error().Err(err).Msg(restErr.ErrTypeError)
		return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "error": err})
	}
	user, err := uuc.users.GetUserByUUID(c.UserContext(), identity.UserUUID.String())
	if err != nil {
		return response.Error(c, err)
	}
	userRecord := &dto.UserRecord{
		UUID: user.UUID, FirstName: user.FirstName, LastName: user.LastName,
		Email: user.Email, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).
		JSON(fiber.Map{"status": "success", "data": fiber.Map{"userRecord": userRecord}})
}
