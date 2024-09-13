package product

import "github.com/google/uuid"

type ProductCreateRequest struct {
	UserID           *uuid.UUID `json:"userID"`
	Name             string     `json:"name" validate:"required"`
	ShortDescription *string    `json:"shortDescription"`
	Description      *string    `json:"description"`
	GoalAmount       int        `json:"goalAmount"`
	CurrentAmount    int        `json:"currentAmount"`
	Slug             string     `json:"slug"`
	BackerAmount     int        `json:"backerCount"`
}

type ProductUpdateRequest struct {
	UserID           *uuid.UUID `json:"userID"`
	Name             string     `json:"name" validate:"required"`
	ShortDescription *string    `json:"shortDescription"`
	Description      *string    `json:"description"`
	GoalAmount       int        `json:"goalAmount"`
	CurrentAmount    int        `json:"currentAmount"`
	Slug             string     `json:"slug"`
	BackerAmount     int        `json:"backerCount"`
}
