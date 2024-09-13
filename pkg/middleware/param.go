package middleware

import (
	"github.com/google/uuid"
)

type Claims struct {
	Data DataClaims `json:"data"`
	Exp  int64      `json:"exp"`
	Iat  int64      `json:"iat"`
}

type DataClaims struct {
	Role   string    `json:"role"`
	UserID uuid.UUID `json:"userID"`
}
