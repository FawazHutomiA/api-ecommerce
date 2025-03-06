package auth

import (
	"time"

	"github.com/google/uuid"
)

type AuthLoginRequest struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"required,min=8"`
}

type AuthRegisterRequest struct {
	RoleID   uuid.UUID `json:"roleID" validate:"required"`
	Name     string    `json:"name"`
	Email    string    `json:"email" validate:"email"`
	Password string    `json:"password" validate:"required,min=8"`
	Phone    string    `json:"phone"`
	Gender   string    `json:"gender"`
	Birth    time.Time `json:"birth"`
}
