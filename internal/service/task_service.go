package service

import (
	"otus/internal/model"
	"otus/internal/repository"
	"otus/internal/repository/logger"
	"time"
)

type TaskService interface {
	CreateTask(title string, userID int) (model.Task, error)
	GetTask(id int) (model.Task, error)
	GetTasks() ([]model.Task, error)
	UpdateTask(id int, title string, userID int) (model.Task, error)
	DeleteTask(id int) error
	GetTasksByUser(userID int) ([]model.Task, error)
	GetTasksByStatus(status string) ([]model.Task, error)
	UpdateTaskStatus(id int, status string) (model.Task, error)

	CreateTaskFull(title string, assigneeID, authorID int, deadline *time.Time) (model.Task, error)
	GetTasksByAuthor(authorID int) ([]model.Task, error)
	CloseTask(id int) (model.Task, error)
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
	t, err := s.repo.AddTask(model.Task{Title: title})
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
	t, err := s.repo.UpdateTask(id, model.Task{Title: title})
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
func (s *taskServiceImpl) GetTasksByUser(userID int) ([]model.Task, error) {
	return s.repo.GetTasksByUserID(userID)
}

func (s *taskServiceImpl) GetTasksByStatus(status string) ([]model.Task, error) {
	return s.repo.GetTasksByStatus(status)
}

func (s *taskServiceImpl) UpdateTaskStatus(id int, status string) (model.Task, error) {
	if status != "pending" && status != "in_progress" && status != "done" && status != "cancelled" {
		return model.Task{}, ErrInvalidStatus
	}
	t, err := s.repo.GetTaskByID(id)
	if err != nil {
		return model.Task{}, err
	}
	t.Status = status
	return s.repo.UpdateTask(id, t)
}

func (s *taskServiceImpl) CreateTaskFull(title string, assigneeID, authorID int, deadline *time.Time) (model.Task, error) {
	if title == "" {
		return model.Task{}, ErrEmptyTitle
	}
	return s.repo.AddTask(model.Task{
		Title:      title,
		UserID:     assigneeID,
		AssignedBy: authorID,
		Status:     "pending",
		Deadline:   deadline,
	})

}

func (s *taskServiceImpl) GetTasksByAuthor(authorID int) ([]model.Task, error) {
	return s.repo.GetTasksByAuthorID(authorID)
}

func (s *taskServiceImpl) CloseTask(id int) (model.Task, error) {
	return s.UpdateTaskStatus(id, "done")
}
