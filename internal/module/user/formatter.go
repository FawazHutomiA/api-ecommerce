package user

import (
	userEntity "example/internal/entity"
)

type UserFormatter struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Occupation string `json:"occupation"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	ImageURL   string `json:"image_url"`
	IsGoogle   bool   `json:"is_google"`
	Token      string `json:"token"`
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

func FormatUser(user userEntity.User, token string) UserFormatter {
	formatter := UserFormatter{
		ID:         user.ID,
		Name:       user.Name,
		Occupation: user.Occupation,
		Role:       user.Role,
		Email:      user.Email,
		ImageURL:   user.AvatarFileName,
		IsGoogle:   user.IsGoogle,
		Token:      token,
	}

	return formatter
}

func FormatGetUser(user userEntity.User) UserGetFormatter {
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

func FormatUsers(users []userEntity.User) []UserGetFormatter {
	usersFormatter := []UserGetFormatter{}

	for _, user := range users {
		userFormatter := FormatGetUser(user)
		usersFormatter = append(usersFormatter, userFormatter)
	}

	return usersFormatter
}
