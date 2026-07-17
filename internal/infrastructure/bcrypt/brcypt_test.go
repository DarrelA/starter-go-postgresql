package bcrypt

import (
	"errors"
	"testing"

	errConst "github.com/DarrelA/starter-go-postgresql/internal/error"
	"golang.org/x/crypto/bcrypt"
)

const (
	password          = "testPassword"
	incorrectPassword = "wrongPassword"
)

func TestHashPassword(t *testing.T) {
	t.Run("successful hash", func(t *testing.T) {
		hashedPassword, err := HashPassword(password)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if hashedPassword == "" {
			t.Errorf("Expected hashed password, got empty string")
		}

		// Verify that the hashed password can be used to compare with the original password
		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		if err != nil {
			t.Errorf("Expected the hashed password to match the original password, got error: %v", err)
		}
	})

	t.Run("bcrypt error", func(t *testing.T) {
		_, err := HashPassword(generateLongPassword(80))
		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
	})
}

func TestVerifyPassword(t *testing.T) {
	hashedPassword, err := HashPassword(password)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err = VerifyPassword(hashedPassword, password)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test with incorrect password
	err = VerifyPassword(hashedPassword, incorrectPassword)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	expectedErr := errors.New(errConst.ErrMsgInvalidCredentials)
	if err.Error() != expectedErr.Error() {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}
}

func generateLongPassword(length int) string {
	password := make([]byte, length)
	for i := range password {
		password[i] = 'a'
	}
	return string(password)
}
