package entity

import (
	"github.com/google/uuid"
)

type Transaction struct {
	ID         uuid.UUID `db:"id" json:"ID"`
	ProductID  uuid.UUID `db:"product_id" json:"productID"`
	UserID     uuid.UUID `db:"user_id" json:"userID"`
	Amount     int       `db:"amount" json:"amount"`
	Status     string    `db:"status" json:"status"`
	Code       *string   `db:"code" json:"code"`
	PaymentURL string    `db:"payment_url" json:"paymentURL"`
}

func (a *Transaction) ToInsert() []interface{} {
	return []interface{}{
		a.ID,
		a.ProductID,
		a.UserID,
		a.Amount,
		a.Status,
		a.Code,
		a.PaymentURL,
	}
}

func (a *Transaction) ToUpdate() []interface{} {
	return []interface{}{
		a.ID,
		a.ProductID,
		a.UserID,
		a.Amount,
		a.Status,
		a.Code,
		a.PaymentURL,
	}
}
