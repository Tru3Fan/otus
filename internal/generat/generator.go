package generat

import (
	"context"
	"otus/internal/model"
	"otus/internal/repository"
	"sync"
)

func GenerateAndCreate(ctx context.Context, out chan<- repository.Storable, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(out)

	for range 4 {

		select {
		case <-ctx.Done():
			return
		case out <- model.User{UserID: 4, Username: "Ivan"}:
		}

		select {
		case <-ctx.Done():
			return
		case out <- model.Task{TaskID: 1, Title: "Dream"}:
		}
	}
}
