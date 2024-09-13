package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type DataToken struct {
	jwt.StandardClaims
	UserID uuid.UUID `json:"userID"`
	Role   string    `json:"role"`
}

type GenerateResponse struct {
	Token string `json:"token"`
	Exp   int64  `json:"exp"`
}
