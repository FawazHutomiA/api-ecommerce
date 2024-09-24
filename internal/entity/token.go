package entity

import (
	"time"

	"github.com/google/uuid"
)

// properties of user
// json attributes will set the field name
type Token struct {
	ID        uuid.UUID `db:"id" json:"ID"`
	UserID    uuid.UUID `db:"user_id" json:"userID"`
	Token     string    `db:"token" json:"token"`
	ExpiredAt time.Time `db:"expired_at" json:"expiredAt"`
}

// dto
func (a *Token) ToInsert() []interface{} {
	return []interface{}{
		a.ID,
		a.UserID,
		a.Token,
		a.ExpiredAt,
	}
}

func (a *Token) ToUpdate() []interface{} {
	return []interface{}{
		a.ID,
		a.UserID,
		a.Token,
		a.ExpiredAt,
	}
}
