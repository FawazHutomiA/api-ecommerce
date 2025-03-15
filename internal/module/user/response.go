package user

import (
	"time"

	"github.com/google/uuid"
)

type UserCreateResponse struct {
	ID       uuid.UUID `json:"ID"`
	RoleID   uuid.UUID `json:"roleID"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password *string   `json:"password"`
	Phone    string    `json:"phone"`
	Gender   string    `json:"gender"`
	Birth    time.Time `json:"birth"`
	IsActive bool      `json:"isActive"`
	Image    *string   `json:"image"`
}

type UserUpdateResponse struct {
	ID       uuid.UUID `json:"ID"`
	RoleID   uuid.UUID `json:"roleID"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password *string   `json:"password"`
	Phone    string    `json:"phone"`
	Gender   string    `json:"gender"`
	Birth    time.Time `json:"birth"`
	IsActive bool      `json:"isActive"`
}

type UserDetailResponse struct {
	ID           uuid.UUID    `json:"ID"`
	RoleResponse RoleResponse `json:"role"`
	Name         string       `json:"name"`
	Email        string       `json:"email"`
	Phone        string       `json:"phone"`
	Gender       string       `json:"gender"`
	Birth        time.Time    `json:"birth"`
	IsActive     bool         `json:"isActive"`
	Image        *string      `json:"image"`
	CreatedAt    *time.Time   `json:"createdAt"`
}

type UserListlResponse struct {
	ID           uuid.UUID    `json:"ID"`
	RoleResponse RoleResponse `json:"role"`
	Name         string       `json:"name"`
	Email        string       `json:"email"`
	Phone        string       `json:"phone"`
	Gender       string       `json:"gender"`
	Birth        time.Time    `json:"birth"`
	IsActive     bool         `json:"isActive"`
	Image        *string      `json:"image"`
	CreatedAt    *time.Time   `json:"createdAt"`
}

type RoleResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
