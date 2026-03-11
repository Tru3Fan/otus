package repository

import (
	"context"
	"errors"
	"fmt"
	"otus/internal/model"
	"sync"
)

type Storable interface {
	ID() int
}

var (
	muUsers sync.Mutex
	muTasks sync.Mutex

	users []model.User
	tasks []model.Task
)

var ErrNotFound = errors.New("not found")

func Add(ctx context.Context, in <-chan Storable, wg *sync.WaitGroup) {

	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case m, ok := <-in:
			if !ok {
				return
			}
			switch v := m.(type) {
			case model.User:
				muUsers.Lock()
				users = append(users, v)
				_ = appendJSON(userFile, v)
				muUsers.Unlock()
				fmt.Println("Added user: ", v.Username)
			case model.Task:
				muTasks.Lock()
				tasks = append(tasks, v)
				_ = appendJSON(taskFile, v)
				muTasks.Unlock()
				fmt.Println("Added task: ", v.Title)
			default:
				fmt.Println("Unknown type: ", v)
			}
		}
	}
}

func GetAllUsers() []model.User {
	muUsers.Lock()
	defer muUsers.Unlock()
	result := make([]model.User, len(users))
	copy(result, users)
	return result
}

func GetUserByID(id int) (model.User, error) {
	muUsers.Lock()
	defer muUsers.Unlock()
	for _, u := range users {
		if u.UserID == id {
			return u, nil
		}
	}
	return model.User{}, ErrNotFound
}

func AddUser(u model.User) error {
	muUsers.Lock()
	defer muUsers.Unlock()

	if u.UserID == 0 {
		u.UserID = nextUserID()
	}
	users = append(users, u)
	return saveAllJSON(userFile, users)
}

func UpdateUser(id int, updated model.User) (model.User, error) {
	muUsers.Lock()
	defer muUsers.Unlock()
	for i, u := range users {
		if u.UserID == id {
			updated.UserID = id
			users[i] = updated
			return updated, saveAllJSON(userFile, users)
		}
	}
	return model.User{}, ErrNotFound
}

func DeleteUser(id int) error {
	muUsers.Lock()
	defer muUsers.Unlock()
	for i, u := range users {
		if u.UserID == id {
			users = append(users[:i], users[i+1:]...)
			return saveAllJSON(userFile, users)
		}
	}
	return ErrNotFound
}

func GetAllTasks() []model.Task {
	muTasks.Lock()
	defer muTasks.Unlock()
	result := make([]model.Task, len(tasks))
	copy(result, tasks)
	return result
}

func GetTaskByID(id int) (model.Task, error) {
	muTasks.Lock()
	defer muTasks.Unlock()
	for _, t := range tasks {
		if t.TaskID == id {
			return t, nil
		}
	}
	return model.Task{}, ErrNotFound
}

func AddTask(t model.Task) error {
	muTasks.Lock()
	defer muTasks.Unlock()
	if t.TaskID == 0 {
		t.TaskID = nextTaskID()
	}
	tasks = append(tasks, t)
	return saveAllJSON(taskFile, tasks)
}

func UpdateTask(id int, updated model.Task) (model.Task, error) {
	muTasks.Lock()
	defer muTasks.Unlock()
	for i, t := range tasks {
		if t.TaskID == id {
			updated.TaskID = id
			tasks[i] = updated
			return updated, saveAllJSON(taskFile, tasks)
		}
	}
	return model.Task{}, ErrNotFound
}

func DeleteTask(id int) error {
	muTasks.Lock()
	defer muTasks.Unlock()
	for i, t := range tasks {
		if t.TaskID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return saveAllJSON(taskFile, tasks)
		}
	}
	return ErrNotFound
}

func nextUserID() int {
	max := 0
	for _, u := range users {
		if u.UserID > max {
			max = u.UserID
		}
	}
	return max + 1
}

func nextTaskID() int {
	max := 0
	for _, t := range tasks {
		if t.TaskID > max {
			max = t.TaskID
		}
	}
	return max + 1
}
