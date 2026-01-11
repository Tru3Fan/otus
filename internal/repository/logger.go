package repository

import (
	"context"
	"log"
	"otus/internal/model"
	"time"
)

func LogNew(ctx context.Context) {
	t := time.NewTicker(200 * time.Millisecond)
	defer t.Stop()

	uLast, tLast := 0, 0

	for {
		select {
		case <-ctx.Done():
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
