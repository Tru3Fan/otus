package service

import (
	"otus/internal/model"
	"otus/internal/repository"
)

type UserService interface {
	CreateUser(username string) (model.User, error)
	GetUser(id int) (model.User, error)
	GetUsers() ([]model.User, error)
	UpdateUser(id int, username string) (model.User, error)
	DeleteUser(id int) error
}

type userServiceImpl struct{}

func NewUserService() UserService {
	return &userServiceImpl{}
}

func (u *userServiceImpl) CreateUser(username string) (model.User, error) {
	if username == "" {
		return model.User{}, ErrEmptyUsername
	}
	return repository.AddUser(model.User{Username: username})
}

func (u *userServiceImpl) GetUser(id int) (model.User, error) {
	return repository.GetUserByID(id)
}
func (u *userServiceImpl) GetUsers() ([]model.User, error) {
	return repository.GetAllUsers()
}
func (u *userServiceImpl) UpdateUser(id int, username string) (model.User, error) {
	if username == "" {
		return model.User{}, ErrEmptyUsername
	}
	return repository.UpdateUser(id, model.User{Username: username})
}
func (u *userServiceImpl) DeleteUser(id int) error {
	return repository.DeleteUser(id)
}
