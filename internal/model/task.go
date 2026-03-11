package model

type Task struct {
	TaskID int    `json:"task_id"`
	Title  string `json:"title"`
}

func (t Task) ID() int {
	return t.TaskID
}
