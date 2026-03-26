package service_test

import (
	"errors"
	"os"
	"otus/internal/model"
	"otus/internal/repository"
	"otus/internal/repository/csv"
	"otus/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockUserRepo struct {
	users []model.User
}

func setupUserService() service.UserService {
	os.Setenv("DATA_DIR", "../../data")
	csv.ResetUsers()
	repo := csv.NewUserRepo()
	return service.NewUserService(repo)
}

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name     string
		username string
		wantErr  bool
		wantUser model.User
	}{
		{
			name:     "valid user",
			username: "Ivan",
			wantErr:  false,
			wantUser: model.User{Username: "Ivan"},
		},
		{
			name:     "empty username",
			username: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := setupUserService()
			u, err := svc.CreateUser(tt.username)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantUser.Username, u.Username)
				assert.NotZero(t, u.UserID)
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		wantErr bool
	}{
		{
			name:    "existing user",
			id:      1,
			wantErr: false,
		},
		{
			name:    "non existing user",
			id:      999,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := setupUserService()
			svc.CreateUser("Ivan")
			_, err := svc.GetUser(tt.id)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, repository.ErrNotFound))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	tests := []struct {
		name     string
		id       int
		username string
		wantErr  bool
	}{
		{
			name:     "valid update",
			id:       1,
			username: "Petr",
			wantErr:  false,
		},
		{
			name:     "empty username",
			id:       1,
			username: "",
			wantErr:  true,
		},
		{
			name:     "non existing user",
			id:       999,
			username: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := setupUserService()
			svc.CreateUser("Ivan")

			u, err := svc.UpdateUser(tt.id, tt.username)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.username, u.Username)
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		wantErr bool
	}{{
		name:    "existing user",
		id:      1,
		wantErr: false,
	}, {
		name:    "non existing user",
		id:      999,
		wantErr: true,
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := setupUserService()
			svc.CreateUser("Ivan")

			err := svc.DeleteUser(tt.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetUsers(t *testing.T) {
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
			name:      "multiple users",
			seedCount: 3,
			wantCount: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := setupUserService()
			for i := 0; i < tt.seedCount; i++ {
				svc.CreateUser("Ivan")
			}
			users, err := svc.GetUsers()
			assert.NoError(t, err)
			assert.Len(t, users, tt.seedCount)
		})
	}
}
