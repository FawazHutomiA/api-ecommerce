package auth

type AuthLoginRequest struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"required,min=8"`
}

type AuthRegisterRequest struct {
	Name       string  `json:"name"`
	Email      string  `json:"email" validate:"email"`
	Occupation *string `json:"occupation"`
	Password   *string `json:"password" validate:"required,min=8"`
	Phone      *string `json:"phone"`
	Role       string  `json:"role"`
	Gender     *string `json:"gender"`
	IsGoogle   bool    `json:"isGoogle"`
}

type GoogleUser struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Email         string `json:"email" binding:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
}
