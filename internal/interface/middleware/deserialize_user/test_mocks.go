// coverage:ignore file
// Test file
package middleware

import (
	"time"

	dto "github.com/DarrelA/starter-go-postgresql/internal/application/dto"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	errConst "github.com/DarrelA/starter-go-postgresql/internal/domain/error"
	restDomainErr "github.com/DarrelA/starter-go-postgresql/internal/domain/error/transport/http"
	restInterfaceErr "github.com/DarrelA/starter-go-postgresql/internal/interface/transport/http/error"
	"github.com/google/uuid"
)

type MockUUIDs struct {
	mockUserUUID  *uuid.UUID
	mockTokenUUID *uuid.UUID
}

func (m *MockUUIDs) initializeMockUUIDEntities() {
	mockUserUUID, _ := uuid.NewV7()
	mockTokenUUID, _ := uuid.NewV7()
	m.mockUserUUID = &mockUserUUID
	m.mockTokenUUID = &mockTokenUUID
}

type MockRedisUserRepository struct{ mid MockUUIDs }

func (m *MockRedisUserRepository) SetUserUUID(tokenUUID string, userUUID string, expiresIn int64) *restDomainErr.RestErr {
	return nil
}

func (m *MockRedisUserRepository) GetUserUUID(tokenUUID string) (string, *restDomainErr.RestErr) {
	return m.mid.mockUserUUID.String(), nil
}

func (m *MockRedisUserRepository) DelUserUUID(tokenUUID string, accessTokenUUID string) (int64, *restDomainErr.RestErr) {
	return 1, nil
}

type MockTokenService struct{ mid MockUUIDs }

func (m *MockTokenService) CreateToken(userUUID string, ttl time.Duration, privateKey string) (
	*entity.Token, *restDomainErr.RestErr) {
	return nil, nil
}

func (m *MockTokenService) ValidateToken(token string, publicKey string) (
	*entity.Token, *restDomainErr.RestErr) {
	// Simulate invalid token
	if token == "" || token == "mockInvalidBearerToken" {
		err := restInterfaceErr.NewUnauthorizedError(errConst.ErrMsgPleaseLoginAgain)
		return nil, err
	}

	// Simulate valid token
	expiresIn := mockExpiresIn
	mockToken := &entity.Token{
		Token:     &token,
		TokenUUID: m.mid.mockTokenUUID.String(),
		UserUUID:  m.mid.mockUserUUID.String(),
		ExpiresIn: &expiresIn,
	}

	return mockToken, nil
}

type UserFactory struct{}

func (m *UserFactory) GetJWTConfig() *entity.JWTConfig { return &entity.JWTConfig{} }

func (m *UserFactory) CreateUser(payload dto.RegisterInput) (*dto.UserResponse, *restDomainErr.RestErr) {
	return nil, nil
}

func (m *UserFactory) GetUserByEmail(u dto.LoginInput) (*dto.UserResponse, *restDomainErr.RestErr) {
	return nil, nil
}

func (m *UserFactory) GetUserByUUID(userUuid string) (*entity.User, *restDomainErr.RestErr) {
	uuidPointer, _ := uuid.Parse(userUuid)
	user := &entity.User{
		ID:        int64(1),
		UUID:      &uuidPointer,
		FirstName: "jie",
		LastName:  "wei",
		Email:     "jiewei@gmail.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return user, nil
}
