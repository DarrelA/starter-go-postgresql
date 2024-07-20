package bcrypt

import (
	"errors"

	errConst "github.com/DarrelA/starter-go-postgresql/internal/domain/error"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

const ErrMsgBCryptError = "bcrypt_error"

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Error().Err(err).Msg(ErrMsgBCryptError)
		return "", errors.New(ErrMsgBCryptError)
	}

	return string(hashedPassword), nil
}

func VerifyPassword(hashedPassword string, inputPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword)); err != nil {
		return errors.New(errConst.ErrMsgInvalidCredentials)
	}

	return nil
}
