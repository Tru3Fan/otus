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

func (s *userServiceImpl) CreateUser(username string) (model.User, error) {
	if username == "" {
		return model.User{}, ErrEmptyUsername
	}
	u, err := repository.MongoAddUser(model.User{Username: username})
	if err != nil {
		return model.User{}, err
	}
	_ = repository.LogAction("create", "user", u.UserID)
	return u, nil
}

func (s *userServiceImpl) GetUser(id int) (model.User, error) {
	return repository.MongoGetUserByID(id)
}
func (s *userServiceImpl) GetUsers() ([]model.User, error) {
	return repository.MongoGetAllUsers()
}
func (s *userServiceImpl) UpdateUser(id int, username string) (model.User, error) {
	if username == "" {
		return model.User{}, ErrEmptyUsername
	}
	u, err := repository.MongoUpdateUser(id, model.User{Username: username})
	if err != nil {
		return model.User{}, err
	}
	_ = repository.LogAction("update", "user", id)
	return u, nil
}
func (s *userServiceImpl) DeleteUser(id int) error {
	err := repository.MongoDeleteUser(id)
	if err != nil {
		return err
	}
	_ = repository.LogAction("delete", "user", id)
	return nil
}
