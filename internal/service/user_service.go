package service

import (
	"otus/internal/model"
	"otus/internal/repository"
	"otus/internal/repository/logger"
)

type UserService interface {
	CreateUser(u model.User) (model.User, error)
	GetUser(id int) (model.User, error)
	GetUsers() ([]model.User, error)
	UpdateUser(id int, username string) (model.User, error)
	DeleteUser(id int) error
	GetUserByTelegramID(telegramID int64) (model.User, error)
}

type userServiceImpl struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userServiceImpl{repo: repo}
}

func (s *userServiceImpl) CreateUser(u model.User) (model.User, error) {
	if u.Username == "" {
		return model.User{}, ErrEmptyUsername
	}
	created, err := s.repo.AddUser(u)
	if err != nil {
		return model.User{}, err
	}
	_ = logger.LogAction("create", "user", created.UserID)
	return created, nil
}

func (s *userServiceImpl) GetUser(id int) (model.User, error) {
	return s.repo.GetUserByID(id)
}
func (s *userServiceImpl) GetUsers() ([]model.User, error) {
	return s.repo.GetAllUsers()
}
func (s *userServiceImpl) UpdateUser(id int, username string) (model.User, error) {
	if username == "" {
		return model.User{}, ErrEmptyUsername
	}
	u, err := s.repo.UpdateUser(id, model.User{Username: username})
	if err != nil {
		return model.User{}, err
	}
	_ = logger.LogAction("update", "user", id)
	return u, nil
}
func (s *userServiceImpl) DeleteUser(id int) error {
	err := s.repo.DeleteUser(id)
	if err != nil {
		return err
	}
	_ = logger.LogAction("delete", "user", id)
	return nil
}
func (s *userServiceImpl) GetUserByTelegramID(telegramID int64) (model.User, error) {
	return s.repo.GetUserByTelegramID(telegramID)
}
