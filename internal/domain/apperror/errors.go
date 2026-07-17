package apperror

import "errors"

var (
	ErrEmailConflict         = errors.New("email already exists")
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrInvalidToken          = errors.New("invalid token")
	ErrSessionNotFound       = errors.New("session not found")
	ErrSessionRevoked        = errors.New("session revoked")
	ErrRefreshTokenReused    = errors.New("refresh token reused")
	ErrInvalidCSRFToken      = errors.New("invalid CSRF token")
	ErrUserNotFound          = errors.New("user not found")
	ErrOAuthProfileInvalid   = errors.New("invalid OAuth profile")
	ErrOAuthIdentityConflict = errors.New("OAuth identity conflict")
)
