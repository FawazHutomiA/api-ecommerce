package entity

import "github.com/google/uuid"

type Role struct {
	ID          uuid.UUID `db:"id" json:"ID"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
}

func (a *Role) ToInsert() []interface{} {
	return []interface{}{
		a.ID,
		a.Name,
		a.Description,
	}
}

func (a *Role) ToUpdate() []interface{} {
	return []interface{}{
		a.Name,
		a.Description,
	}
}
