package domain

import (
	"time"
)

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"password,omitempty" validate:"required,min=8"`
	FirstName string    `json:"first_name" validate:"required,alpha"`
	LastName  string    `json:"last_name" validate:"required,alpha"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
