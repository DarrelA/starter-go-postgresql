package repository

import (
	"context"

	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
)

// TokenRepository stores revocable token sessions until their expiry.
type TokenRepository interface {
	Create(ctx context.Context, session entity.TokenSession) error
	Rotate(ctx context.Context, currentRefreshTokenUUID string, replacement entity.TokenSession) error
	GetUserUUID(ctx context.Context, familyUUID, accessTokenUUID string) (string, error)
	Revoke(ctx context.Context, familyUUID, refreshTokenUUID string) error
}
