package service

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	"github.com/DarrelA/starter-go-postgresql/internal/utils/err_rest"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type Token struct {
	t entity.Token
}

func NewUserFactory(t entity.Token) *Token {
	return &Token{t}
}

func (t *Token) CreateToken(user_uuid string, ttl time.Duration, privateKey string) (*Token, *err_rest.RestErr) {
	now := time.Now().UTC()
	t.t.ExpiresIn = new(int64)
	t.t.Token = new(string)
	*t.t.ExpiresIn = now.Add(ttl).Unix()

	id, err := uuid.NewV7()
	if err != nil {
		log.Error().Err(err).Msg("uuid_error")
		return nil, err_rest.NewUnprocessableEntityError("something went wrong")
	}

	t.t.TokenUUID = id.String()
	t.t.UserUUID = user_uuid

	decodePrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		log.Error().Err(err).Msg("DecodeString_privateKey_error")
		return nil, err_rest.NewUnprocessableEntityError("something went wrong")
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodePrivateKey)

	if err != nil {
		log.Error().Err(err).Msg("ParseRSA_privateKey_error")
		return nil, err_rest.NewUnprocessableEntityError("something went wrong")
	}

	atClaims := make(jwt.MapClaims)
	atClaims["sub"] = user_uuid
	atClaims["token_uuid"] = t.t.TokenUUID
	atClaims["exp"] = t.t.ExpiresIn
	atClaims["iat"] = now.Unix() // Issued at
	atClaims["nbf"] = now.Unix() // Not before

	*t.t.Token, err = jwt.NewWithClaims(jwt.SigningMethodRS256, atClaims).SignedString(key)
	if err != nil {
		log.Error().Err(err).Msg("sign_key_error")
		return nil, err_rest.NewUnprocessableEntityError("something went wrong")
	}

	return t, nil
}

func (t *Token) ValidateToken(token string, publicKey string) (*Token, *err_rest.RestErr) {
	decodePublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		log.Error().Err(err).Msg("DecodeString_publicKey_error")
		return nil, err_rest.NewInternalServerError("something went wrong")
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(decodePublicKey)
	if err != nil {
		log.Error().Err(err).Msg("ParseRSA_publicKey_error")
		return nil, err_rest.NewInternalServerError("something went wrong")
	}

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
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

	return &Token{
		t: entity.Token{
			TokenUUID: fmt.Sprint(claims["token_uuid"]),
			UserUUID:  fmt.Sprint(claims["subs"]),
		},
	}, nil
}
