package service

import (
	"context"
	"errors"
	"fmt"

	dto "github.com/DarrelA/starter-go-postgresql/internal/application/dto"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/apperror"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/repository"
	"github.com/google/uuid"
)

// UserService defines user operations exposed to transport adapters.
type UserService interface {
	CreateUser(ctx context.Context, payload dto.RegisterInput) (*dto.UserResponse, error)
	Authenticate(ctx context.Context, input dto.LoginInput) (*dto.UserResponse, error)
	GetUserByUUID(ctx context.Context, userUUID string) (*entity.User, error)
}

// PasswordService defines the password operations required by userService.
type PasswordService interface {
	Hash(password string) (string, error)
	Verify(hashedPassword, inputPassword string) error
}

type userService struct {
	users     repository.UserRepository
	passwords PasswordService
}

func NewUserService(users repository.UserRepository, passwords PasswordService) UserService {
	return &userService{users: users, passwords: passwords}
}

func (s *userService) CreateUser(ctx context.Context, payload dto.RegisterInput) (*dto.UserResponse, error) {
	user := &entity.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  payload.Password,
	}

	hashedPassword, err := s.passwords.Hash(user.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}
	user.Password = hashedPassword

	if err := s.users.Save(ctx, user); err != nil {
		return nil, err
	}

	return toUserResponse(user), nil
}

func (s *userService) Authenticate(ctx context.Context, input dto.LoginInput) (*dto.UserResponse, error) {
	user, err := s.users.GetByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, apperror.ErrUserNotFound) {
			return nil, apperror.ErrInvalidCredentials
		}
		return nil, err
	}

	if err := s.passwords.Verify(user.Password, input.Password); err != nil {
		return nil, apperror.ErrInvalidCredentials
	}

	return toUserResponse(user), nil
}

func (s *userService) GetUserByUUID(ctx context.Context, userUUID string) (*entity.User, error) {
	id, err := uuid.Parse(userUUID)
	if err != nil {
		return nil, fmt.Errorf("parse user UUID: %w", err)
	}

	return s.users.GetByUUID(ctx, id)
}

func toUserResponse(user *entity.User) *dto.UserResponse {
	return &dto.UserResponse{
		UUID:      user.UUID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}
}
