package http

import (
	"context"
	stdhttp "net/http"
	"net/http/httptest"
	"testing"

	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	authmw "github.com/DarrelA/starter-go-postgresql/internal/interface/middleware/deserialize_user"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type userReaderStub struct {
	userUUID uuid.UUID
	calls    int
}

func (s *userReaderStub) GetUserByUUID(_ context.Context, userUUID string) (*entity.User, error) {
	s.calls++
	if userUUID != s.userUUID.String() {
		return nil, nil
	}
	return &entity.User{UUID: &s.userUUID, Email: "user@example.com"}, nil
}

func TestGetUserRecord(t *testing.T) {
	userUUID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")

	t.Run("missing request identity", func(t *testing.T) {
		users := &userReaderStub{userUUID: userUUID}
		app := fiber.New()
		app.Get("/me", NewUserHandler(users).GetUserRecord)

		resp, err := app.Test(httptest.NewRequest(stdhttp.MethodGet, "/me", nil))
		if err != nil {
			t.Fatalf("execute request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != stdhttp.StatusInternalServerError {
			t.Fatalf("expected status %d, got %d", stdhttp.StatusInternalServerError, resp.StatusCode)
		}
		if users.calls != 0 {
			t.Fatalf("missing identity caused %d user lookup(s)", users.calls)
		}
	})

	t.Run("loads user for verified identity", func(t *testing.T) {
		users := &userReaderStub{userUUID: userUUID}
		app := fiber.New()
		app.Use(func(c *fiber.Ctx) error {
			authmw.SetRequestIdentity(c, authmw.RequestIdentity{UserUUID: userUUID})
			return c.Next()
		})
		app.Get("/me", NewUserHandler(users).GetUserRecord)

		resp, err := app.Test(httptest.NewRequest(stdhttp.MethodGet, "/me", nil))
		if err != nil {
			t.Fatalf("execute request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != stdhttp.StatusOK {
			t.Fatalf("expected status %d, got %d", stdhttp.StatusOK, resp.StatusCode)
		}
		if users.calls != 1 {
			t.Fatalf("expected one user lookup, got %d", users.calls)
		}
	})
}
