// coverage:ignore file
// Test file
package middleware

import (
	"time"

	dto "github.com/DarrelA/starter-go-postgresql/internal/application/dto"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	errConst "github.com/DarrelA/starter-go-postgresql/internal/error"
	restErr "github.com/DarrelA/starter-go-postgresql/internal/error"
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

type mockRedisUserRepository struct{ mid mockUUIDs }

func (m *mockRedisUserRepository) SetUserUUID(tokenUUID string, userUUID string, expiresIn int64) *restErr.RestErr {
	return nil
}

func (m *mockRedisUserRepository) GetUserUUID(tokenUUID string) (string, *restErr.RestErr) {
	return m.mid.mockUserUUID.String(), nil
}

func (m *mockRedisUserRepository) DelUserUUID(tokenUUID string, accessTokenUUID string) (int64, *restErr.RestErr) {
	return 1, nil
}

type mockTokenService struct{ mid mockUUIDs }

func (m *mockTokenService) CreateToken(userUUID string, ttl time.Duration, privateKey string) (
	*entity.Token, *restErr.RestErr) {
	return nil, nil
}

func (m *mockTokenService) ValidateToken(token string, publicKey string) (
	*entity.Token, *restErr.RestErr) {
	// Simulate invalid token
	if token == "" || token == "mockInvalidBearerToken" {
		err := restErr.NewUnauthorizedError(errConst.ErrMsgPleaseLoginAgain)
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

type mockUserService struct{}

func (m *mockUserService) GetJWTConfig() *entity.JWTConfig { return &entity.JWTConfig{} }

func (m *mockUserService) CreateUser(payload dto.RegisterInput) (*dto.UserResponse, *restErr.RestErr) {
	return nil, nil
}

func (m *mockUserService) GetUserByEmail(u dto.LoginInput) (*dto.UserResponse, *restErr.RestErr) {
	return nil, nil
}

func (m *mockUserService) GetUserByUUID(userUuid string) (*entity.User, *restErr.RestErr) {
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
