package transaction

import "github.com/google/uuid"

type TransactionListResponse struct {
	ID         uuid.UUID `json:"ID"`
	ProductID  uuid.UUID `json:"productID"`
	UserID     uuid.UUID `json:"userID"`
	Amount     int       `json:"amount"`
	Status     string    `json:"status"`
	Code       *string   `json:"code"`
	PaymentURL string    `json:"paymentURL"`
}

type TransactionDetailResponse struct {
	ID         uuid.UUID `json:"ID"`
	ProductID  uuid.UUID `json:"productID"`
	UserID     uuid.UUID `json:"userID"`
	Amount     int       `json:"amount"`
	Status     string    `json:"status"`
	Code       *string   `json:"code"`
	PaymentURL string    `json:"paymentURL"`
}

type TransactionCreateResponse struct {
	ProductID  uuid.UUID `json:"productID"`
	UserID     uuid.UUID `json:"userID"`
	Amount     int       `json:"amount"`
	Status     string    `json:"status"`
	Code       *string   `json:"code"`
	PaymentURL string    `json:"paymentURL"`
}

type TransactionUpdateResponse struct {
	ProductID  uuid.UUID `json:"productID"`
	UserID     uuid.UUID `json:"userID"`
	Amount     int       `json:"amount"`
	Status     string    `json:"status"`
	Code       *string   `json:"code"`
	PaymentURL string    `json:"paymentURL"`
}
