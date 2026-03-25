package service

import (
	"otus/internal/model"
	"otus/internal/repository"
)

type TaskService interface {
	CreateTask(title string) (model.Task, error)
	GetTask(id int) (model.Task, error)
	GetTasks() ([]model.Task, error)
	UpdateTask(id int, title string) (model.Task, error)
	DeleteTask(id int) error
}

type taskServiceImpl struct{}

func NewTaskService() TaskService {
	return &taskServiceImpl{}
}

func (u *taskServiceImpl) CreateTask(title string) (model.Task, error) {
	if title == "" {
		return model.Task{}, ErrEmptyTitle
	}
	return repository.AddTask(model.Task{Title: title})
}

func (u *taskServiceImpl) GetTask(id int) (model.Task, error) {
	return repository.GetTaskByID(id)
}
func (u *taskServiceImpl) GetTasks() ([]model.Task, error) {
	return repository.GetAllTasks()
}

func (u *taskServiceImpl) UpdateTask(id int, title string) (model.Task, error) {
	if title == "" {
		return model.Task{}, ErrEmptyTitle
	}
	return repository.UpdateTask(id, model.Task{Title: title})
}

func (u *taskServiceImpl) DeleteTask(id int) error {
	return repository.DeleteTask(id)
}
