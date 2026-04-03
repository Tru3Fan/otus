package service_test

import (
	"otus/internal/model"
	"otus/internal/repository"
)

type mockUserRepo struct {
	users   []model.User
	pending []string
	nextID  int
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{nextID: 1}
}

func (r *mockUserRepo) AddUser(u model.User) (model.User, error) {
	u.UserID = r.nextID
	r.nextID++
	r.users = append(r.users, u)
	return u, nil
}

func (r *mockUserRepo) GetUserByID(id int) (model.User, error) {
	for _, u := range r.users {
		if u.UserID == id {
			return u, nil
		}
	}
	return model.User{}, repository.ErrNotFound
}

func (r *mockUserRepo) GetAllUsers() ([]model.User, error) {
	return r.users, nil
}
func (r *mockUserRepo) UpdateUser(id int, updated model.User) (model.User, error) {
	for i, u := range r.users {
		if u.UserID == id {
			updated.UserID = id
			r.users[i] = updated
			return updated, nil
		}
	}
	return model.User{}, repository.ErrNotFound
}

func (r *mockUserRepo) DeleteUser(id int) error {
	for i, u := range r.users {
		if u.UserID == id {
			r.users = append(r.users[:i], r.users[i+1:]...)
			return nil
		}
	}
	return repository.ErrNotFound
}

func (r *mockUserRepo) GetUserByTelegramID(telegramID int64) (model.User, error) {
	for _, u := range r.users {
		if u.TelegramUserID == telegramID {
			return u, nil
		}
	}
	return model.User{}, repository.ErrNotFound
}

func (r *mockUserRepo) AddPendingUser(username string) error {
	r.pending = append(r.pending, username)
	return nil
}

func (r *mockUserRepo) IsPendingUser(username string) (bool, error) {
	for _, u := range r.pending {
		if u == username {
			return true, nil
		}
	}
	return false, nil
}

func (r *mockUserRepo) DeletePendingUser(username string) error {
	for i, u := range r.pending {
		if u == username {
			r.pending = append(r.pending[:i], r.pending[i+1:]...)
			return nil
		}
	}
	return nil
}
