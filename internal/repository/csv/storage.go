package csv

import (
	"encoding/csv"
	"errors"
	"os"
	"otus/internal/model"
	"strconv"
)

//const (
//	userFile = "./data/users.csv"
//	taskFile = "./data/tasks.csv"
//)

func userFilePath() string {
	if dir := os.Getenv("DATA_DIR"); dir != "" {
		return dir + "/users.csv"
	}
	return "./data/users.csv"
}

func taskFilePath() string {
	if dir := os.Getenv("DATA_DIR"); dir != "" {
		return dir + "/tasks.csv"
	}
	return "./data/tasks.csv"
}

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
	userData, err := os.Open(userFilePath())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	defer userData.Close()

	//sc := bufio.NewScanner(userData)
	//for sc.Scan() {
	//	var u model.User
	//	if err := json.Unmarshal(sc.Bytes(), &u); err != nil {
	//		return err
	//	}
	rows, err := csv.NewReader(userData).ReadAll()
	if err != nil {
		return err
	}

	for _, row := range rows[1:] {
		id, _ := strconv.Atoi(row[0])
		muUsers.Lock()
		users = append(users, model.User{UserID: id, Username: row[1]})
		muUsers.Unlock()
	}
	return nil
}

func loadTask() error {
	taskData, err := os.Open(taskFilePath())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	defer taskData.Close()

	rows, err := csv.NewReader(taskData).ReadAll()
	if err != nil {
		return err
	}
	for _, row := range rows[1:] {
		id, _ := strconv.Atoi(row[0])
		muTasks.Lock()
		tasks = append(tasks, model.Task{TaskID: id, Title: row[1]})
		muTasks.Unlock()
	}
	return nil
}

func appendCSV(path string, row []string) error {

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()
	return w.Write(row)
}

func saveAllCSV(path string, header []string, rows [][]string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	if err := w.Write(header); err != nil {
		return err
	}
	return w.WriteAll(rows)
}
