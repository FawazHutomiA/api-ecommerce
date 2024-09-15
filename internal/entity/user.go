package entity

import (
	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `db:"id" json:"id"`
	Name       string    `db:"name" json:"name"`
	Email      string    `db:"email" json:"email"`
	Occupation *string   `db:"occupation" json:"occupation"`
	Password   *string   `db:"password" json:"password,omitempty"`
	Phone      *string   `db:"phone" json:"phone"`
	Gender     *string   `db:"gender" json:"gender"`
	Role       string    `db:"role" json:"role"`
	Token      string    `db:"token" json:"token"`
	IsGoogle   bool      `db:"is_google" json:"isGoogle"`
	IsActive   bool      `db:"is_active" json:"isActive"`
	IsVerify   bool      `db:"is_verify" json:"isVerify"`
}

func (a *User) ToInsert() []interface{} {
	return []interface{}{
		a.ID,
		a.Name,
		a.Email,
		a.Occupation,
		a.Password,
		a.Phone,
		a.Role,
		a.Gender,
		a.IsGoogle,
		a.Token,
	}
}

func (a *User) ToUpdate() []interface{} {
	return []interface{}{
		a.ID,
		a.Token,
	}
}

func (a *User) ToUpdateVerifyStatus() []interface{} {
	return []interface{}{
		a.ID,
		a.IsVerify,
	}
}
