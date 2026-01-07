package main

import (
	"fmt"
	"otus/internal/model"
)

func main() {
	a := model.Task{
		ID:          1,
		Title:       "Ivan",
		Description: "Ivan",
		AssigneeID:  23,
	}

	b := model.User{
		ID:       8654,
		Username: "john",
	}

	fmt.Println(a)
	fmt.Println(b)
}
