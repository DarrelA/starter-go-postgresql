package service

import (
	"context"
	"strings"

	"github.com/DarrelA/starter-go-postgresql/internal/application/dto"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/apperror"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/repository"
)

// OAuthService resolves a verified provider profile to a local user.
type OAuthService interface {
	Authenticate(ctx context.Context, profile dto.OAuthProfile) (*dto.UserResponse, error)
}

type oauthService struct {
	identities repository.ProviderIdentityRepository
}

func NewOAuthService(identities repository.ProviderIdentityRepository) OAuthService {
	return &oauthService{identities: identities}
}

func (s *oauthService) Authenticate(ctx context.Context, profile dto.OAuthProfile) (*dto.UserResponse, error) {
	profile.Provider = strings.ToLower(strings.TrimSpace(profile.Provider))
	profile.Subject = strings.TrimSpace(profile.Subject)
	profile.Email = strings.ToLower(strings.TrimSpace(profile.Email))
	profile.FirstName = strings.TrimSpace(profile.FirstName)
	profile.LastName = strings.TrimSpace(profile.LastName)
	if profile.Provider == "" || profile.Subject == "" || profile.Email == "" || !profile.EmailVerified {
		return nil, apperror.ErrOAuthProfileInvalid
	}

	user, err := s.identities.FindLinkOrCreate(ctx, entity.ProviderIdentity{
		Provider: profile.Provider,
		Subject:  profile.Subject,
	}, &entity.User{
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
		Email:     profile.Email,
	})
	if err != nil {
		return nil, err
	}
	return toUserResponse(user), nil
}
