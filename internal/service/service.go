package service

import (
	"otus/internal/model"
	"otus/internal/repository"
)

func GenerateAndCreate(out chan<- repository.Storable) {

	for range 10 {
		out <- model.User{1234, "Dmitriy"}
		out <- model.Task{1211, "Sleep"}
	}

	close(out)
}
