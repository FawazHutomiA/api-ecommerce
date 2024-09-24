package user

import (
	"time"

	"github.com/google/uuid"
)

type UserCreateResponse struct {
	ID          uuid.UUID `json:"ID"`
	RoleID      uuid.UUID `json:"roleID"`
	WarehouseID uuid.UUID `json:"warehouseID"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Password    *string   `json:"password"`
	Phone       string    `json:"phone"`
	Gender      string    `json:"gender"`
	Birth       time.Time `json:"birth"`
	IsActive    bool      `json:"isActive"`
}

type UserUpdateResponse struct {
	ID          uuid.UUID `json:"ID"`
	RoleID      uuid.UUID `json:"roleID"`
	WarehouseID uuid.UUID `json:"warehouseID"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Password    *string   `json:"password"`
	Phone       string    `json:"phone"`
	Gender      string    `json:"gender"`
	Birth       time.Time `json:"birth"`
	IsActive    bool      `json:"isActive"`
}

type UserDetailResponse struct {
	ID          uuid.UUID `json:"ID"`
	RoleID      uuid.UUID `json:"roleID"`
	WarehouseID uuid.UUID `json:"warehouseID"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Password    *string   `json:"password"`
	Phone       string    `json:"phone"`
	Gender      string    `json:"gender"`
	Birth       time.Time `json:"birth"`
	IsActive    bool      `json:"isActive"`
}
