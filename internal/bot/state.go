package bot

import "time"

type Step int

const (
	StepNone Step = iota
	StepWaitingTaskTitle
	StepWaitingTaskBody
	StepWaitingAssignee
	StepWaitingDeadline
)

type DraftTask struct {
	Title      string
	Body       string
	AssigneeID int
	Deadline   *time.Time
}

type UserState struct {
	Step  Step
	Draft DraftTask
}

func newState() *UserState {
	return &UserState{Step: StepNone}
}

func (b *Bot) getState(userID int64) *UserState {
	if s, ok := b.state[userID]; ok {
		return s
	}
	s := newState()
	b.state[userID] = s
	return s
}
func (b *Bot) resetState(userID int64) {
	b.state[userID] = newState()
}
