package service

import (
	"otus/internal/model"
	"otus/internal/repository"
)

func GenerateAndStore() {
	u := model.User{UserID: 1, Username: "Ivan"}
	t := model.Task{TaskID: 1, Title: "Sleep"}

	repository.Add(u)
	repository.Add(t)
}
