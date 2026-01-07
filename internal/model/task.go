package model

import "time"

type Status string

const (
	StatusPending   Status = "pending"
	StatusCompleted Status = "completed"
	statusFailed    Status = "failed"
)

type Task struct {
	ID          int64
	Title       string
	Description string
	AssigneeID  int64
	Deadline    time.Time
	status      Status
	createdAt   time.Time
}
