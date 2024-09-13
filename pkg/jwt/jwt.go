package jwt

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateToken(data DataToken) (resp GenerateResponse, err error) {
	expTokenStr := os.Getenv("EXP_TOKEN")
	secretKey := []byte(os.Getenv("SECRET_KEY"))

	expToken, err := strconv.Atoi(expTokenStr)
	if err != nil {
		return resp, err
	}

	expTokenDate := time.Now().Add(time.Second * time.Duration(expToken)).UTC().Unix()
	claim := jwt.MapClaims{
		"exp": expTokenDate,
		"iat": time.Now().Unix(),
		"data": map[string]interface{}{
			"userID": data.UserID,
			"role":   data.Role,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return resp, err
	}

	resp = GenerateResponse{
		Token: signedToken,
		Exp:   expTokenDate,
	}

	return resp, nil
}

func ValidateToken(encodedToken string) (*jwt.Token, error) {
	secretKey := []byte(os.Getenv("SECRET_KEY"))
	// Parse token
	token, err := jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		// Check token method
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("Invalid token")
		}

		return []byte(secretKey), nil
	})

	if err != nil {
		return token, err
	}

	return token, nil
}
