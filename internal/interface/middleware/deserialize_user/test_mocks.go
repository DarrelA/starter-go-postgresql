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

func (m *mockRedisUserRepository) SetUserUUID(tokenUUID string, userUUID string, expiresIn int64) *restDomainErr.RestErr {
	return nil
}

func (m *mockRedisUserRepository) GetUserUUID(tokenUUID string) (string, *restDomainErr.RestErr) {
	return m.mid.mockUserUUID.String(), nil
}

func (m *mockRedisUserRepository) DelUserUUID(tokenUUID string, accessTokenUUID string) (int64, *restDomainErr.RestErr) {
	return 1, nil
}

type mockTokenService struct{ mid mockUUIDs }

func (m *mockTokenService) CreateToken(userUUID string, ttl time.Duration, privateKey string) (
	*entity.Token, *restDomainErr.RestErr) {
	return nil, nil
}

func (m *mockTokenService) ValidateToken(token string, publicKey string) (
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

type mockUserService struct{}

func (m *mockUserService) GetJWTConfig() *entity.JWTConfig { return &entity.JWTConfig{} }

func (m *mockUserService) CreateUser(payload dto.RegisterInput) (*dto.UserResponse, *restDomainErr.RestErr) {
	return nil, nil
}

func (m *mockUserService) GetUserByEmail(u dto.LoginInput) (*dto.UserResponse, *restDomainErr.RestErr) {
	return nil, nil
}

func (m *mockUserService) GetUserByUUID(userUuid string) (*entity.User, *restDomainErr.RestErr) {
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
