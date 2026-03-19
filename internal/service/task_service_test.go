package service_test

import (
	"errors"
	"os"
	"otus/internal/repository"
	"otus/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupTaskService() service.TaskService {
	os.Setenv("DATA_DIR", "../../data")
	repository.ResetTasks()
	return service.NewTaskService()
}

func TestCreateTask(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		wantErr bool
	}{
		{
			name:    "valid task",
			title:   "Dream",
			wantErr: false,
		},
		{
			name:    "empty title",
			title:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := setupTaskService()
			task, err := svc.CreateTask(tt.title)
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
		{
			name:    "existing task",
			id:      1,
			wantErr: false,
		},
		{
			name:    "non existing task",
			id:      999,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := setupTaskService()
			svc.CreateTask("Dream")

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

func TestUpdateTask(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		title   string
		wantErr bool
	}{
		{
			name:    "valid update",
			id:      1,
			title:   "New Dream",
			wantErr: false,
		},
		{
			name:    "empty title",
			id:      1,
			title:   "",
			wantErr: true,
		},
		{
			name:    "non existing task",
			id:      999,
			title:   "New Dream",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := setupTaskService()
			svc.CreateTask("Dream")

			task, err := svc.UpdateTask(tt.id, tt.title)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.title, task.Title)
			}
		})
	}
}

func TestDeleteTask(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		wantErr bool
	}{
		{
			name:    "existing task",
			id:      1,
			wantErr: false,
		},
		{
			name:    "non existing task",
			id:      999,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := setupTaskService()
			svc.CreateTask("Dream")

			err := svc.DeleteTask(tt.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetTasks(t *testing.T) {
	tests := []struct {
		name      string
		seedCount int
		wantCount int
	}{
		{
			name:      "empty list",
			seedCount: 0,
			wantCount: 0,
		},
		{
			name:      "multiple tasks",
			seedCount: 3,
			wantCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := setupTaskService()
			for i := 0; i < tt.seedCount; i++ {
				svc.CreateTask("Dream")
			}

			tasks, err := svc.GetTasks()
			assert.NoError(t, err)
			assert.Len(t, tasks, tt.wantCount)
		})
	}
}
