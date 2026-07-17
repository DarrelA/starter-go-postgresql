package service

import (
	"context"
	"fmt"
	"time"

	"github.com/DarrelA/starter-go-postgresql/internal/application/dto"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/repository"
	domainservice "github.com/DarrelA/starter-go-postgresql/internal/domain/service"
	"github.com/google/uuid"
)

// SessionService issues the same application session for every login method.
type SessionService interface {
	Issue(ctx context.Context, userUUID string) (*dto.AuthSession, error)
}

type sessionService struct {
	tokens domainservice.TokenService
	csrf   domainservice.CSRFService
	repo   repository.TokenRepository
	config *entity.JWTConfig
}

func NewSessionService(
	tokens domainservice.TokenService,
	csrf domainservice.CSRFService,
	repo repository.TokenRepository,
	config *entity.JWTConfig,
) SessionService {
	return &sessionService{tokens: tokens, csrf: csrf, repo: repo, config: config}
}

func (s *sessionService) Issue(ctx context.Context, userUUID string) (*dto.AuthSession, error) {
	familyUUID, err := uuid.NewV7()
	if err != nil { // coverage:ignore
		return nil, fmt.Errorf("create token family UUID: %w", err)
	}
	access, err := s.tokens.CreateAccessToken(userUUID, familyUUID.String(), s.config.AccessTokenExpiredIn)
	if err != nil {
		return nil, err
	}
	refresh, err := s.tokens.CreateRefreshToken(userUUID, familyUUID.String(), s.config.RefreshTokenExpiredIn)
	if err != nil {
		return nil, err
	}
	csrfToken, err := s.csrf.Create(refresh.TokenUUID)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Create(ctx, tokenSession(access, refresh)); err != nil {
		return nil, err
	}
	return &dto.AuthSession{
		AccessToken: *access.Token, RefreshToken: *refresh.Token, CSRFToken: csrfToken,
	}, nil
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
