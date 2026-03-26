package model

import "time"

type Task struct {
	TaskID    int       `json:"task_id"`
	Title     string    `json:"title"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (t Task) ID() int {
	return t.TaskID
}
