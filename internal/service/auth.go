package service

import (
	"errors"

	userEntity "example/internal/entity"
	userInput "example/internal/input"
	userRepository "example/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	RegisterUser(input userInput.RegisterUserInput) (userEntity.User, error)
	Login(input userInput.LoginInput) (userEntity.User, error)
	IsEmailAvailable(input userInput.CheckEmailInput) (bool, error)
	CheckEmail(input userInput.CheckEmailInput) (userEntity.User, error)
}

type authService struct {
	repository userRepository.UserRepository
}

func AuthNewService(repository userRepository.UserRepository) *authService {
	return &authService{repository}
}

func (s *authService) RegisterUser(input userInput.RegisterUserInput) (userEntity.User, error) {
	user := userEntity.User{}
	user.Name = input.Name
	user.Occupation = input.Occupation
	user.Email = input.Email
	user.IsGoogle = input.IsGoogle
	if !input.IsGoogle {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
		if err != nil {
			return user, err
		}
		user.PasswordHash = string(passwordHash)
		user.Role = "user"
	} else {
		user.PasswordHash = ""
	}

	newUser, err := s.repository.Save(user)
	if err != nil {
		return newUser, err
	}

	return newUser, nil
}

func (s *authService) Login(input userInput.LoginInput) (userEntity.User, error) {
	email := input.Email
	password := input.Password

	user, err := s.repository.FindByEmail(email)
	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, errors.New("no user found on that email")
	}

	if !user.IsGoogle {
		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
		if err != nil {
			return user, err
		}
	}

	return user, nil
}

func (s *authService) IsEmailAvailable(input userInput.CheckEmailInput) (bool, error) {
	email := input.Email

	user, err := s.repository.FindByEmail(email)
	if err != nil {
		return false, err
	}

	if user.ID == 0 {
		return true, nil
	}

	return false, nil
}

func (s *authService) CheckEmail(input userInput.CheckEmailInput) (userEntity.User, error) {
	email := input.Email

	user, err := s.repository.FindByEmail(email)
	if err != nil {
		return user, err
	}

	return user, nil
}
