package entity

import "github.com/google/uuid"

type PaymentRequest struct {
	ID     uuid.UUID `json:"ID"`
	Amount int       `json:"amount"`
}
