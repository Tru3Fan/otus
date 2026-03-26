package csv

import (
	"context"
	"log"
	"otus/internal/model"
	"sync"
	"time"
)

func Counts() (usersCount, tasksCount int) {
	muUsers.Lock()
	usersCount = len(users)
	muUsers.Unlock()

	muTasks.Lock()
	tasksCount = len(tasks)
	muTasks.Unlock()

	return
}

func LogNew(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	t := time.NewTicker(200 * time.Millisecond)
	defer t.Stop()

	uLast, tLast := Counts()

	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			muUsers.Lock()
			if uLast > len(users) {
				uLast = len(users)
			}
			uNew := append([]model.User(nil), users[uLast:]...)
			uLast = len(users)
			muUsers.Unlock()

			muTasks.Lock()
			if tLast > len(tasks) {
				tLast = len(tasks)
			}

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
