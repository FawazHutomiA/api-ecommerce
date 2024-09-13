package auth

type AuthLoginResponse struct {
	ExpiredAt int64  `json:"exp"`
	Token     string `json:"token"`
}

type AuthRegisterResponse struct {
	ExpiredAt int64  `json:"exp"`
	Token     string `json:"token"`
}
