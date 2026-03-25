package service

import (
	"otus/internal/model"
	"otus/internal/repository"
)

type TaskService interface {
	CreateTask(title string, userID int) (model.Task, error)
	GetTask(id int) (model.Task, error)
	GetTasks() ([]model.Task, error)
	UpdateTask(id int, title string, userID int) (model.Task, error)
	GetTasksByUser(userID int) ([]model.Task, error)
	DeleteTask(id int) error
}

type taskServiceImpl struct{}

func NewTaskService() TaskService {
	return &taskServiceImpl{}
}

func (s *taskServiceImpl) CreateTask(title string, userID int) (model.Task, error) {
	if title == "" {
		return model.Task{}, ErrEmptyTitle
	}

	t, err := repository.PgAddTask(model.Task{Title: title, UserID: userID})
	if err != nil {
		return model.Task{}, err
	}
	_ = repository.LogAction("create", "task", t.TaskID)
	return t, nil

}

func (s *taskServiceImpl) GetTask(id int) (model.Task, error) {
	return repository.PgGetTaskByID(id)
}
func (s *taskServiceImpl) GetTasks() ([]model.Task, error) {
	return repository.PgGetAllTasks()
}
func (s *taskServiceImpl) GetTasksByUser(userID int) ([]model.Task, error) {
	return repository.PgGetTasksByUserID(userID)
}

func (s *taskServiceImpl) UpdateTask(id int, title string, userID int) (model.Task, error) {
	if title == "" {
		return model.Task{}, ErrEmptyTitle
	}
	t, err := repository.PgUpdateTask(id, model.Task{Title: title, UserID: userID})
	if err != nil {
		return model.Task{}, err
	}
	_ = repository.LogAction("update", "task", t.TaskID)
	return t, nil
}

func (s *taskServiceImpl) DeleteTask(id int) error {
	err := repository.PgDeleteTask(id)
	if err != nil {
		return err
	}
	_ = repository.LogAction("delete", "task", id)
	return nil
}
