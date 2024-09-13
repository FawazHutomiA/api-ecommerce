package entity

import (
	"github.com/google/uuid"
)

type Product struct {
	ID               uuid.UUID  `db:"id" json:"ID"`
	UserID           *uuid.UUID `db:"user_id" json:"userID"`
	Name             string     `db:"name" json:"name"`
	ShortDescription *string    `db:"short_description" json:"shortDescription"`
	Description      *string    `db:"description" json:"description"`
	GoalAmount       int        `db:"goal_amount" json:"goalAmount"`
	CurrentAmount    int        `db:"current_amount" json:"currentAmount"`
	Slug             string     `db:"slug" json:"slug"`
	BackerAmount     int        `db:"backer_amount" json:"backerAmount"`
}

func (a *Product) ToInsert() []interface{} {
	return []interface{}{
		a.ID,
		a.UserID,
		a.Name,
		a.ShortDescription,
		a.Description,
		a.GoalAmount,
		a.CurrentAmount,
		a.Slug,
		a.BackerAmount,
	}
}

func (a *Product) ToUpdate() []interface{} {
	return []interface{}{
		a.ID,
		a.UserID,
		a.Name,
		a.ShortDescription,
		a.Description,
		a.GoalAmount,
		a.CurrentAmount,
		a.Slug,
		a.BackerAmount,
	}
}
