package service_test

import (
	"otus/internal/model"
	"otus/internal/repository"
)

type MockTaskRepo struct {
	tasks  []model.Task
	nestID int
}

func NewMockTaskRepo() *MockTaskRepo {
	return &MockTaskRepo{nestID: 1}
}

func (r *MockTaskRepo) AddTask(t model.Task) (model.Task, error) {
	t.TaskID = r.nestID
	r.nestID++
	r.tasks = append(r.tasks, t)
	return t, nil
}

func (r *MockTaskRepo) GetTaskByID(id int) (model.Task, error) {
	for _, t := range r.tasks {
		if t.TaskID == id {
			return t, nil
		}
	}
	return model.Task{}, repository.ErrNotFound
}

func (r *MockTaskRepo) GetAllTasks() ([]model.Task, error) {
	return r.tasks, nil
}

func (r *MockTaskRepo) UpdateTask(id int, updated model.Task) (model.Task, error) {
	for i, t := range r.tasks {
		if t.TaskID == id {
			updated.TaskID = id
			r.tasks[i] = updated
			return updated, nil
		}
	}
	return model.Task{}, repository.ErrNotFound
}

func (r *MockTaskRepo) DeleteTask(id int) error {
	for i, t := range r.tasks {
		if t.TaskID == id {
			r.tasks = append(r.tasks[:i], r.tasks[i+1:]...)
			return nil
		}
	}
	return repository.ErrNotFound
}

func (r *MockTaskRepo) GetTasksByUserID(userId int) ([]model.Task, error) {
	var tasks []model.Task
	for _, t := range r.tasks {
		if t.UserID == userId {
			tasks = append(tasks, t)
		}
	}
	return tasks, nil
}

func (r *MockTaskRepo) GetTasksByStatus(status string) ([]model.Task, error) {
	var result []model.Task
	for _, t := range r.tasks {
		if t.Status == status {
			result = append(result, t)
		}
	}
	return result, nil
}

func (r *MockTaskRepo) GetTasksByAuthorID(authorId int) ([]model.Task, error) {
	var result []model.Task
	for _, t := range r.tasks {
		if t.AssignedBy == authorId {
			result = append(result, t)
		}
	}
	return result, nil
}
