package user

import "github.com/google/uuid"

type UserListResponse struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Occupation *string   `json:"occopation"`
	Phone      *string   `json:"phone"`
	Gender     *string   `json:"gender"`
	Role       string    `json:"role"`
	IsGoogle   bool      `json:"isGoogle"`
	IsActive   bool      `json:"isActive"`
	IsVerify   bool      `json:"isVerify"`
}

type UserDetailResponse struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Occupation *string   `json:"occopation"`
	Phone      *string   `json:"phone"`
	Gender     *string   `json:"gender"`
	Role       string    `json:"role"`
	IsGoogle   bool      `json:"isGoogle"`
	IsActive   bool      `json:"isActive"`
	IsVerify   bool      `json:"isVerify"`
}
