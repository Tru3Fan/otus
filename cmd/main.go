package main

import (
	"fmt"
	"otus/internal/model/task"
	"otus/internal/model/user"
)

func main() {
	a := task.Task{
		ID:          1,
		Title:       "Ivan",
		Description: "Ivan",
		AssigneeID:  23,
	}

	b := user.User{
		ID:       8654,
		Username: "john",
	}

	fmt.Println(a)
	fmt.Println(b)
}
