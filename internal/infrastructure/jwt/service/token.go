package jwt

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/service"
	"github.com/DarrelA/starter-go-postgresql/internal/utils/err_rest"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

/*
- The `TokenService` should be a stateless service that performs operations related to tokens. It does not need to manage the lifecycle of the `Token` entity itself but rather uses it.
- You might inject dependencies into `TokenService` if it interacts with other services or repositories.
*/
type TokenService struct{}

func NewTokenService() service.TokenService {
	return &TokenService{}
}

func (ts *TokenService) CreateToken(userUUID string, ttl time.Duration, privateKey string) (*entity.Token, *err_rest.RestErr) {
	now := time.Now().UTC()
	t := &entity.Token{
		ExpiresIn: new(int64),
		Token:     new(string),
	}

	id, err := uuid.NewV7()
	if err != nil {
		log.Error().Err(err).Msg("uuid_error")
		return nil, err_rest.NewUnprocessableEntityError("something went wrong")
	}

	t.TokenUUID = id.String()
	t.UserUUID = userUUID
	*t.ExpiresIn = now.Add(ttl).Unix()

	decodedPrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		log.Error().Err(err).Msg("DecodeString_privateKey_error")
		return nil, err_rest.NewUnprocessableEntityError("something went wrong")
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)

	if err != nil {
		log.Error().Err(err).Msg("ParseRSA_privateKey_error")
		return nil, err_rest.NewUnprocessableEntityError("something went wrong")
	}

	atClaims := jwt.MapClaims{
		"sub":        userUUID,
		"token_uuid": t.TokenUUID,
		"exp":        *t.ExpiresIn,
		"iat":        now.Unix(), // Issued at
		"nbf":        now.Unix(), // Not before
	}

	*t.Token, err = jwt.NewWithClaims(jwt.SigningMethodRS256, atClaims).SignedString(key)
	if err != nil {
		log.Error().Err(err).Msg("sign_key_error")
		return nil, err_rest.NewUnprocessableEntityError("something went wrong")
	}

	return t, nil
}

func (ts *TokenService) ValidateToken(tokenStr string, publicKey string) (*entity.Token, *err_rest.RestErr) {
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		log.Error().Err(err).Msg("DecodeString_publicKey_error")
		return nil, err_rest.NewInternalServerError("something went wrong")
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)
	if err != nil {
		log.Error().Err(err).Msg("ParseRSA_publicKey_error")
		return nil, err_rest.NewInternalServerError("something went wrong")
	}

	parsedToken, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
		}
		return key, nil
	})

	if err != nil {
		log.Error().Err(err).Msg("validate_token_error")
		return nil, err_rest.NewForbiddenError("please login again")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		err := err_rest.NewForbiddenError("validation: invalid_token")
		log.Error().Err(err).Msg("")
		return nil, err_rest.NewForbiddenError("please login again")
	}

	return &entity.Token{
		TokenUUID: fmt.Sprint(claims["token_uuid"]),
		UserUUID:  fmt.Sprint(claims["sub"]),
	}, nil
}
