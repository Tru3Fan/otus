package csv

import (
	"context"
	"fmt"
	"otus/internal/model"
	"otus/internal/repository"
	"strconv"
	"sync"
)

type UserRepo struct{}
type TaskRepo struct{}

func NewUserRepo() repository.UserRepository {
	return &UserRepo{}
}

func NewTaskRepo() repository.TaskRepository {
	return &TaskRepo{}
}

type Storable interface {
	ID() int
}

var (
	muUsers sync.Mutex
	muTasks sync.Mutex

	users []model.User
	tasks []model.Task
)

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
				_ = appendCSV(userFilePath(), []string{strconv.Itoa(v.UserID), v.Username})
				muUsers.Unlock()
				fmt.Println("Added user: ", v.Username)
			case model.Task:
				muTasks.Lock()
				tasks = append(tasks, v)
				_ = appendCSV(taskFilePath(), []string{strconv.Itoa(v.TaskID), v.Title})
				muTasks.Unlock()
				fmt.Println("Added task: ", v.Title)
			default:
				fmt.Println("Unknown type: ", v)
			}
		}
	}
}
func (r *UserRepo) AddUser(u model.User) (model.User, error) {
	muUsers.Lock()
	defer muUsers.Unlock()

	if u.UserID == 0 {
		u.UserID = nextUserID()
	}
	users = append(users, u)
	return u, saveAllCSV(userFilePath(), []string{"user_id", "username"}, usersToRows())
}

func (r *UserRepo) GetUserByID(id int) (model.User, error) {
	muUsers.Lock()
	defer muUsers.Unlock()
	for _, u := range users {
		if u.UserID == id {
			return u, nil
		}
	}
	return model.User{}, repository.ErrNotFound
}

func (r *UserRepo) GetAllUsers() ([]model.User, error) {
	muUsers.Lock()
	defer muUsers.Unlock()
	result := make([]model.User, len(users))
	copy(result, users)
	return result, nil
}

func (r *UserRepo) UpdateUser(id int, updated model.User) (model.User, error) {
	muUsers.Lock()
	defer muUsers.Unlock()
	for i, u := range users {
		if u.UserID == id {
			updated.UserID = id
			users[i] = updated
			return updated, saveAllCSV(userFilePath(), []string{"user_id", "username"}, usersToRows())
		}
	}
	return model.User{}, repository.ErrNotFound
}

func (r *UserRepo) DeleteUser(id int) error {
	muUsers.Lock()
	defer muUsers.Unlock()
	for i, u := range users {
		if u.UserID == id {
			users = append(users[:i], users[i+1:]...)
			return saveAllCSV(userFilePath(), []string{"user_id", "username"}, usersToRows())
		}
	}
	return repository.ErrNotFound
}

func (r *UserRepo) GetUserByTelegramID(telegramID int64) (model.User, error) {
	muUsers.Lock()
	defer muUsers.Unlock()
	for _, u := range users {
		if u.TelegramUserID == telegramID {
			return u, nil
		}
	}
	return model.User{}, repository.ErrNotFound
}

func (r *TaskRepo) AddTask(t model.Task) (model.Task, error) {
	muTasks.Lock()
	defer muTasks.Unlock()
	if t.TaskID == 0 {
		t.TaskID = nextTaskID()
	}
	tasks = append(tasks, t)
	return t, saveAllCSV(taskFilePath(), []string{"task_id", "title"}, tasksToRows())
}

func (r *TaskRepo) GetTaskByID(id int) (model.Task, error) {
	muTasks.Lock()
	defer muTasks.Unlock()
	for _, t := range tasks {
		if t.TaskID == id {
			return t, nil
		}
	}
	return model.Task{}, repository.ErrNotFound
}

func (r *TaskRepo) GetAllTasks() ([]model.Task, error) {
	muTasks.Lock()
	defer muTasks.Unlock()
	result := make([]model.Task, len(tasks))
	copy(result, tasks)
	return result, nil
}

func (r *TaskRepo) UpdateTask(id int, updated model.Task) (model.Task, error) {
	muTasks.Lock()
	defer muTasks.Unlock()
	for i, t := range tasks {
		if t.TaskID == id {
			updated.TaskID = id
			tasks[i] = updated
			return updated, saveAllCSV(taskFilePath(), []string{"task_id", "title"}, tasksToRows())
		}
	}
	return model.Task{}, repository.ErrNotFound
}

func (r *TaskRepo) DeleteTask(id int) error {
	muTasks.Lock()
	defer muTasks.Unlock()
	for i, t := range tasks {
		if t.TaskID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return saveAllCSV(taskFilePath(), []string{"task_id", "title"}, tasksToRows())
		}
	}
	return repository.ErrNotFound
}

func (r *TaskRepo) GetTasksByUserID(userID int) ([]model.Task, error) {
	muTasks.Lock()
	defer muTasks.Unlock()
	var result []model.Task
	for _, t := range tasks {
		if t.UserID == userID {
			result = append(result, t)
		}
	}
	return result, nil
}

func (r *TaskRepo) GetTasksByStatus(status string) ([]model.Task, error) {
	muTasks.Lock()
	defer muTasks.Unlock()
	var result []model.Task
	for _, t := range tasks {
		if t.Status == status {
			result = append(result, t)
		}
	}
	return result, nil
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

func usersToRows() [][]string {
	rows := make([][]string, len(users))
	for i, u := range users {
		rows[i] = []string{strconv.Itoa(u.UserID), u.Username}
	}
	return rows
}

func tasksToRows() [][]string {
	rows := make([][]string, len(tasks))
	for i, u := range tasks {
		rows[i] = []string{strconv.Itoa(u.TaskID), u.Title}
	}
	return rows
}

func ResetUsers() {
	muUsers.Lock()
	defer muUsers.Unlock()
	users = []model.User{}
}

func ResetTasks() {
	muTasks.Lock()
	defer muTasks.Unlock()
	tasks = []model.Task{}

}
