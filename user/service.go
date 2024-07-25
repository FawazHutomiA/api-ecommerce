package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegisterUser(input RegisterUserInput) (User, error)
	Login(input LoginInput) (User, error)
	IsEmailAvailable(input CheckEmailInput) (bool, error)
	CheckEmail(input CheckEmailInput) (User, error)
	SaveAvatar(ID int, fileLocation string) (User, error)
	GetUserByID(ID int) (User, error)
	GetUsers() ([]User, error)
	GetOrSaveUser(userGoogle GoogleUser) (User, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) RegisterUser(input RegisterUserInput) (User, error) {
	user := User{}
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

func (s *service) Login(input LoginInput) (User, error) {
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

func (s *service) IsEmailAvailable(input CheckEmailInput) (bool, error) {
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

func (s *service) CheckEmail(input CheckEmailInput) (User, error) {
	email := input.Email

	user, err := s.repository.FindByEmail(email)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *service) SaveAvatar(ID int, fileLocation string) (User, error) {

	user, err := s.repository.FindById(ID)
	if err != nil {
		return user, err
	}

	user.AvatarFileName = fileLocation

	updatedUser, err := s.repository.Update(user)
	if err != nil {
		return updatedUser, err
	}

	return updatedUser, nil
}

func (s *service) GetUserByID(ID int) (User, error) {

	user, err := s.repository.FindById(ID)
	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, errors.New("no user found with that ID")
	}

	return user, nil
}

func (s *service) GetUsers() ([]User, error) {
	users, err := s.repository.FindAll()
	if err != nil {
		return users, err
	}
	return users, nil
}

func (s *service) GetOrSaveUser(userGoogle GoogleUser) (User, error) {
	userByEmail, err := s.repository.FindByEmail(userGoogle.Email)

	if err != nil {
		return userByEmail, err
	}

	if userByEmail.ID == 0 {
		user := User{
			Name:       userGoogle.Name,
			Email:      userGoogle.Email,
			Occupation: "",
			IsGoogle:   true,
		}

		newUser, err := s.repository.Save(user)
		if err != nil {
			return newUser, err
		}

		return newUser, nil
	}

	return userByEmail, nil
}
