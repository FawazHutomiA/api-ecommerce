package repository

import (
	"gorm.io/gorm"

	userEntity "example/internal/entity"
)

type UserRepository interface {
	Save(user userEntity.User) (userEntity.User, error)
	FindByEmail(email string) (userEntity.User, error)
	FindById(ID int) (userEntity.User, error)
	Update(user userEntity.User) (userEntity.User, error)
	FindAll() ([]userEntity.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func UserNewRepository(db *gorm.DB) *userRepository {
	return &userRepository{db}
}

func (r *userRepository) Save(user userEntity.User) (userEntity.User, error) {
	err := r.db.Create(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *userRepository) FindByEmail(email string) (userEntity.User, error) {
	var user userEntity.User

	err := r.db.Where("email = ?", email).Find(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *userRepository) FindById(ID int) (userEntity.User, error) {
	var user userEntity.User

	err := r.db.Where("id = ?", ID).Find(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *userRepository) Update(user userEntity.User) (userEntity.User, error) {
	err := r.db.Save(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *userRepository) FindAll() ([]userEntity.User, error) {
	var users []userEntity.User

	err := r.db.Order("id desc").Find(&users).Error

	if err != nil {
		return users, err
	}

	return users, nil
}
