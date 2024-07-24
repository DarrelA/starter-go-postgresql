package jwt

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	errConst "github.com/DarrelA/starter-go-postgresql/internal/domain/error"
	restDomainErr "github.com/DarrelA/starter-go-postgresql/internal/domain/error/transport/http"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/service"
	restInterfaceErr "github.com/DarrelA/starter-go-postgresql/internal/interface/transport/http/error"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

const (
	errMsgDecodeStringError = "DecodeString error"
	errMsgParseRSAKeyError  = "ParseRSA error"
	errMsgSignKeyError      = "sign key error"
	errMsgUnexpectedMethod  = "unexpected signing method: %s"
)

/*
The `TokenService` should be a stateless service that performs operations related to tokens.
It does not need to manage the lifecycle of the `Token` entity itself but rather uses it.
*/
type TokenService struct{}

func NewTokenService() service.TokenService {
	return &TokenService{}
}

func (ts *TokenService) CreateToken(userUUID string, ttl time.Duration, privateKey string) (
	*entity.Token, *restDomainErr.RestErr) {
	now := time.Now().UTC()
	t := &entity.Token{
		ExpiresIn: new(int64),
		Token:     new(string),
	}

	id, err := uuid.NewV7()
	if err != nil { // coverage:ignore
		log.Error().Err(err).Msg(errConst.ErrUUIDError)
		return nil, restInterfaceErr.NewUnprocessableEntityError(errConst.ErrMsgSomethingWentWrong)
	}

	t.TokenUUID = id.String()
	t.UserUUID = userUUID
	*t.ExpiresIn = now.Add(ttl).Unix()

	decodedPrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		log.Error().Err(err).Msg(errMsgDecodeStringError)
		return nil, restInterfaceErr.NewUnprocessableEntityError(errConst.ErrMsgSomethingWentWrong)
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)

	if err != nil {
		log.Error().Err(err).Msg(errMsgParseRSAKeyError)
		return nil, restInterfaceErr.NewUnprocessableEntityError(errConst.ErrMsgSomethingWentWrong)
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
		log.Error().Err(err).Msg(errMsgSignKeyError)
		return nil, restInterfaceErr.NewUnprocessableEntityError(errConst.ErrMsgSomethingWentWrong)
	}

	return t, nil
}

func (ts *TokenService) ValidateToken(tokenStr string, publicKey string) (
	*entity.Token, *restDomainErr.RestErr) {
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		log.Error().Err(err).Msg(errMsgDecodeStringError)
		return nil, restInterfaceErr.NewInternalServerError(errConst.ErrMsgSomethingWentWrong)
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)
	if err != nil {
		log.Error().Err(err).Msg(errMsgParseRSAKeyError)
		return nil, restInterfaceErr.NewInternalServerError(errConst.ErrMsgSomethingWentWrong)
	}

	parsedToken, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf(errMsgUnexpectedMethod, t.Header["alg"])
		}
		return key, nil
	})

	if err != nil {
		log.Error().Err(err).Msg(errConst.ErrMsgInvalidToken)
		return nil, restInterfaceErr.NewForbiddenError(errConst.ErrMsgPleaseLoginAgain)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		err := restInterfaceErr.NewForbiddenError(errConst.ErrMsgInvalidToken)
		log.Error().Err(err).Msg("")
		return nil, restInterfaceErr.NewForbiddenError(errConst.ErrMsgPleaseLoginAgain)
	}

	return &entity.Token{
		TokenUUID: fmt.Sprint(claims["token_uuid"]),
		UserUUID:  fmt.Sprint(claims["sub"]),
	}, nil
}