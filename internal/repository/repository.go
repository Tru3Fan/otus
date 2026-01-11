package repository

import (
	"fmt"
	"log"
	"otus/internal/model"
	"sync"
	"time"
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

func Add(in <-chan Storable) {

	for m := range in {
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

func LogNew(stop <-chan struct{}) {
	t := time.NewTicker(200 * time.Millisecond)
	defer t.Stop()

	uLast, tLast := 0, 0

	for {
		select {
		case <-stop:
			return
		case <-t.C:
			muUsers.Lock()
			uNew := append([]model.User(nil), users[uLast:]...)
			uLast = len(users)
			muUsers.Unlock()

			muTasks.Lock()
			tNew := append([]model.Task(nil), tasks[tLast:]...)
			tLast = len(tasks)
			muTasks.Unlock()

			for _, u := range uNew {
				log.Println("Add user: ", u.Username)
			}
			for _, t := range tNew {
				log.Println("Add task: ", t.Title)
			}
		}
	}
}
