package services

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/DarrelA/starter-go-postgresql/internal/utils/errors"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type TokenDetails struct {
	Token     *string
	TokenUUID string
	UserUUID  string
	ExpiresIn *int64
}

func CreateToken(user_uuid string, ttl time.Duration, privateKey string) (*TokenDetails, *errors.RestErr) {
	now := time.Now().UTC()
	td := &TokenDetails{
		ExpiresIn: new(int64),
		Token:     new(string),
	}

	*td.ExpiresIn = now.Add(ttl).Unix()

	id, err := uuid.NewV7()
	if err != nil {
		log.Error().Err(err).Msg("uuid_error")
		return nil, errors.NewUnprocessableEntityError("something went wrong")
	}
	td.TokenUUID = id.String()

	td.UserUUID = user_uuid

	decodePrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		log.Error().Err(err).Msg("DecodeString_privateKey_error")
		return nil, errors.NewUnprocessableEntityError("something went wrong")
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodePrivateKey)

	if err != nil {
		log.Error().Err(err).Msg("ParseRSA_privateKey_error")
		return nil, errors.NewUnprocessableEntityError("something went wrong")
	}

	atClaims := make(jwt.MapClaims)
	atClaims["sub"] = user_uuid
	atClaims["token_uuid"] = td.TokenUUID
	atClaims["exp"] = td.ExpiresIn
	atClaims["iat"] = now.Unix() // Issued at
	atClaims["nbf"] = now.Unix() // Not before

	*td.Token, err = jwt.NewWithClaims(jwt.SigningMethodRS256, atClaims).SignedString(key)
	if err != nil {
		log.Error().Err(err).Msg("sign_key_error")
		return nil, errors.NewUnprocessableEntityError("something went wrong")
	}

	return td, nil
}

func ValidateToken(token string, publicKey string) (*TokenDetails, *errors.RestErr) {
	decodePublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		log.Error().Err(err).Msg("DecodeString_publicKey_error")
		return nil, errors.NewInternalServerError("something went wrong")
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(decodePublicKey)
	if err != nil {
		log.Error().Err(err).Msg("ParseRSA_publicKey_error")
		return nil, errors.NewInternalServerError("something went wrong")
	}

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
		}
		return key, nil
	})

	if err != nil {
		log.Error().Err(err).Msg("validate_token_error")
		return nil, errors.NewForbiddenError("please login again")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		err := errors.NewForbiddenError("validation: invalid_token")
		log.Error().Err(err).Msg("")
		return nil, errors.NewForbiddenError("please login again")
	}

	return &TokenDetails{
		TokenUUID: fmt.Sprint(claims["token_uuid"]),
		UserUUID:  fmt.Sprint(claims["subs"]),
	}, nil
}
