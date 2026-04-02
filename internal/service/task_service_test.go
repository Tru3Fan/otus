package service_test

import (
	"errors"
	"otus/internal/repository"
	"otus/internal/service"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func setupTaskService() service.TaskService {
	return service.NewTaskService(NewMockTaskRepo())
}

func TestCreateTask(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		wantErr bool
	}{
		{"valid task", "Dream", false},
		{"empty title", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := setupTaskService()
			task, err := svc.CreateTask(tt.title, 0)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.title, task.Title)
				assert.NotZero(t, task.TaskID)
			}
		})
	}
}

func TestGetTask(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		wantErr bool
	}{
		{"existing task", 1, false},
		{"non existing task", 999, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := setupTaskService()
			svc.CreateTask("Dream", 0)
			_, err := svc.GetTask(tt.id)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, repository.ErrNotFound))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdateTaskStatus(t *testing.T) {
	tests := []struct {
		name    string
		status  string
		wantErr bool
	}{
		{"pending", "pending", false},
		{"in_progress", "in_progress", false},
		{"done", "done", false},
		{"cancelled", "cancelled", false},
		{"invalid", "неверный", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := setupTaskService()
			task, _ := svc.CreateTask("Dream", 0)
			updated, err := svc.UpdateTaskStatus(task.TaskID, tt.status)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.status, updated.Status)
			}
		})
	}
}

func TestCloseTask(t *testing.T) {
	svc := setupTaskService()
	task, _ := svc.CreateTask("Dream", 0)
	closed, err := svc.CloseTask(task.TaskID)
	assert.NoError(t, err)
	assert.Equal(t, "done", closed.Status)
}

func TestCreateTaskFull(t *testing.T) {
	svc := setupTaskService()
	d := time.Now().AddDate(0, 0, 3)
	task, err := svc.CreateTaskFull("Fix bug", "описание", 2, 1, &d)
	assert.NoError(t, err)
	assert.Equal(t, "Fix bug", task.Title)
	assert.Equal(t, "pending", task.Status)
	assert.NotZero(t, task.TaskID)
}

func TestDeleteTask(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		wantErr bool
	}{
		{"existing task", 1, false},
		{"non existing task", 999, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := setupTaskService()
			svc.CreateTask("Dream", 0)
			err := svc.DeleteTask(tt.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
