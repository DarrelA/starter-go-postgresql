package auth

import (
	"strings"
	"testing"
	"time"

	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	"github.com/google/uuid"
)

func TestTokenServiceCreatesAndValidatesBothTokenTypes(t *testing.T) {
	accessPrivate, accessPublic := generateRSAKeyPair(t)
	refreshPrivate, refreshPublic := generateRSAKeyPair(t)
	service, err := NewJWTService(&entity.JWTConfig{
		AccessTokenPrivateKey:  accessPrivate,
		AccessTokenPublicKey:   accessPublic,
		RefreshTokenPrivateKey: refreshPrivate,
		RefreshTokenPublicKey:  refreshPublic,
	})
	if err != nil {
		t.Fatalf("create token service: %v", err)
	}

	userUUID := uuid.Must(uuid.NewV7()).String()
	familyUUID := uuid.Must(uuid.NewV7()).String()
	accessToken, err := service.CreateAccessToken(userUUID, familyUUID, testTTL)
	if err != nil {
		t.Fatalf("create access token: %v", err)
	}
	accessClaims, err := service.ValidateAccessToken(*accessToken.Token)
	if err != nil {
		t.Fatalf("validate access token: %v", err)
	}
	if accessClaims.FamilyUUID != familyUUID {
		t.Fatalf("expected family UUID %q, got %q", familyUUID, accessClaims.FamilyUUID)
	}

	refreshToken, err := service.CreateRefreshToken(userUUID, familyUUID, testTTL)
	if err != nil {
		t.Fatalf("create refresh token: %v", err)
	}
	if _, err := service.ValidateRefreshToken(*refreshToken.Token); err != nil {
		t.Fatalf("validate refresh token: %v", err)
	}
	if _, err := service.ValidateRefreshToken(*accessToken.Token); err == nil {
		t.Fatal("expected access token to be rejected by refresh key")
	}
}

func TestTokenServiceRejectsExpiredToken(t *testing.T) {
	accessPrivate, accessPublic := generateRSAKeyPair(t)
	refreshPrivate, refreshPublic := generateRSAKeyPair(t)
	service, err := NewJWTService(&entity.JWTConfig{
		AccessTokenPrivateKey:  accessPrivate,
		AccessTokenPublicKey:   accessPublic,
		RefreshTokenPrivateKey: refreshPrivate,
		RefreshTokenPublicKey:  refreshPublic,
	})
	if err != nil {
		t.Fatalf("create token service: %v", err)
	}
	token, err := service.CreateRefreshToken(
		uuid.Must(uuid.NewV7()).String(),
		uuid.Must(uuid.NewV7()).String(),
		-time.Second,
	)
	if err != nil {
		t.Fatalf("create expired token: %v", err)
	}
	if _, err := service.ValidateRefreshToken(*token.Token); err == nil {
		t.Fatal("expected expired refresh token to be rejected")
	}
}

func TestNewJWTServiceRejectsInvalidKeyConfiguration(t *testing.T) {
	accessPrivate, accessPublic := generateRSAKeyPair(t)
	refreshPrivate, refreshPublic := generateRSAKeyPair(t)
	_, differentPublic := generateRSAKeyPair(t)

	tests := []struct {
		name    string
		config  *entity.JWTConfig
		message string
	}{
		{name: "missing config", message: "configuration is required"},
		{
			name: "invalid private key",
			config: &entity.JWTConfig{
				AccessTokenPrivateKey: "invalid", AccessTokenPublicKey: accessPublic,
				RefreshTokenPrivateKey: refreshPrivate, RefreshTokenPublicKey: refreshPublic,
			},
			message: "decode access private key",
		},
		{
			name: "mismatched pair",
			config: &entity.JWTConfig{
				AccessTokenPrivateKey: accessPrivate, AccessTokenPublicKey: differentPublic,
				RefreshTokenPrivateKey: refreshPrivate, RefreshTokenPublicKey: refreshPublic,
			},
			message: "access public and private keys do not match",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewJWTService(test.config)
			if err == nil || !strings.Contains(err.Error(), test.message) {
				t.Fatalf("expected %q, got %v", test.message, err)
			}
		})
	}
}
