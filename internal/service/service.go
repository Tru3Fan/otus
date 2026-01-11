package service

import (
	"context"
	"otus/internal/model"
	"otus/internal/repository"
)

func GenerateAndCreate(ctx context.Context, out chan<- repository.Storable) {
	defer close(out)

	for range 10 {

		select {
		case <-ctx.Done():
			return
		case out <- model.User{1234, "Dmitriy"}:
		}

		select {
		case <-ctx.Done():
			return
		case out <- model.Task{1211, "Sleep"}:
		}
	}
}
