package auth

import (
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/DarrelA/starter-go-postgresql/internal/domain/apperror"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/service"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/google/uuid"
)

const (
	errMsgSignKeyError     = "sign key error"
	errMsgUnexpectedMethod = "unexpected signing method: %s"
)

// JWTService creates and validates JWTs without retaining session state.
type JWTService struct {
	accessPrivateKey  *rsa.PrivateKey
	accessPublicKey   *rsa.PublicKey
	refreshPrivateKey *rsa.PrivateKey
	refreshPublicKey  *rsa.PublicKey
}

// NewJWTService parses and verifies all signing keys before serving requests.
func NewJWTService(config *entity.JWTConfig) (service.TokenService, error) {
	if config == nil {
		return nil, errors.New("JWT configuration is required")
	}
	accessPrivate, accessPublic, err := parseKeyPair(
		"access", config.AccessTokenPrivateKey, config.AccessTokenPublicKey,
	)
	if err != nil {
		return nil, err
	}
	refreshPrivate, refreshPublic, err := parseKeyPair(
		"refresh", config.RefreshTokenPrivateKey, config.RefreshTokenPublicKey,
	)
	if err != nil {
		return nil, err
	}
	return &JWTService{
		accessPrivateKey:  accessPrivate,
		accessPublicKey:   accessPublic,
		refreshPrivateKey: refreshPrivate,
		refreshPublicKey:  refreshPublic,
	}, nil
}

func (ts *JWTService) CreateAccessToken(userUUID, familyUUID string, ttl time.Duration) (*entity.Token, error) {
	return ts.createToken(userUUID, familyUUID, ttl, ts.accessPrivateKey)
}

func (ts *JWTService) CreateRefreshToken(userUUID, familyUUID string, ttl time.Duration) (*entity.Token, error) {
	return ts.createToken(userUUID, familyUUID, ttl, ts.refreshPrivateKey)
}

func (ts *JWTService) ValidateAccessToken(token string) (*entity.Token, error) {
	return ts.validateToken(token, ts.accessPublicKey)
}

func (ts *JWTService) ValidateRefreshToken(token string) (*entity.Token, error) {
	return ts.validateToken(token, ts.refreshPublicKey)
}

func (ts *JWTService) createToken(userUUID, familyUUID string, ttl time.Duration, privateKey *rsa.PrivateKey) (
	*entity.Token, error) {
	now := time.Now().UTC()
	t := &entity.Token{
		ExpiresIn: new(int64),
		Token:     new(string),
	}

	id, err := uuid.NewV7()
	if err != nil { // coverage:ignore
		return nil, fmt.Errorf("create token UUID: %w", err)
	}

	t.TokenUUID = id.String()
	t.UserUUID = userUUID
	t.FamilyUUID = familyUUID
	*t.ExpiresIn = now.Add(ttl).Unix()

	atClaims := jwt.MapClaims{
		"sub":         userUUID,
		"token_uuid":  t.TokenUUID,
		"family_uuid": familyUUID,
		"exp":         *t.ExpiresIn,
		"iat":         now.Unix(), // Issued at
		"nbf":         now.Unix(), // Not before
	}

	*t.Token, err = jwt.NewWithClaims(jwt.SigningMethodRS256, atClaims).SignedString(privateKey)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errMsgSignKeyError, err)
	}

	return t, nil
}

func (ts *JWTService) validateToken(tokenStr string, publicKey *rsa.PublicKey) (
	*entity.Token, error) {
	parsedToken, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf(errMsgUnexpectedMethod, t.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, errors.Join(apperror.ErrInvalidToken, err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, apperror.ErrInvalidToken
	}

	return &entity.Token{
		TokenUUID:  fmt.Sprint(claims["token_uuid"]),
		UserUUID:   fmt.Sprint(claims["sub"]),
		FamilyUUID: fmt.Sprint(claims["family_uuid"]),
	}, nil
}

func parseKeyPair(name, encodedPrivate, encodedPublic string) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privatePEM, err := base64.StdEncoding.DecodeString(encodedPrivate)
	if err != nil {
		return nil, nil, fmt.Errorf("decode %s private key: %w", name, err)
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePEM)
	if err != nil {
		return nil, nil, fmt.Errorf("parse %s private key: %w", name, err)
	}

	publicPEM, err := base64.StdEncoding.DecodeString(encodedPublic)
	if err != nil {
		return nil, nil, fmt.Errorf("decode %s public key: %w", name, err)
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicPEM)
	if err != nil {
		return nil, nil, fmt.Errorf("parse %s public key: %w", name, err)
	}
	if privateKey.PublicKey.E != publicKey.E || privateKey.PublicKey.N.Cmp(publicKey.N) != 0 {
		return nil, nil, fmt.Errorf("%s public and private keys do not match", name)
	}
	return privateKey, publicKey, nil
}
