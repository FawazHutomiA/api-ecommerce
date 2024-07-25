package user

type UserFormatter struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Occupation string `json:"occupation"`
	Email      string `json:"email"`
	Token      string `json:"token"`
	ImageURL   string `json:"image_url"`
	IsGoogle   bool   `json:"is_google"`
}

type UserGetFormatter struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Occupation string `json:"occupation"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	ImageURL   string `json:"image_url"`
	IsGoogle   bool   `json:"is_google"`
}

func FormatUser(user User, token string) UserFormatter {
	formatter := UserFormatter{
		ID:         user.ID,
		Name:       user.Name,
		Occupation: user.Occupation,
		Email:      user.Email,
		Token:      token,
		ImageURL:   user.AvatarFileName,
		IsGoogle:   user.IsGoogle,
	}

	return formatter
}

func FormatGetUser(user User) UserGetFormatter {
	formatter := UserGetFormatter{
		ID:         user.ID,
		Name:       user.Name,
		Occupation: user.Occupation,
		Email:      user.Email,
		ImageURL:   user.AvatarFileName,
		Role:       user.Role,
		IsGoogle:   user.IsGoogle,
	}

	return formatter
}

func FormatUsers(users []User) []UserGetFormatter {
	usersFormatter := []UserGetFormatter{}

	for _, user := range users {
		userFormatter := FormatGetUser(user)
		usersFormatter = append(usersFormatter, userFormatter)
	}

	return usersFormatter
}
