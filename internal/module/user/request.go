package user

import (
	"time"

	"github.com/google/uuid"
)

type UserCreateRequest struct {
	ID       uuid.UUID `json:"ID"`
	RoleID   uuid.UUID `json:"roleID"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password *string   `json:"password" validate:"required,min=8"`
	Phone    string    `json:"phone" validate:"min=11,max=13"`
	Gender   string    `json:"gender"`
	Birth    time.Time `json:"birth"`
	IsActive bool      `json:"isActive"`
	Image    *string   `json:"image"`
}

type UserUpdateRequest struct {
	RoleID   uuid.UUID `json:"roleID"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password *string   `json:"password" validate:"required,min=8"`
	Phone    string    `json:"phone" validate:"min=11,max=13"`
	Gender   string    `json:"gender"`
	Birth    time.Time `json:"birth"`
	IsActive bool      `json:"isActive"`
}
