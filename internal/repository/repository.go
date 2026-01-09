package repository

import (
	"fmt"
	"otus/internal/model"
)

type Storable interface {
	ID() int
}

var (
	users []model.User
	tasks []model.Task
)

func Add(m Storable) {
	switch v := m.(type) {
	case model.User:
		users = append(users, v)
		fmt.Println("Added User:", v.Username)
	case model.Task:
		tasks = append(tasks, v)
		fmt.Println("Added Task:", v.Title)
	default:
		fmt.Println("Unknown Type")
	}

}
