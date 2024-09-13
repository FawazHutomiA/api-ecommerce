package transaction

import "github.com/google/uuid"

type TransactionCreateRequest struct {
	ProductID  uuid.UUID `json:"productID"`
	UserID     uuid.UUID `json:"userID"`
	Amount     int       `json:"amount"`
	Status     string    `json:"status"`
	Code       *string   `json:"code"`
	PaymentURL string    `json:"paymentURL"`
}

type TransactionUpdateRequest struct {
	ProductID  uuid.UUID `json:"productID"`
	UserID     uuid.UUID `json:"userID"`
	Amount     int       `json:"amount"`
	Status     string    `json:"status"`
	Code       *string   `json:"code"`
	PaymentURL string    `json:"paymentURL"`
}

type TransactionNotificationInput struct {
	TransactionStatus string    `json:"transaction_status"`
	OrderID           uuid.UUID `json:"order_id"`
	PaymentType       string    `json:"payment_type"`
	FraudStatus       string    `json:"fraud_status"`
}
