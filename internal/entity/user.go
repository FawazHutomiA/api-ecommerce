package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID  `db:"id" json:"ID"`
	RoleID    uuid.UUID  `db:"role_id" json:"roleID"`
	Name      string     `db:"name" json:"name"`
	Email     string     `db:"email" json:"email"`
	Password  *string    `db:"password" json:"password,omitempty"`
	Phone     string     `db:"phone" json:"phone"`
	Gender    string     `db:"gender" json:"gender"`
	Birth     time.Time  `db:"birth" json:"birth"`
	IsActive  bool       `db:"is_active" json:"isActive"`
	Image     *string    `db:"image" json:"image"`
	CreatedAt *time.Time `db:"created_at" json:"createdAt"`
}

type UserRole struct {
	ID       uuid.UUID `db:"id" json:"id"`
	Email    string    `db:"email" json:"email"`
	Role     string    `db:"role" json:"role"`
	Password *string   `db:"password" json:"password"`
}

func (a *User) ToInsert() []interface{} {
	return []interface{}{
		a.ID,
		a.RoleID,
		a.Name,
		a.Email,
		a.Password,
		a.Phone,
		a.Gender,
		a.Birth,
		a.IsActive,
		a.Image,
	}
}

func (a *User) ToUpdate() []interface{} {
	return []interface{}{
		a.ID,
		a.RoleID,
		a.Name,
		a.Email,
		a.Password,
		a.Phone,
		a.Gender,
		a.Birth,
		a.IsActive,
	}
}
