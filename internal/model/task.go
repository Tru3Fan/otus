package model

import "time"

type Task struct {
	TaskID     int        `json:"task_id"`
	Title      string     `json:"title"`
	Body       string     `json:"body,omitempty"`
	UserID     int        `json:"user_id"`
	Status     string     `json:"status"`
	Deadline   *time.Time `json:"deadline,omitempty"`
	AssignedBy int        `json:"assigned_by,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

func (t Task) ID() int {
	return t.TaskID
}
