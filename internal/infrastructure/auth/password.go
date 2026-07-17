package auth

import (
	"fmt"

	"github.com/DarrelA/starter-go-postgresql/internal/domain/apperror"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

const errMsgBCryptError = "bcrypt_error"

// PasswordService hashes and verifies password credentials with bcrypt.
type PasswordService struct{}

// NewPasswordService constructs a bcrypt-backed password service.
func NewPasswordService() *PasswordService {
	return &PasswordService{}
}

func (PasswordService) Hash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Error().Err(err).Msg(errMsgBCryptError)
		return "", fmt.Errorf("%s: %w", errMsgBCryptError, err)
	}

	return string(hashedPassword), nil
}

func (PasswordService) Verify(hashedPassword string, inputPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword)); err != nil {
		return apperror.ErrInvalidCredentials
	}

	return nil
}
