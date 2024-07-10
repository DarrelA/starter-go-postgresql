package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        int64      `json:"ID"`
	UUID      *uuid.UUID `json:"uuid"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Password  string     `json:"password"`
	Email     string     `json:"email"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

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
