package product

import "github.com/google/uuid"

type ProductListResponse struct {
	ID               uuid.UUID  `json:"ID"`
	UserID           *uuid.UUID `json:"userID"`
	Name             string     `json:"name"`
	ShortDescription *string    `json:"shortDescription"`
	Description      *string    `json:"description"`
	GoalAmount       int        `json:"goalAmount"`
	CurrentAmount    int        `json:"currentAmount"`
	Slug             string     `json:"slug"`
	BackerAmount     int        `json:"backerCount"`
}

type ProductDetailResponse struct {
	ID               uuid.UUID  `json:"ID"`
	UserID           *uuid.UUID `json:"userID"`
	Name             string     `json:"name"`
	ShortDescription *string    `json:"shortDescription"`
	Description      *string    `json:"description"`
	GoalAmount       int        `json:"goalAmount"`
	CurrentAmount    int        `json:"currentAmount"`
	Slug             string     `json:"slug"`
	BackerAmount     int        `json:"backerCount"`
}

type ProductCreateResponse struct {
	UserID           *uuid.UUID `json:"userID"`
	Name             string     `json:"name"`
	ShortDescription *string    `json:"shortDescription"`
	Description      *string    `json:"description"`
	GoalAmount       int        `json:"goalAmount"`
	CurrentAmount    int        `json:"currentAmount"`
	Slug             string     `json:"slug"`
	BackerAmount     int        `json:"backerCount"`
}

type ProductUpdateResponse struct {
	UserID           *uuid.UUID `json:"userID"`
	Name             string     `json:"name"`
	ShortDescription *string    `json:"shortDescription"`
	Description      *string    `json:"description"`
	GoalAmount       int        `json:"goalAmount"`
	CurrentAmount    int        `json:"currentAmount"`
	Slug             string     `json:"slug"`
	BackerAmount     int        `json:"backerCount"`
}
