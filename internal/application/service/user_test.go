package service

import (
	"context"
	"errors"
	"testing"

	"github.com/DarrelA/starter-go-postgresql/internal/application/dto"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/apperror"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	"github.com/google/uuid"
)

type userRepositoryStub struct {
	saved       *entity.User
	byEmail     *entity.User
	byEmailErr  error
	byUUID      *entity.User
	byUUIDErr   error
	receivedCtx context.Context
}

func (r *userRepositoryStub) Save(ctx context.Context, user *entity.User) error {
	r.receivedCtx, r.saved = ctx, user
	id := uuid.MustParse("018f0000-0000-7000-8000-000000000001")
	user.UUID = &id
	return nil
}

func (r *userRepositoryStub) GetByEmail(ctx context.Context, _ string) (*entity.User, error) {
	r.receivedCtx = ctx
	return r.byEmail, r.byEmailErr
}

func (r *userRepositoryStub) GetByUUID(ctx context.Context, _ uuid.UUID) (*entity.User, error) {
	r.receivedCtx = ctx
	return r.byUUID, r.byUUIDErr
}

type passwordServiceStub struct {
	hash      string
	hashErr   error
	verifyErr error
}

type userServiceTestContextKey struct{}

func (p passwordServiceStub) Hash(string) (string, error) { return p.hash, p.hashErr }
func (p passwordServiceStub) Verify(string, string) error { return p.verifyErr }

func TestCreateUserHashesPasswordAndForwardsContext(t *testing.T) {
	repository := &userRepositoryStub{}
	service := NewUserService(repository, passwordServiceStub{hash: "hashed"})
	ctx := context.WithValue(context.Background(), userServiceTestContextKey{}, "request")

	result, err := service.CreateUser(ctx, dto.RegisterInput{
		FirstName: "Test", LastName: "User", Email: "test@example.com", Password: "secret",
	})
	if err != nil {
		t.Fatalf("create user: %v", err)
	}
	if repository.receivedCtx != ctx {
		t.Fatal("repository did not receive the request context")
	}
	if repository.saved.Password != "hashed" {
		t.Fatalf("expected hashed password, got %q", repository.saved.Password)
	}
	if result.UUID == nil || result.Email != "test@example.com" {
		t.Fatalf("unexpected response: %#v", result)
	}
}

func TestAuthenticateDoesNotRevealUnknownEmail(t *testing.T) {
	repository := &userRepositoryStub{byEmailErr: apperror.ErrUserNotFound}
	service := NewUserService(repository, passwordServiceStub{})

	_, err := service.Authenticate(context.Background(), dto.LoginInput{Email: "unknown@example.com"})
	if !errors.Is(err, apperror.ErrInvalidCredentials) {
		t.Fatalf("expected invalid credentials, got %v", err)
	}
}

func TestAuthenticateRejectsWrongPassword(t *testing.T) {
	repository := &userRepositoryStub{byEmail: &entity.User{Password: "hashed"}}
	service := NewUserService(repository, passwordServiceStub{verifyErr: apperror.ErrInvalidCredentials})

	_, err := service.Authenticate(context.Background(), dto.LoginInput{Password: "wrong"})
	if !errors.Is(err, apperror.ErrInvalidCredentials) {
		t.Fatalf("expected invalid credentials, got %v", err)
	}
}
