package service

import (
	"otus/internal/model"
	"otus/internal/repository"
	"otus/internal/repository/logger"
)

type TaskService interface {
	CreateTask(title string, userID int) (model.Task, error)
	GetTask(id int) (model.Task, error)
	GetTasks() ([]model.Task, error)
	UpdateTask(id int, title string, userID int) (model.Task, error)
	DeleteTask(id int) error
}

type taskServiceImpl struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) TaskService {
	return &taskServiceImpl{repo: repo}
}

func (s *taskServiceImpl) CreateTask(title string, userID int) (model.Task, error) {
	if title == "" {
		return model.Task{}, ErrEmptyTitle
	}
	t, err := s.repo.AddTask(model.Task{Title: title, UserID: userID})
	if err != nil {
		return model.Task{}, err
	}
	_ = logger.LogAction("create", "task", t.TaskID)
	return t, nil

}

func (s *taskServiceImpl) GetTask(id int) (model.Task, error) {
	return s.repo.GetTaskByID(id)
}
func (s *taskServiceImpl) GetTasks() ([]model.Task, error) {
	return s.repo.GetAllTasks()
}

func (s *taskServiceImpl) UpdateTask(id int, title string, userID int) (model.Task, error) {
	if title == "" {
		return model.Task{}, ErrEmptyTitle
	}
	t, err := s.repo.UpdateTask(id, model.Task{Title: title, UserID: userID})
	if err != nil {
		return model.Task{}, err
	}
	_ = logger.LogAction("update", "task", t.TaskID)
	return t, nil
}

func (s *taskServiceImpl) DeleteTask(id int) error {
	err := s.repo.DeleteTask(id)
	if err != nil {
		return err
	}
	_ = logger.LogAction("delete", "task", id)
	return nil
}
