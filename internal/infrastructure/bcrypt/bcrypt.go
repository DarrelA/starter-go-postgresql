package bcrypt

import (
	"fmt"

	"github.com/DarrelA/starter-go-postgresql/internal/domain/apperror"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

const errMsgBCryptError = "bcrypt_error"

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (Service) Hash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Error().Err(err).Msg(errMsgBCryptError)
		return "", fmt.Errorf("%s: %w", errMsgBCryptError, err)
	}

	return string(hashedPassword), nil
}

func (Service) Verify(hashedPassword string, inputPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword)); err != nil {
		return apperror.ErrInvalidCredentials
	}

	return nil
}

func HashPassword(password string) (string, error) {
	return Service{}.Hash(password)
}

func VerifyPassword(hashedPassword string, inputPassword string) error {
	return Service{}.Verify(hashedPassword, inputPassword)
}
