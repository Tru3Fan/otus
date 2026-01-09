package model

type Task struct {
	TaskID int
	Title  string
}

func (t Task) ID() int {
	return t.TaskID
}
