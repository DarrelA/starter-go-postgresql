package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/DarrelA/starter-go-postgresql/internal/domain/apperror"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/service"
)

const randomTokenSize = 32

// CSRFService creates and validates session-bound signed double-submit tokens.
type CSRFService struct{ secret []byte }

// NewCSRFService constructs an HMAC-backed CSRF service.
func NewCSRFService(secret string) (service.CSRFService, error) {
	if len(secret) < 32 {
		return nil, errors.New("CSRF secret must contain at least 32 characters")
	}
	return &CSRFService{secret: []byte(secret)}, nil
}

func (s *CSRFService) Create(sessionID string) (string, error) {
	if sessionID == "" {
		return "", errors.New("CSRF session ID is required")
	}
	random := make([]byte, randomTokenSize)
	if _, err := rand.Read(random); err != nil {
		return "", fmt.Errorf("create CSRF random value: %w", err)
	}
	randomValue := base64.RawURLEncoding.EncodeToString(random)
	return hex.EncodeToString(s.signature(sessionID, randomValue)) + "." + randomValue, nil
}

func (s *CSRFService) Validate(sessionID, token string) error {
	if sessionID == "" {
		return apperror.ErrInvalidCSRFToken
	}
	parts := strings.Split(token, ".")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return apperror.ErrInvalidCSRFToken
	}
	provided, err := hex.DecodeString(parts[0])
	if err != nil {
		return apperror.ErrInvalidCSRFToken
	}
	expected := s.signature(sessionID, parts[1])
	if len(provided) != len(expected) || subtle.ConstantTimeCompare(provided, expected) != 1 {
		return apperror.ErrInvalidCSRFToken
	}
	return nil
}

func (s *CSRFService) signature(sessionID, randomValue string) []byte {
	message := fmt.Sprintf("%d!%s!%d!%s", len(sessionID), sessionID, len(randomValue), randomValue)
	mac := hmac.New(sha256.New, s.secret)
	_, _ = mac.Write([]byte(message))
	return mac.Sum(nil)
}
