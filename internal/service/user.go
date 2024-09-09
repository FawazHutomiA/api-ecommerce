package service

import (
	"errors"

	userEntity "example/internal/entity"
	userInput "example/internal/input"
	userRepository "example/internal/repository"
)

type UserService interface {
	SaveAvatar(ID int, fileLocation string) (userEntity.User, error)
	GetUserByID(ID int) (userEntity.User, error)
	GetUsers() ([]userEntity.User, error)
	GetOrSaveUser(userGoogle userInput.GoogleUser) (userEntity.User, error)
}

type userService struct {
	repository userRepository.UserRepository
}

func UserNewService(repository userRepository.UserRepository) *userService {
	return &userService{repository}
}

func (s *userService) SaveAvatar(ID int, fileLocation string) (userEntity.User, error) {

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

func (s *userService) GetUserByID(ID int) (userEntity.User, error) {

	user, err := s.repository.FindById(ID)
	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, errors.New("no user found with that ID")
	}

	return user, nil
}

func (s *userService) GetUsers() ([]userEntity.User, error) {
	users, err := s.repository.FindAll()
	if err != nil {
		return users, err
	}
	return users, nil
}

func (s *userService) GetOrSaveUser(userGoogle userInput.GoogleUser) (userEntity.User, error) {
	userByEmail, err := s.repository.FindByEmail(userGoogle.Email)

	if err != nil {
		return userByEmail, err
	}

	if userByEmail.ID == 0 {
		user := userEntity.User{
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
