package service

import (
	"time"

	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
)

// TokenService creates and validates signed authentication tokens.
type TokenService interface {
	CreateAccessToken(userUUID, familyUUID string, ttl time.Duration) (*entity.Token, error)
	CreateRefreshToken(userUUID, familyUUID string, ttl time.Duration) (*entity.Token, error)
	ValidateAccessToken(token string) (*entity.Token, error)
	ValidateRefreshToken(token string) (*entity.Token, error)
}

type CSRFService interface {
	Create(sessionID string) (string, error)
	Validate(sessionID, token string) error
}
