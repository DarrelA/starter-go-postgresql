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

type providerIdentityRepositoryStub struct {
	identity  entity.ProviderIdentity
	candidate *entity.User
	result    *entity.User
	err       error
}

func (r *providerIdentityRepositoryStub) FindLinkOrCreate(
	_ context.Context,
	identity entity.ProviderIdentity,
	candidate *entity.User,
) (*entity.User, error) {
	r.identity, r.candidate = identity, candidate
	return r.result, r.err
}

func TestOAuthServiceResolvesVerifiedProfile(t *testing.T) {
	id := uuid.MustParse("018f0000-0000-7000-8000-000000000010")
	repository := &providerIdentityRepositoryStub{result: &entity.User{
		UUID: &id, Email: "user@example.com", FirstName: "Test", LastName: "User",
	}}
	service := NewOAuthService(repository)

	result, err := service.Authenticate(context.Background(), dto.OAuthProfile{
		Provider: " GOOGLE ", Subject: " subject-1 ", Email: " USER@EXAMPLE.COM ",
		EmailVerified: true, FirstName: " Test ", LastName: " User ",
	})
	if err != nil {
		t.Fatalf("authenticate OAuth profile: %v", err)
	}
	if repository.identity.Provider != entity.OAuthProviderGoogle || repository.identity.Subject != "subject-1" {
		t.Fatalf("unexpected identity: %#v", repository.identity)
	}
	if repository.candidate.Email != "user@example.com" || result.UUID == nil {
		t.Fatalf("unexpected user resolution: candidate=%#v result=%#v", repository.candidate, result)
	}
}

func TestOAuthServiceRejectsInvalidProfiles(t *testing.T) {
	service := NewOAuthService(&providerIdentityRepositoryStub{})
	for _, profile := range []dto.OAuthProfile{
		{Subject: "subject", Email: "user@example.com", EmailVerified: true},
		{Provider: "google", Email: "user@example.com", EmailVerified: true},
		{Provider: "google", Subject: "subject", EmailVerified: true},
		{Provider: "google", Subject: "subject", Email: "user@example.com"},
	} {
		if _, err := service.Authenticate(context.Background(), profile); !errors.Is(err, apperror.ErrOAuthProfileInvalid) {
			t.Fatalf("expected invalid profile for %#v, got %v", profile, err)
		}
	}
}

func TestOAuthServicePropagatesIdentityConflict(t *testing.T) {
	service := NewOAuthService(&providerIdentityRepositoryStub{err: apperror.ErrOAuthIdentityConflict})
	_, err := service.Authenticate(context.Background(), dto.OAuthProfile{
		Provider: "google", Subject: "subject", Email: "user@example.com", EmailVerified: true,
	})
	if !errors.Is(err, apperror.ErrOAuthIdentityConflict) {
		t.Fatalf("expected identity conflict, got %v", err)
	}
}
