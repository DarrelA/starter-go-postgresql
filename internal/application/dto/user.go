package dto

import (
	"time"

	"github.com/google/uuid"
)

type RegisterInput struct {
	FirstName string `json:"first_name" validate:"required,min=2,max=50,alpha"`
	LastName  string `json:"last_name" validate:"required,min=2,max=50,alpha"`
	Email     string `json:"email" validate:"required,min=5,max=64,email"`
	Password  string `json:"password" validate:"required,min=8,passwd"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,max=100,email"`
	Password string `json:"password" validate:"required,max=100"`
}

type UserResponse struct {
	UUID      *uuid.UUID `json:"uuid"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Email     string     `json:"email"`
}

type UserRecord struct {
	UUID      *uuid.UUID `json:"uuid"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Email     string     `json:"email"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// OAuthProfile is the provider-neutral identity accepted by the application.
type OAuthProfile struct {
	Provider      string
	Subject       string
	Email         string
	EmailVerified bool
	FirstName     string
	LastName      string
}

// GoogleProfile is the typed response returned by Google's user-info endpoint.
type GoogleProfile struct {
	Subject       string `json:"id"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"verified_email"`
	FirstName     string `json:"given_name"`
	LastName      string `json:"family_name"`
}
