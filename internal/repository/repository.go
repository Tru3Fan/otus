package repository

import (
	"errors"
	"otus/internal/model"
)

var ErrNotFound = errors.New("not found")

type UserRepository interface {
	AddUser(u model.User) (model.User, error)
	GetUserByID(id int) (model.User, error)
	GetAllUsers() ([]model.User, error)
	UpdateUser(id int, u model.User) (model.User, error)
	DeleteUser(id int) error
}

type TaskRepository interface {
	AddTask(t model.Task) (model.Task, error)
	GetTaskByID(id int) (model.Task, error)
	GetAllTasks() ([]model.Task, error)
	UpdateTask(id int, t model.Task) (model.Task, error)
	DeleteTask(id int) error
	GetTasksByUserID(userID int) ([]model.Task, error)
}
