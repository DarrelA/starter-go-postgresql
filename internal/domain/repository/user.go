package repository

import (
	"context"

	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	"github.com/google/uuid"
)

// UserRepository defines the user persistence required by application services.
type UserRepository interface {
	Save(ctx context.Context, user *entity.User) error
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetByUUID(ctx context.Context, userUUID uuid.UUID) (*entity.User, error)
}
