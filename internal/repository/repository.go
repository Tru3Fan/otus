package repository

import (
	"context"
	"fmt"
	"otus/internal/model"
	"sync"
)

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
				muUsers.Unlock()
				fmt.Println("Added user: ", v.Username)
			case model.Task:
				muTasks.Lock()
				tasks = append(tasks, v)
				muTasks.Unlock()
				fmt.Println("Added task: ", v.Title)
			default:
				fmt.Println("Unknown type: ", v)
			}
		}
	}
}
