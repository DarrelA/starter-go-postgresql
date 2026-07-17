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
