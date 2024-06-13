package services

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/DarrelA/starter-go-postgresql/internal/utils/errors"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/google/uuid"
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
		errors.NewInternalServerError("failed to generate UUID: " + err.Error())
	}
	td.TokenUUID = id.String()

	td.UserUUID = user_uuid

	decodePrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return nil, errors.NewInternalServerError("could not decode token private key: " + err.Error())
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodePrivateKey)

	if err != nil {
		return nil, errors.NewInternalServerError("create: parse token private key: " + err.Error())
	}

	atClaims := make(jwt.MapClaims)
	atClaims["sub"] = user_uuid
	atClaims["token_uuid"] = td.TokenUUID
	atClaims["exp"] = td.ExpiresIn
	atClaims["iat"] = now.Unix() // Issued at
	atClaims["nbf"] = now.Unix() // Not before

	*td.Token, err = jwt.NewWithClaims(jwt.SigningMethodRS256, atClaims).SignedString(key)
	if err != nil {
		return nil, errors.NewInternalServerError("create: sign token: " + err.Error())
	}

	return td, nil
}

func ValidateToken(token string, publicKey string) (*TokenDetails, error) {
	decodePublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return nil, fmt.Errorf("could not decode: %w", err)
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(decodePublicKey)
	if err != nil {
		return nil, fmt.Errorf("validate: parse key: %w", err)
	}

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
		}
		return key, nil
	})

	if err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("validate: invalid token")
	}

	return &TokenDetails{
		TokenUUID: fmt.Sprint(claims["token_uuid"]),
		UserUUID:  fmt.Sprint(claims["subs"]),
	}, nil
}
