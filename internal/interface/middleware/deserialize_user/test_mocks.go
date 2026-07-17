// coverage:ignore file
// Test file
package middleware

import (
	"context"
	"time"

	"github.com/DarrelA/starter-go-postgresql/internal/domain/apperror"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	"github.com/google/uuid"
)

type mockUUIDs struct {
	mockUserUUID  *uuid.UUID
	mockTokenUUID *uuid.UUID
}

func (m *mockUUIDs) initializeMockUUIDEntities() {
	mockUserUUID, _ := uuid.NewV7()
	mockTokenUUID, _ := uuid.NewV7()
	m.mockUserUUID = &mockUserUUID
	m.mockTokenUUID = &mockTokenUUID
}

type mockRedisUserRepository struct {
	mid   mockUUIDs
	err   error
	calls int
}

func (m *mockRedisUserRepository) Create(context.Context, entity.TokenSession) error {
	return nil
}

func (m *mockRedisUserRepository) Rotate(context.Context, string, entity.TokenSession) error {
	return nil
}

func (m *mockRedisUserRepository) GetUserUUID(context.Context, string, string) (string, error) {
	m.calls++
	return m.mid.mockUserUUID.String(), m.err
}

func (m *mockRedisUserRepository) Revoke(context.Context, string, string) error {
	return nil
}

type mockTokenService struct{ mid mockUUIDs }

func (m *mockTokenService) CreateAccessToken(string, string, time.Duration) (*entity.Token, error) {
	return nil, nil
}

func (m *mockTokenService) CreateRefreshToken(string, string, time.Duration) (*entity.Token, error) {
	return nil, nil
}

func (m *mockTokenService) ValidateAccessToken(token string) (*entity.Token, error) {
	// Simulate invalid token
	if token == "" || token == "mockInvalidBearerToken" {
		return nil, apperror.ErrInvalidToken
	}

	// Simulate valid token
	expiresIn := mockExpiresIn
	mockToken := &entity.Token{
		Token:      &token,
		TokenUUID:  m.mid.mockTokenUUID.String(),
		FamilyUUID: m.mid.mockTokenUUID.String(),
		UserUUID:   m.mid.mockUserUUID.String(),
		ExpiresIn:  &expiresIn,
	}

	return mockToken, nil
}

func (m *mockTokenService) ValidateRefreshToken(token string) (*entity.Token, error) {
	return m.ValidateAccessToken(token)
}
