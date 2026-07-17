package service

import (
	"context"
	"testing"
	"time"

	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
)

type sessionTokenServiceStub struct{}

func (sessionTokenServiceStub) CreateAccessToken(userUUID, familyUUID string, ttl time.Duration) (*entity.Token, error) {
	return testSessionToken("access", "access-id", userUUID, familyUUID, ttl), nil
}
func (sessionTokenServiceStub) CreateRefreshToken(userUUID, familyUUID string, ttl time.Duration) (*entity.Token, error) {
	return testSessionToken("refresh", "refresh-id", userUUID, familyUUID, ttl), nil
}
func (sessionTokenServiceStub) ValidateAccessToken(string) (*entity.Token, error)  { return nil, nil }
func (sessionTokenServiceStub) ValidateRefreshToken(string) (*entity.Token, error) { return nil, nil }

type sessionCSRFServiceStub struct{}

func (sessionCSRFServiceStub) Create(sessionID string) (string, error) {
	return "csrf-" + sessionID, nil
}
func (sessionCSRFServiceStub) Validate(string, string) error { return nil }

type sessionTokenRepositoryStub struct{ created entity.TokenSession }

func (r *sessionTokenRepositoryStub) Create(_ context.Context, session entity.TokenSession) error {
	r.created = session
	return nil
}
func (*sessionTokenRepositoryStub) Rotate(context.Context, string, entity.TokenSession) error {
	return nil
}
func (*sessionTokenRepositoryStub) GetUserUUID(context.Context, string, string) (string, error) {
	return "", nil
}
func (*sessionTokenRepositoryStub) Revoke(context.Context, string, string) error { return nil }

func TestSessionServiceIssuesCompleteSession(t *testing.T) {
	repository := &sessionTokenRepositoryStub{}
	service := NewSessionService(sessionTokenServiceStub{}, sessionCSRFServiceStub{}, repository, &entity.JWTConfig{
		AccessTokenExpiredIn: time.Minute, RefreshTokenExpiredIn: 2 * time.Minute,
	})

	session, err := service.Issue(context.Background(), "user-id")
	if err != nil {
		t.Fatalf("issue session: %v", err)
	}
	if session.AccessToken != "access" || session.RefreshToken != "refresh" || session.CSRFToken != "csrf-refresh-id" {
		t.Fatalf("unexpected session: %#v", session)
	}
	if repository.created.UserUUID != "user-id" || repository.created.FamilyUUID == "" ||
		repository.created.AccessTokenUUID != "access-id" || repository.created.RefreshTokenUUID != "refresh-id" {
		t.Fatalf("unexpected persisted session: %#v", repository.created)
	}
}

func testSessionToken(value, tokenUUID, userUUID, familyUUID string, ttl time.Duration) *entity.Token {
	expires := time.Now().Add(ttl).Unix()
	return &entity.Token{
		Token: &value, TokenUUID: tokenUUID, UserUUID: userUUID, FamilyUUID: familyUUID, ExpiresIn: &expires,
	}
}
