package bcrypt

import (
	"errors"

	errConst "github.com/DarrelA/starter-go-postgresql/internal/error"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

const errMsgBCryptError = "bcrypt_error"

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Error().Err(err).Msg(errMsgBCryptError)
		return "", errors.New(errMsgBCryptError)
	}

	return string(hashedPassword), nil
}

func VerifyPassword(hashedPassword string, inputPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword)); err != nil {
		return errors.New(errConst.ErrMsgInvalidCredentials)
	}

	return nil
}
