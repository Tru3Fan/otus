package service_test

import (
	"errors"
	"otus/internal/model"
	"otus/internal/repository"
	"otus/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupUserService() service.UserService {
	return service.NewUserService(newMockUserRepo())
}

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name     string
		username string
		wantErr  bool
	}{
		{"valid user", "Ivan", false},
		{"empty username", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := setupUserService()
			u, err := svc.CreateUser(model.User{Username: tt.username})
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.username, u.Username)
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
		{"existing user", 1, false},
		{"non existing user", 999, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := setupUserService()
			svc.CreateUser(model.User{Username: "Ivan"})
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
		{"valid update", 1, "Petr", false},
		{"empty username", 1, "", true},
		{"non existing user", 999, "Petr", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := setupUserService()
			svc.CreateUser(model.User{Username: "Ivan"})
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
	}{
		{"existing user", 1, false},
		{"non existing user", 999, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := setupUserService()
			svc.CreateUser(model.User{Username: "Ivan"})
			err := svc.DeleteUser(tt.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAddAndConfirmPendingUser(t *testing.T) {
	svc := setupUserService()
	err := svc.AddPendingUser("testuser")
	assert.NoError(t, err)

	ok, err := svc.IsPendingUser("testuser")
	assert.NoError(t, err)
	assert.True(t, ok)

	u, err := svc.ConfirmPendingUser(123456, "testuser")
	assert.NoError(t, err)
	assert.Equal(t, "testuser", u.Username)
	assert.Equal(t, int64(123456), u.TelegramUserID)

	ok, err = svc.IsPendingUser("testuser")
	assert.NoError(t, err)
	assert.False(t, ok)
}
