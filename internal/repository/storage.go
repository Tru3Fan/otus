package repository

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"otus/internal/model"
)

const (
	userFile = "./data/users.json"
	TaskFile = "./data/tasks.json"
)

func LoadAllData() error {
	if err := loadUser(); err != nil {
		return err
	}
	if err := loadTask(); err != nil {
		return err
	}
	return nil
}

func loadUser() error {
	userData, err := os.Open(userFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	defer userData.Close()

	sc := bufio.NewScanner(userData)
	for sc.Scan() {
		var u model.User
		if err := json.Unmarshal(sc.Bytes(), &u); err != nil {
			return err
		}
		muUsers.Lock()
		users = append(users, u)
		muUsers.Unlock()
	}
	return sc.Err()
}

func loadTask() error {
	TaskData, err := os.Open(TaskFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	defer TaskData.Close()

	sc := bufio.NewScanner(TaskData)
	for sc.Scan() {
		var t model.Task
		if err := json.Unmarshal(sc.Bytes(), &t); err != nil {
			return err
		}
		muTasks.Lock()
		tasks = append(tasks, t)
		muTasks.Unlock()
	}
	return sc.Err()
}

func appendJSON(path string, v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(append(data, '\n'))
	return err
}
